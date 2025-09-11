package areship

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/areship-go/entity"
)

// 面单服务
type shippingLabelService service

type ShippingLabelDetailRequest struct {
	OrderCode   string `json:"order_code"`   // 订单号
	ReferenceNo string `json:"reference_no"` // 参考号
}

func (m ShippingLabelDetailRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.OrderCode, validation.When(m.ReferenceNo == "", validation.Required.Error("订单号不能为空"))),
		validation.Field(&m.ReferenceNo, validation.When(m.OrderCode == "", validation.Required.Error("参考号不能为空"))),
	)
}

// Detail 获取面单
// http://doc.areship.cn/api-68024258
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
	if err != nil {
		return label, err
	}

	if err = recheckError(resp, res.NormalResponse, err); err != nil {
		return label, err
	}
	return res.Result, nil
}
