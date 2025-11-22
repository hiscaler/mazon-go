package mazon

import (
	"context"
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/hiscaler/mazon-go/entity"
)

// 订单服务
type orderService service

type OrderBox struct {
	Length          float64 `json:"box_length"`                  // 长（支持两位小数）
	Width           float64 `json:"box_width"`                   // 宽（支持两位小数）
	Height          float64 `json:"box_height"`                  // 高（支持两位小数）
	ActualWeight    float64 `json:"box_actual_weight"`           // 箱子重量（单位 kg，支持两位小数）
	Sku             string  `json:"sku,omitempty"`               // SKU
	CnName          string  `json:"cn_name,omitempty"`           // 中文名称
	EngName         string  `json:"eng_name,omitempty"`          // 英文名称
	ApplyCompany    string  `json:"apply_company,omitempty"`     // 申报单位
	ApplyNumber     int     `json:"apply_number,omitempty"`      // 申报数量
	ApplyUnitPrice  float64 `json:"apply_unit_price,omitempty"`  // 申报价格
	ApplyUnitWeight float64 `json:"apply_unit_weight,omitempty"` // 申报重量
	GoodDetail      string  `json:"good_detail,omitempty"`       // 配货信息
	CustomsCode     string  `json:"customs_code,omitempty"`      // 海关编码
	SaleUrl         string  `json:"sale_url,omitempty"`          // 销售链接
	CnMaterial      string  `json:"cn_material,omitempty"`       // 中文材质
	EngMaterial     string  `json:"eng_material,omitempty"`      // 英文材质
	ProduceCountry  string  `json:"produce_country,omitempty"`   // 生产国家
	Remark          string  `json:"remark,omitempty"`            // 生产国家
}

func (m OrderBox) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Length,
			validation.Required.Error("长不能为空"),
			validation.Min(0.01).Error("长不能小于 {{.min}}"),
			validation.Max(999999.99).Error("长不能大于 {{.max}}"),
		),
		validation.Field(&m.Width,
			validation.Required.Error("宽不能为空"),
			validation.Min(0.01).Error("宽不能小于 {{.min}}"),
			validation.Max(999999.99).Error("宽不能大于 {{.max}}"),
		),
		validation.Field(&m.Height,
			validation.Required.Error("高不能为空"),
			validation.Min(0.01).Error("高不能小于 {{.min}}"),
			validation.Max(999999.99).Error("高不能大于 {{.max}}"),
		),
		validation.Field(&m.ActualWeight,
			validation.Required.Error("重量不能为空"),
			validation.Min(0.01).Error("重量不能小于 {{.min}}"),
			validation.Max(999999.99).Error("重量不能大于 {{.max}}"),
		),
		validation.Field(&m.Sku, validation.When(m.Sku != "", validation.Length(1, 35).Error("SKU 不能超过 {{.max}} 个字符"))),
		validation.Field(&m.CnName, validation.When(m.CnName != "", validation.Length(1, 35).Error("中文名称不能超过 {{.max}} 个字符"))),
		validation.Field(&m.EngName, validation.When(m.EngName != "", validation.Length(1, 35).Error("英文名称不能超过 {{.max}} 个字符"))),
		validation.Field(&m.ApplyCompany, validation.When(m.ApplyCompany != "", validation.Length(1, 35).Error("申报单位不能超过 {{.max}} 个字符"))),
		validation.Field(&m.ApplyNumber, validation.When(m.ApplyNumber > 0, validation.Min(1).Error("申报数量不能小于 {{.min}}"))),
		validation.Field(&m.ApplyUnitPrice, validation.When(m.ApplyUnitPrice > 0, validation.Min(0.01).Error("申报价格不能小于 {{.min}}"))),
		validation.Field(&m.ApplyUnitWeight, validation.When(m.ApplyUnitWeight > 0, validation.Min(0.01).Error("申报重量不能小于 {{.min}}"))),
		validation.Field(&m.GoodDetail, validation.When(m.GoodDetail != "", validation.Length(1, 35).Error("配货信息不能超过 {{.max}} 个字符"))),
		validation.Field(&m.CustomsCode, validation.When(m.CustomsCode != "", validation.Length(1, 35).Error("海关编码不能超过 {{.max}} 个字符"))),
		validation.Field(&m.SaleUrl, validation.When(m.SaleUrl != "", validation.Length(1, 35).Error("销售链接不能超过 {{.max}} 个字符"))),
		validation.Field(&m.CnMaterial, validation.When(m.CnMaterial != "", validation.Length(1, 35).Error("中文材质不能超过 {{.max}} 个字符"))),
		validation.Field(&m.EngMaterial, validation.When(m.EngMaterial != "", validation.Length(1, 35).Error("英文材质不能超过 {{.max}} 个字符"))),
		validation.Field(&m.ProduceCountry, validation.When(m.ProduceCountry != "", validation.Length(1, 35).Error("生产国家不能超过 {{.max}} 个字符"))),
		validation.Field(&m.Remark, validation.When(m.Remark != "", validation.Length(1, 35).Error("生产国家不能超过 {{.max}} 个字符"))),
	)
}

