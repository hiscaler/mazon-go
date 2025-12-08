package mazon

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/aar"
	"github.com/hiscaler/mazon-go/config"
	"github.com/hiscaler/mazon-go/entity"
)

const (
	OK              = 200 // 无错误
	BadRequestError = 400 // 请求错误
	InvalidToken    = 401 // 无效的 Token
	InternalError   = 500 // 内部错误，数据库异常
)

const (
	Version   = "0.0.1"
	userAgent = "Mazon API Client-Golang/" + Version + " (https://github.com/hiscaler/mazon-go)"
	baseUrl   = "https://api.mazonlabel.com/api/svc" // 美正无测试地址
)

type Client struct {
	config      *config.Config // 配置
	httpClient  *resty.Client  // Resty Client
	accessToken string         // AccessToken
	retry       bool           // 是否重新发起请求，如果是重新发起的，需要重新获取 token
	Services    services       // API Services
}

func NewClient(ctx context.Context, cfg config.Config) *Client {
	logger := log.New(os.Stdout, "[ Client ] ", log.LstdFlags|log.Llongfile)
	mazonClient := &Client{
		config: &cfg,
	}
	httpClient := resty.New().
		SetDebug(cfg.Debug).
		SetBaseURL(baseUrl).
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		}).
		SetTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialContext: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).DialContext,
		}).
		SetTimeout(time.Duration(cfg.Timeout) * time.Second).
		OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
			token := mazonClient.accessToken
			if token == "" || mazonClient.retry {
				ar, err := aar.New("mazon.access.token.%s.%s", cfg.AppKey, cfg.AppToken)
				if err == nil {
					token, _ = ar.SetDuration(time.Duration(min(max(cfg.TokenDuration, 1), 10)) * time.Hour).Read()
				}
				if token == "" {
					// 重新获取 Token
					if err = mazonClient.getAccessToken(ctx); err != nil {
						return err
					}
					token = mazonClient.accessToken
					if err = ar.Write([]byte(token)); err != nil {
						logger.Println("[ Error ]", err)
					}
				}
			}
			client.SetHeader("Authorization", token)
			return nil
		}).
		OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
			logger.Printf(`
Header: %#v
Raw Request: %#v
Request: %#v
Raw Response: %#v
Response: %#v
`,
				response.Request.Header,
				response.Request.RawRequest,
				response.Request,
				response.RawResponse,
				response,
			)
			fmt.Printf(`
Header: %#v
Raw Request: %#v
Request: %#v
Raw Response: %#v
Response: %#v
`,
				response.Request.Header,
				response.Request.RawRequest,
				response.Request,
				response.RawResponse,
				response,
			)
			return nil
		}).
		SetRetryCount(2).
		SetRetryWaitTime(2 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second).
		AddRetryCondition(func(response *resty.Response, err error) bool {
			if response == nil {
				return true
			}
			var r NormalResponse
			retry := json.Unmarshal(response.Body(), &r) == nil && r.Code == InvalidToken
			mazonClient.retry = retry
			return retry
		})
	mazonClient.httpClient = httpClient

	xService := service{
		config:     &cfg,
		logger:     logger,
		httpClient: mazonClient.httpClient,
	}
	mazonClient.Services = services{
		Order:         (orderService)(xService),
		Rate:          (rateService)(xService),
		User:          (userService)(xService),
		ShippingLabel: (shippingLabelService)(xService),
		ScanForm:      (scanFormService)(xService),
	}
	return mazonClient
}

type NormalResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
	Result  any    `json:"result"`
}

// accessToken 获取 Access Token 值
func (c *Client) getAccessToken(ctx context.Context) (err error) {
	if c.accessToken != "" {
		return nil
	}

	result := struct {
		NormalResponse
		Result *entity.Token `json:"result"`
	}{}
	httpClient := resty.New().
		SetDebug(c.config.Debug).
		SetBaseURL(baseUrl).
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		}).
		SetTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialContext: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).DialContext,
		})
	resp, err := httpClient.R().
		SetContext(ctx).
		SetBody(map[string]string{
			"app_key":   c.config.AppKey,
			"app_token": c.config.AppToken,
		}).
		SetResult(&result).
		Post("/getToken")
	if err = recheckError(resp, result.NormalResponse, err); err != nil {
		return
	}
	if result.Result != nil {
		c.accessToken = result.Result.AccessToken
	}
	return
}

// errorWrap 错误包装
func errorWrap(code int, message string) error {
	if code == OK || code == 0 {
		return nil
	}

	switch code {
	case InvalidToken:
		message = "无效的 Token"
	default:
		if code == InternalError {
			if message == "" {
				message = "内部错误，请联系美正客服"
			}
		} else {
			message = strings.TrimSpace(message)
			if message == "" {
				message = "Unknown error"
			}
		}
	}
	return fmt.Errorf("%d: %s", code, message)
}

func invalidInput(e error) error {
	var errs validation.Errors
	if !errors.As(e, &errs) {
		return e
	}

	if len(errs) == 0 {
		return nil
	}

	fields := make([]string, 0)
	messages := make([]string, 0)
	for field := range errs {
		fields = append(fields, field)
	}
	sort.Strings(fields)

	for _, field := range fields {
		e1 := errs[field]
		if e1 == nil {
			continue
		}

		var errObj validation.ErrorObject
		if errors.As(e1, &errObj) {
			e1 = errObj
		} else {
			var errs1 validation.Errors
			if errors.As(e1, &errs1) {
				e1 = invalidInput(errs1)
				if e1 == nil {
					continue
				}
			}
		}

		messages = append(messages, e1.Error())
	}
	return errors.New(strings.Join(messages, "; "))
}

func recheckError(resp *resty.Response, result NormalResponse, e error) error {
	if e != nil {
		if errors.Is(e, http.ErrHandlerTimeout) {
			return errorWrap(http.StatusRequestTimeout, e.Error())
		}
		return e
	}

	if resp.IsError() {
		return errorWrap(resp.StatusCode(), resp.Error().(string))
	}

	if result.Code != http.StatusOK {
		return errorWrap(result.Code, result.Message)
	}
	return nil
}
