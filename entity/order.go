package entity

// Order 订单
type Order struct {
	ReferenceNo    string `json:"reference_no"`
	OrderCode      string `json:"order_code"`
	AddTime        string `json:"add_time"`
	OrderStatus    int    `json:"order_status"`
	Remark         string `json:"remark"`
	Firstname      string `json:"firstname"`
	Company        string `json:"company"`
	Country        string `json:"country"`
	Postcode       string `json:"postcode"`
	State          string `json:"state"`
	City           string `json:"city"`
	StreetAddress1 string `json:"street_address1"`
	TelPhone       string `json:"telphone"`
}
