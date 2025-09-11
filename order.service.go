package areship

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/areship-go/entity"
	"gopkg.in/guregu/null.v4"
)

type orderService service

type OrderBox struct {
	Length       float64 `json:"box_length"`        // 长（支持两位小数）
	Width        float64 `json:"box_width"`         // 宽（支持两位小数）
	Height       float64 `json:"box_height"`        // 高（支持两位小数）
	ActualWeight float64 `json:"box_actual_weight"` // 箱子重量（单位 kg，支持两位小数）
}

type CreateOrderRequest struct {
	ReferenceNO      string         `json:"reference_no"`       // 订单参考号，唯一
	SMCode           string         `json:"sm_code"`            // 物流产品代码，请咨询您的销售代表获取
	Remark           null.String    `json:"remark"`             // 订单备注
	OAFirstname      string         `json:"oa_firstname"`       // 收件人姓名
	OACompany        null.String    `json:"oa_company"`         // 收件人公司
	OAStreetAddress1 string         `json:"oa_street_address1"` // 收件人地址 1
	OAStreetAddress2 null.String    `json:"oa_street_address2"` // 收件人地址 2
	OAPostcode       string         `json:"oa_postcode"`        // 收件人邮编
	OAState          string         `json:"oa_state"`           // 收件人州/省
	OACity           string         `json:"oa_city"`            // 收件人城市
	OACountry        string         `json:"oa_country"`         // 收件人国家（国家二字码）
	OATelephone      string         `json:"oa_telphone"`        // 收件人电话
	IsMoreBox        int            `json:"is_more_box"`        // 包裹类型
	SignatureService null.String    `json:"signature_service"`  // 签名服务（是否需要签名服务：ASS为 成人签名 ，SSF为 普通签名，不需要可以不传该字段）
	WeightUnitType   int            `json:"weight_unit_type"`   // 包裹单位类型（1-英制(INCH/LBS) 2-公制(CM/KG) 默认为2）
	BoxList          []OrderBox     `json:"box_list"`           // 包裹信息
	ShipperAddress   entity.Address `json:"shipper_address"`    // 发件人信息（发件人信息必须与我司备案信息完全一致）
	ShipperCode      string         `json:"shipper_code"`       // 发件人编码（发件人信息与发件人编码同时存在时，以发件人信息为准）
}

func (m CreateOrderRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ReferenceNO,
			validation.Required.Error("订单参考号不能为空"),
			validation.Length(1, 35).Error("订单参考号不能超过 {.max} 个字符"),
		),
		validation.Field(&m.SMCode, validation.Required.Error("物流产品代码不能为空")),
		validation.Field(&m.Remark, validation.When(m.Remark.Valid, validation.Length(1, 35).Error("备注不能超过 {.max} 个字符"))),
		validation.Field(&m.OAFirstname,
			validation.Required.Error("收件人不能为空"),
			validation.Length(3, 35).Error("收件人长度必须在 {.min} ~ {.max} 个字符"),
		),
		validation.Field(&m.OACompany, validation.When(m.OACompany.Valid, validation.Length(1, 35).Error("收件人公司不能超过 {.max} 个字符"))),
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
		validation.Field(&m.ShipperAddress, validation.Required.Error("发件人信息不能为空")),
		validation.Field(&m.ShipperCode, validation.Required.Error("发件人编码不能为空")),
	)
}

// Create 创建订单
// http://doc.areship.cn/api-68024102
func (s orderService) Create(ctx context.Context, req CreateOrderRequest) (createRes entity.CreateOrderResult, err error) {
	if err = req.Validate(); err != nil {
		err = invalidInput(err)
		return
	}

	res := struct {
		NormalResponse
		Result entity.CreateOrderResult `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(req).
		Post("/createOrder")
	if err = recheckError(resp, res.NormalResponse, err); err != nil {
		return
	}
	return res.Result, nil
}

type OrderQueryRequest struct {
	Type        int    `json:"type"`                   // 类型（1 代表按时间搜索、2 代表按票搜索）
	OrderCode   string `json:"order_code,omitempty"`   // 订单号
	ReferenceNo string `json:"reference_no,omitempty"` // 参考号
	DateFrom    string `json:"date_from,omitempty"`    // 开始时间（格式：2023-06-01 00:00:00）
	DateTo      string `json:"date_to,omitempty"`      // 结束时间（格式：2023-06-01 00:00:00）
}

func (m OrderQueryRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Type, validation.In(1, 2).Error("类型参数错误")),
		validation.Field(&m.DateFrom, validation.When(m.DateFrom != "", validation.Date(time.DateTime).Error("开始时间格式错误"))),
		validation.Field(&m.DateTo, validation.When(m.DateTo != "", validation.Date(time.DateTime).Error("结束时间格式错误"))),
	)
}

// Query 查询订单
// http://doc.areship.cn/api-106489828
func (s orderService) Query(ctx context.Context, req OrderQueryRequest) ([]entity.Order, error) {
	if err := req.validate(); err != nil {
		return nil, invalidInput(err)
	}

	res := struct {
		NormalResponse
		Result []entity.Order `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&res).
		Post("/getOrderInfo")
	if err != nil {
		return nil, err
	}

	if err = recheckError(resp, res.NormalResponse, err); err != nil {
		return nil, err
	}
	return res.Result, nil
}

type CancelOrderRequest struct {
	OrderCode   string `json:"order_code"`   // 订单号
	ReferenceNo string `json:"reference_no"` // 参考号
}

func (m CancelOrderRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.OrderCode, validation.When(m.ReferenceNo == "", validation.Required.Error("订单号不能为空"))),
		validation.Field(&m.ReferenceNo, validation.When(m.OrderCode == "", validation.Required.Error("参考号不能为空"))),
	)
}

// Cancel 取消订单
// http://doc.areship.cn/api-68024502
// 在订单草稿、已预报、已提交（未在预报执行中）状态时可以进行取消订单操作。
func (s orderService) Cancel(ctx context.Context, req CancelOrderRequest) ([]string, error) {
	if err := req.validate(); err != nil {
		return nil, invalidInput(err)
	}

	res := struct {
		NormalResponse
		Result []string `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&res).
		Post("/cancelOrder")
	if err != nil {
		return nil, err
	}

	if err = recheckError(resp, res.NormalResponse, err); err != nil {
		return nil, err
	}
	return res.Result, nil
}
