package mazon

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/mazon-go/entity"
	"gopkg.in/guregu/null.v4"
)

// 运费服务
type rateService service

type RateCalcOrderBox struct {
	Length       float64 `json:"box_length"`        // 长（支持两位小数）
	Width        float64 `json:"box_width"`         // 宽（支持两位小数）
	Height       float64 `json:"box_height"`        // 高（支持两位小数）
	ActualWeight float64 `json:"box_actual_weight"` // 箱子重量（单位 kg，支持两位小数）
}

type RateCalcRequest struct {
	ReferenceNO      string                 `json:"reference_no"`                 // 订单参考号，唯一
	SMCode           string                 `json:"sm_code"`                      // 物流产品代码，请咨询您的销售代表获取
	Remark           string                 `json:"remark,omitempty"`             // 订单备注
	OAFirstname      string                 `json:"oa_firstname"`                 // 收件人姓名
	OACompany        string                 `json:"oa_company,omitempty"`         // 收件人公司
	OATelephone      string                 `json:"oa_telphone"`                  // 收件人电话
	OACountry        string                 `json:"oa_country"`                   // 收件人国家（国家二字码）
	OAState          string                 `json:"oa_state"`                     // 收件人州/省
	OACity           string                 `json:"oa_city"`                      // 收件人城市
	OAPostcode       string                 `json:"oa_postcode"`                  // 收件人邮编
	OAStreetAddress1 string                 `json:"oa_street_address1,omitempty"` // 收件人地址 1
	OAStreetAddress2 null.String            `json:"oa_street_address2,omitempty"` // 收件人地址 2
	IsMoreBox        int                    `json:"is_more_box"`                  // 包裹类型
	SignatureService null.String            `json:"signature_service,omitempty"`  // 签名服务（是否需要签名服务：ASS为 成人签名 ，SSF为 普通签名，不需要可以不传该字段）
	PickUp           int                    `json:"pick_up,omitempty"`            // 是否提货 1：是，0：否，不传默认为否, 传1（是）需要物流产品支持，物流产品不支持传1(是)也无效
	WeightUnitType   int                    `json:"weight_unit_type,omitempty"`   // 包裹单位类型（1-英制(INCH/LBS) 2-公制(CM/KG) 默认为2）
	BoxList          []RateCalcOrderBox     `json:"box_list"`                     // 包裹信息
	ShipperAddress   *entity.ShipperAddress `json:"shipper_address,omitempty"`    // 发件人信息（发件人信息必须与我司备案信息完全一致）
	ShipperCode      string                 `json:"shipper_code,omitempty"`       // 发件人编码（发件人信息与发件人编码同时存在时，以发件人信息为准）
}

func (m RateCalcRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ReferenceNO,
			validation.Required.Error("订单参考号不能为空"),
			validation.Length(1, 35).Error("订单参考号不能超过 {.max} 个字符"),
		),
		validation.Field(&m.SMCode, validation.Required.Error("物流产品代码不能为空")),
		validation.Field(&m.Remark, validation.When(m.Remark != "", validation.Length(1, 35).Error("备注不能超过 {.max} 个字符"))),
		validation.Field(&m.OAFirstname,
			validation.Required.Error("收件人不能为空"),
			validation.Length(3, 35).Error("收件人长度必须在 {.min} ~ {.max} 个字符"),
		),
		validation.Field(&m.OACompany, validation.When(m.OACompany != "", validation.Length(1, 35).Error("收件人公司不能超过 {.max} 个字符"))),
		validation.Field(&m.OAStreetAddress1,
			validation.Required.Error("收件人地址1不能为空"),
			validation.Length(1, 35).Error("收件人地址1长度不能超过 {.max} 个字符"),
		),
		validation.Field(&m.OAStreetAddress2,
			validation.When(m.OAStreetAddress2.Valid, validation.Length(1, 35).Error("收件人地址2长度不能超过 {.max} 个字符")),
		),
		validation.Field(&m.OAPostcode, validation.Required.Error("收件人邮编不能为空")),
		validation.Field(&m.OAState, validation.Required.Error("收件人州不能为空")),
		validation.Field(&m.OACity, validation.Required.Error("收件人城市不能为空")),
		validation.Field(&m.OACountry, validation.Required.Error("收件人国家（国家二字码）不能为空")),
		validation.Field(&m.OATelephone,
			validation.Required.Error("收件人电话不能为空"),
			validation.Length(10, 15).Error("收件人电话长度必须在 {.min} ~ {.max} 个字符"),
		),
		validation.Field(&m.SignatureService,
			validation.When(m.SignatureService.Valid, validation.In("ASS", "SSF").Error("签名服务参数错误")),
		),
		validation.Field(&m.WeightUnitType,
			validation.In(1, 2).Error("包裹单位类型参数错误")),
		validation.Field(&m.BoxList, validation.Required.Error("包裹信息不能为空")),
		validation.Field(&m.ShipperAddress, validation.When(m.ShipperCode == "", validation.Required.Error("发件人信息和编码必须填写一个"))),
		validation.Field(&m.ShipperCode, validation.When(m.ShipperAddress == nil, validation.Required.Error("发件人信息和编码必须填写一个"))),
	)
}

// Calc 提交订单预报参数进行费用试算
// https://www.mazonlabel.com/docs/orderapi/%E8%B4%B9%E7%94%A8%E8%AF%95%E7%AE%97.html
func (s rateService) Calc(ctx context.Context, req RateCalcRequest) (calcResult entity.RateCalcResult, err error) {
	if err = req.Validate(); err != nil {
		err = invalidInput(err)
		return
	}

	res := struct {
		NormalResponse
		Result entity.RateCalcResult `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&res).
		Post("/rates")
	if err = recheckError(resp, res.NormalResponse, err); err != nil {
		return
	}
	return res.Result, nil
}
