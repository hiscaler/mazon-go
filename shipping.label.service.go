package mazon

import (
	"context"
	"errors"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/mazon-go/entity"
)

// 面单服务
type shippingLabelService service

type ShippingLabelDetailRequest struct {
	OrderCode   string `json:"order_code,omitempty"`   // 订单号
	ReferenceNo string `json:"reference_no,omitempty"` // 参考号
}

func (m ShippingLabelDetailRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.OrderCode, validation.When(m.ReferenceNo == "", validation.Required.Error("订单号不能为空"))),
		validation.Field(&m.ReferenceNo, validation.When(m.OrderCode == "", validation.Required.Error("参考号不能为空"))),
	)
}

// Detail 获取面单
// https://www.mazonlabel.com/docs/orderapi/%E8%8E%B7%E5%8F%96%E9%9D%A2%E5%8D%95.html
func (s shippingLabelService) Detail(ctx context.Context, req ShippingLabelDetailRequest) (label entity.ShippingLabel, err error) {
	if err = req.validate(); err != nil {
		return label, invalidInput(err)
	}

	res := struct {
		NormalResponse
		Result entity.ShippingLabel `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&res).
		Post("/getLabel")
	if err = recheckError(resp, res.NormalResponse, err); err != nil {
		return label, err
	}
	return res.Result, nil
}

// Query 根据物流单号获取面单信息
// https://www.mazonlabel.com/docs/orderapi/%E6%A0%B9%E6%8D%AE%E7%89%A9%E6%B5%81%E5%8D%95%E5%8F%B7%E8%8E%B7%E5%8F%96%E9%9D%A2%E5%8D%95%E4%BF%A1%E6%81%AF.html
func (s shippingLabelService) Query(ctx context.Context, trackingNumbers ...string) (labels []entity.LogisticsLabel, err error) {
	numbers := make([]string, 0, len(trackingNumbers))
	for _, number := range trackingNumbers {
		number = strings.TrimSpace(number)
		if number != "" {
			numbers = append(numbers, number)
		}
	}
	if len(numbers) == 0 {
		return labels, errors.New("无效的跟踪号")
	}

	res := struct {
		NormalResponse
		Result []entity.LogisticsLabel `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]string{"tracking_number": strings.Join(numbers, ",")}).
		SetResult(&res).
		Post("/getLabelInfo")
	if err = recheckError(resp, res.NormalResponse, err); err != nil {
		return labels, err
	}
	return res.Result, nil
}
