package entity

// Order 订单
type Order struct {
	ReferenceNo    string `json:"reference_no"`        // 参考号
	OrderCode      string `json:"order_code"`          // 订单号
	AddTime        string `json:"add_time"`            // 添加时间
	OrderStatus    int    `json:"order_status,string"` // 订单状态
	Remark         string `json:"remark"`              // 备注
	Firstname      string `json:"firstname"`           // 收件人姓名
	Company        string `json:"company"`             // 收件人公司
	Country        string `json:"country"`             // 收件人国家
	Postcode       string `json:"postcode"`            // 收件人邮编
	State          string `json:"state"`               // 收件人州
	City           string `json:"city"`                // 收件人城市
	StreetAddress1 string `json:"street_address1"`     // 收件人街道
	TelPhone       string `json:"telphone"`            // 收件人电话号码
}
