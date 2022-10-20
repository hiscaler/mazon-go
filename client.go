package areship

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/areship-go/config"
	"log"
	"os"
	"strings"
	"time"
)

const (
	OK                   = 200 // 无错误
	ServiceNotFoundError = 400 // 服务不存在
	InternalError        = 500 // 内部错误，数据库异常
)

const (
	Version   = "0.0.1"
	userAgent = "AreShip API Client-Golang/" + Version + " (https://github.com/hiscaler/areship-go)"
)

type AreShip struct {
	config      *config.Config // 配置
	httpClient  *resty.Client  // Resty Client
	accessToken string         // AccessToken
	Services    services       // API Services
}

func NewClient(config config.Config) *AreShip {
	logger := log.New(os.Stdout, "[ AreShip ] ", log.LstdFlags|log.Llongfile)
	areShipClient := &AreShip{
		config: &config,
	}
	httpClient := resty.New().
		SetDebug(areShipClient.config.Debug).
		SetBaseURL("http://www.areship.cn/api/svc").
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
			"User-Agent":   userAgent,
		})

	httpClient.SetTimeout(time.Duration(config.Timeout) * time.Second).
		OnBeforeRequest(func(client *resty.Client, request *resty.Request) error {
			if areShipClient.accessToken == "" {
				areShipClient.getAccessToken()
			}
			client.SetAuthToken(areShipClient.accessToken)
			return nil
		}).
		SetRetryCount(2).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(10 * time.Second)

	areShipClient.httpClient = httpClient

	xService := service{
		config:     &config,
		logger:     logger,
		httpClient: areShipClient.httpClient,
	}
	areShipClient.Services = services{
		Order: (orderService)(xService),
	}
	return areShipClient
}

type NormalResponse struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}

type tokenResult struct {
	AccessToken string `json:"access_token"` // 授权 Token
	UserInfo    struct {
		UID           string `json:"u_id"`            // 用户 ID
		UAccount      string `json:"u_account"`       // 用户账号
		UCustomerCode string `json:"u_customer_code"` // 用户客户代码
	} `json:"user_info"` // 用户信息
}

// accessToken 获取 Token 值
// force 参数为 true 的情况下，会强制重新获取 token，为 false 的情况下根据已有的 token 数据是否过期而采取重新获取或者续期处理。
// 当前通过测试发现领星对 token 的过期时间处理并不是很准确，故当前总是重新获取 token.
func (lx *AreShip) getAccessToken() (err error) {
	if lx.accessToken != "" {
		return nil
	}

	result := struct {
		NormalResponse
		Result tokenResult `json:"result"`
	}{}
	httpClient := resty.New().
		SetDebug(lx.config.Debug).
		SetBaseURL("http://www.areship.cn/api/svc").
		SetHeaders(map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
			"User-Agent":   userAgent,
		})

	resp, err := httpClient.R().SetResult(&result).Post("/getToken")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if result.Code == OK {
			lx.accessToken = result.Result.AccessToken
			return nil
		}
		err = ErrorWrap(result.Code, result.Message)
	} else {
		err = fmt.Errorf("%s: %s", resp.Status(), string(resp.Body()))
	}
	return
}

// ErrorWrap 错误包装
func ErrorWrap(code int, message string) error {
	if code == OK || code == 0 {
		return nil
	}

	switch code {
	case ServiceNotFoundError:
		message = "服务不存在"
	default:
		if code == InternalError {
			if message == "" {
				message = "内部错误，请联系 AreShip 客服"
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
