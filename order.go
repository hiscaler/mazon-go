package areship

import (
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type orderService service

type OrderBox struct {
	Length float64 `json:"box_length"`        // 箱子长（单位 cm）
	Width  float64 `json:"box_width"`         // 箱子宽（单位 cm）
	Height float64 `json:"box_height"`        // 箱子高（单位 cm）
	Weight float64 `json:"box_actual_weight"` // 箱子重量（单位 kg）
}

type CreateOrderRequest struct {
	ReferenceNO         string     `json:"reference_no"`          // 订单参考号，唯一
	SMCode              string     `json:"sm_code"`               // 物流产品代码，请咨询您的销售代表获取
	Remark              string     `json:"remark"`                // 订单备注
	CallbackURL         string     `json:"callback_url"`          // 订单状态变更回调地址，必须是可外网访问的地址
	POCode              string     `json:"po_code"`               // PO code
	VATCode             string     `json:"vat_code"`              // VAT code
	ParcelQuantity      int        `json:"parcel_quantity"`       // 内件数
	ParcelDeclaredValue float64    `json:"parcel_declared_value"` // 申报价值
	OAFirstname         string     `json:"oa_firstname"`          // 收件人
	OACompany           string     `json:"oa_company"`            // 收件人公司
	OAStreetAddress1    string     `json:"oa_street_address1"`    // 收件人地址 1
	OAStreetAddress2    string     `json:"oa_street_address2"`    // 收件人地址 2
	OAStreetAddress3    string     `json:"oa_street_address3"`    // 收件人地址 3
	OAPostcode          string     `json:"oa_postcode"`           // 收件人邮编
	OAState             string     `json:"oa_state"`              // 收件人州/省
	OACity              string     `json:"oa_city"`               // 收件人城市
	OACountry           string     `json:"oa_country"`            // 收件人国家（国家二字码）
	OADoorplate         string     `json:"oa_doorplate"`          // 收件人门牌号
	OATelephone         string     `json:"oa_telphone"`           // 收件人电话
	IsMoreBox           bool       `json:"is_more_box"`           // 是否为一票多箱（0：否、1：是。当为一票多箱时，box_list 必填）
	BoxList             []OrderBox `json:"box_list"`
	SOActualWeight      float64    `json:"so_actual_weight"` // 单票重量（单位 kg）
	SOLength            float64    `json:"so_length"`        // 单票长（单位 cm）
	SOWidth             float64    `json:"so_width"`         // 单票宽（单位 cm）
	SOHeight            float64    `json:"so_height"`        // 单票高（单位 cm）
}

func (m CreateOrderRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ReferenceNO, validation.Required.Error("订单参考号不能为空")),
		validation.Field(&m.SMCode, validation.Required.Error("物流产品代码不能为空")),
		validation.Field(&m.OAFirstname, validation.Required.Error("收件人不能为空")),
		validation.Field(&m.OAStreetAddress1, validation.Required.Error("收件人地址 1 不能为空")),
		validation.Field(&m.OAPostcode, validation.Required.Error("收件人邮编不能为空")),
		validation.Field(&m.OAState, validation.Required.Error("收件人州/省不能为空")),
		validation.Field(&m.OACity, validation.Required.Error("收件人城市不能为空")),
		validation.Field(&m.OACountry, validation.Required.Error("收件人国家（国家二字码）不能为空")),
		validation.Field(&m.OATelephone, validation.Required.Error("收件人电话不能为空")),
		validation.Field(&m.BoxList, validation.When(m.IsMoreBox, validation.Required.Error("箱子数据不能为空"))),
	)
}

type CreateOrderResult struct {
	OrderCode string `json:"order_code"` // 订单号，后续的操作需要以此订单为依据条件
	Fee       []struct {
		FTCode       string  `json:"ft_code"`       // 费用代码
		FTName       string  `json:"ft_name"`       // 费用名称
		CurrencyCode string  `json:"currency_code"` // 币种
		Amount       float64 `json:"amount"`        // 费用金额
	} `json:"fee"` // 订单费用信息
	Labels struct {
		TrackingNumber string `json:"tracking_number"` // 跟踪号
		LabelURL       string `json:"label_url"`       // 面单文件 URL 地址
		FileType       string `json:"file_type"`       // 面单文件类型（比如 pdf 等）
	} `json:"labels"` // 当创建实时预报时返回跟踪号及面单文件，否则为 null
}

// Create creates an order
func (s orderService) Create(req CreateOrderRequest) (createRes CreateOrderResult, err error) {
	if err = req.Validate(); err != nil {
		return
	}

	res := struct {
		NormalResponse
		Result CreateOrderResult `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetBody(req).
		Post("/createOrder")
	if err != nil {
		return
	}

	if err = json.Unmarshal(resp.Body(), &res); err == nil {
		createRes = res.Result
	}
	return
}