type CreateOrderRequest struct {
	ReferenceNO        string                 `json:"reference_no"`                    // 订单参考号，唯一
	SMCode             string                 `json:"sm_code"`                         // 物流产品代码，请咨询您的销售代表获取
	Remark             string                 `json:"remark,omitempty"`                // 订单备注
	OAFirstname        string                 `json:"oa_firstname"`                    // 收件人姓名
	OACompany          string                 `json:"oa_company,omitempty"`            // 收件人公司
	OATelephone        string                 `json:"oa_telphone"`                     // 收件人电话
	OACountry          string                 `json:"oa_country"`                      // 收件人国家（国家二字码）
	OAState            string                 `json:"oa_state"`                        // 收件人州/省
	OACity             string                 `json:"oa_city"`                         // 收件人城市
	OAPostcode         string                 `json:"oa_postcode"`                     // 收件人邮编
	OAStreetAddress1   string                 `json:"oa_street_address1,omitempty"`    // 收件人地址 1
	OAStreetAddress2   string                 `json:"oa_street_address2,omitempty"`    // 收件人地址 2
	IsMoreBox          int                    `json:"is_more_box"`                     // 包裹类型: 固定值，请传 1
	SignatureService   string                 `json:"signature_service,omitempty"`     // 签名服务（是否需要签名服务：ASS为 成人签名 ，SSF为 普通签名，不需要可以不传该字段）
	PickUp             int                    `json:"pick_up,omitempty"`               // 是否提货 1：是，0：否，不传默认为否, 传1（是）需要物流产品支持，物流产品不支持传1(是)也无效
	WeightUnitType     int                    `json:"weight_unit_type,omitempty"`      // 包裹单位类型（1-英制(INCH/LBS) 2-公制(CM/KG) 默认为2）
	LabelCustomType    string                 `json:"label_custom_type,omitempty"`     // 自定义面单打印类型: 1为都打印 2为打印参考号 3为仅仅打印备注 默认为1
	MailingDate        string                 `json:"mailing_date,omitempty"`          // 发货日期 格式为yyyy-MM-dd
	LabelImageFormat   string                 `json:"label_image_format,omitempty"`    // 面单格式: PDF、ZPL(打印机格式) 默认为PDF
	HasUpsLabelCropped string                 `json:"has_ups_label_cropped,omitempty"` // UPS面单是否裁剪,true为裁剪,false为不裁剪 不传默认为true
	GenerateGxEvent    string                 `json:"generate_gx_event,omitempty"`     // 是否生成gx预报轨迹 不传默认为true
	BoxList            []OrderBox             `json:"box_list"`                        // 包裹信息
	ShipperAddress     *entity.ShipperAddress `json:"shipper_address,omitempty"`       // 发件人信息（发件人信息必须与我司备案信息完全一致）
	ShipperCode        string                 `json:"shipper_code,omitempty"`          // 发件人编码（发件人信息与发件人编码同时存在时，以发件人信息为准）
	ReturnAddress      *entity.ReturnAddress  `json:"return_address"`                  // 退件地址信息
}

