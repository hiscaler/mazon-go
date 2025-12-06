package mazon

import (
	"context"
	"crypto/md5"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-resty/resty/v2"
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
)

type Client struct {
	config      *config.Config // 配置
	httpClient  *resty.Client  // Resty Client
	accessToken string         // AccessToken
	Services    services       // API Services
}

func NewClient(ctx context.Context, cfg config.Config) *Client {
	logger := log.New(os.Stdout, "[ Client ] ", log.LstdFlags|log.Llongfile)
	mazonClient := &Client{
		config: &cfg,
	}
	httpClient := resty.New().
		SetDebug(cfg.Debug).
		SetBaseURL("https://api.mazonlabel.com/api/svc").
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
			"User-Agent":   userAgent,
		}).
		SetTransport(&http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			DialContext: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).DialContext,
		}).
		SetTimeout(time.Duration(cfg.Timeout) * time.Second).
		OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
			if mazonClient.accessToken == "" {
				refreshAccessToken := true
				var filename string
				h := md5.New()
				_, err := io.WriteString(h, fmt.Sprintf("mazon.access.token.%s.%s", cfg.AppKey, cfg.AppToken))
				if err == nil {
					filename = path.Join(os.TempDir(), fmt.Sprintf("%x", h.Sum(nil)))
					var finfo fs.FileInfo
					// AccessToken 实际过期时间为 12 小时，当前 10 个小时算过期，会重新去获取 Token
					if finfo, err = os.Stat(filename); !os.IsNotExist(err) && !finfo.IsDir() && finfo.ModTime().Add(10*time.Hour).After(time.Now()) {
						var b []byte
						if b, err = os.ReadFile(filename); err == nil {
							refreshAccessToken = false
							mazonClient.accessToken = string(b)
						}
					}
				}
				if refreshAccessToken {
					err = mazonClient.getAccessToken(ctx)
					if err != nil {
						return err
					}
					if filename != "" {
						err = os.WriteFile(filename, []byte(mazonClient.accessToken), 0644)
						if err != nil {
							logger.Println("[ Error ]", err)
						}
					}
				}
			}
			client.SetHeader("Authorization", mazonClient.accessToken)
			return nil
		}).
		SetRetryCount(2).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second)

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
		SetBaseURL("https://www.mazonlabel.com/api/svc").
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
			"User-Agent":   userAgent,
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
