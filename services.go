package mazon

import (
	"log/slog"

	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/mazon-go/config"
)

type service struct {
	config     *config.Config // Config
	logger     *slog.Logger   // Log
	httpClient *resty.Client  // HTTP client
}

// API Services
type services struct {
	Order         orderService         // 订单服务
	Rate          rateService          // 运费服务
	User          userService          // 用户服务
	ShippingLabel shippingLabelService // 面单服务
	ScanForm      scanFormService      // ScanForm 服务
}