func (m CreateOrderRequest) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.ReferenceNO,
			validation.Required.Error("订单参考号不能为空"),
			validation.Length(1, 35).Error("订单参考号不能超过 {{.max}} 个字符"),
		),
		validation.Field(&m.SMCode, validation.Required.Error("物流产品代码不能为空")),
		validation.Field(&m.Remark, validation.When(m.Remark != "", validation.Length(1, 35).Error("备注不能超过 {{.max}} 个字符"))),
		validation.Field(&m.OAFirstname,
			validation.Required.Error("收件人不能为空"),
			validation.Length(3, 35).Error("收件人长度必须在 {{.min}} ~ {{.max}} 个字符"),
		),
		validation.Field(&m.OACompany, validation.When(m.OACompany != "", validation.Length(1, 35).Error("收件人公司不能超过 {{.max}} 个字符"))),
		validation.Field(&m.OAStreetAddress1,
			validation.Required.Error("收件人地址1不能为空"),
			validation.Length(1, 35).Error("收件人地址1长度不能超过 {{.max}} 个字符"),
		),
		validation.Field(&m.OAStreetAddress2,
			validation.When(m.OAStreetAddress2 != "", validation.Length(1, 35).Error("收件人地址2长度不能超过 {{.max}} 个字符")),
		),
		validation.Field(&m.OAPostcode, validation.Required.Error("收件人邮编不能为空")),
		validation.Field(&m.OAState, validation.Required.Error("收件人州不能为空")),
		validation.Field(&m.OACity, validation.Required.Error("收件人城市不能为空")),
		validation.Field(&m.OACountry, validation.Required.Error("收件人国家（国家二字码）不能为空")),
		validation.Field(&m.OATelephone,
			validation.Required.Error("收件人电话不能为空"),
			validation.Length(10, 15).Error("收件人电话长度必须在 {{.min}} ~ {{.max}} 个字符"),
		),
		validation.Field(&m.SignatureService,
			validation.When(m.SignatureService != "", validation.In("ASS", "SSF").Error("签名服务参数错误")),
		),
		validation.Field(&m.WeightUnitType, validation.In(1, 2).Error("包裹单位类型参数错误")),
		validation.Field(&m.BoxList, validation.Required.Error("包裹信息不能为空")),
		validation.Field(&m.ShipperAddress, validation.When(m.ShipperCode == "", validation.Required.Error("发件人信息和编码必须填写一个"))),
		validation.Field(&m.ShipperCode, validation.When(m.ShipperAddress == nil, validation.Required.Error("发件人信息和编码必须填写一个"))),
	)
}

// Create 创建订单
// https://www.mazonlabel.com/docs/orderapi/%E5%88%9B%E5%BB%BA%E8%AE%A2%E5%8D%95.html
func (s orderService) Create(ctx context.Context, req CreateOrderRequest) (createRes entity.OrderCreateResult, err error) {
	if err = req.Validate(); err != nil {
		err = invalidInput(err)
		return
	}

	res := struct {
		NormalResponse
		Result entity.OrderCreateResult `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&res).
		Post("/createOrder")
	if err = recheckError(resp, res.NormalResponse, err); err != nil {
		return
	}
	if res.Result.LabelStatus == 0 {
		return createRes, errors.New(res.Message)
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

// Query 根据查询条件筛选符合条件的订单列表数据
// https://www.mazonlabel.com/docs/orderapi/%E8%8E%B7%E5%8F%96%E8%AE%A2%E5%8D%95%E4%BF%A1%E6%81%AF.html
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
	OrderCode   string `json:"order_code,omitempty"`   // 订单号
	ReferenceNo string `json:"reference_no,omitempty"` // 参考号
}

func (m CancelOrderRequest) validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.OrderCode, validation.When(m.ReferenceNo == "", validation.Required.Error("订单号不能为空"))),
		validation.Field(&m.ReferenceNo, validation.When(m.OrderCode == "", validation.Required.Error("参考号不能为空"))),
	)
}

// Cancel 取消订单
// https://www.mazonlabel.com/docs/orderapi/%E5%8F%96%E6%B6%88%E8%AE%A2%E5%8D%95.html
// 在订单草稿、已预报、已提交（未在预报执行中）状态时可以进行取消订单操作。
//
// 返回值
// 5：取消中、6：已取消
func (s orderService) Cancel(ctx context.Context, req CancelOrderRequest) (int, error) {
	if err := req.validate(); err != nil {
		return -1, invalidInput(err)
	}

	res := struct {
		NormalResponse
		Result int `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(req).
		SetResult(&res).
		Post("/cancelOrder")
	if err = recheckError(resp, res.NormalResponse, err); err != nil {
		return -1, err
	}
	return res.Result, nil
}
