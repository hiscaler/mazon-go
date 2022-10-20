package areship

import (
	"github.com/go-resty/resty/v2"
	"github.com/hiscaler/areship-go/config"
	"log"
)

type service struct {
	config     *config.Config // Config
	logger     *log.Logger    // Logger
	httpClient *resty.Client  // HTTP client
}

// API Services
type services struct {
	Order orderService // 订单
}
