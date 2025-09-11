package areship

import (
	"log"

	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/areship-go/config"
)

type service struct {
	config     *config.Config // Config
	logger     *log.Logger    // Logger
	httpClient *resty.Client  // HTTP client
}

// API Services
type services struct {
	Order         orderService         // 订单
	UserInfo      userInfoService      // 用户信息
	ShippingLabel shippingLabelService // 面单
}
