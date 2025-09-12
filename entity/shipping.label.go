package entity

import "gopkg.in/guregu/null.v4"

// ShippingLabel 面单信息
type ShippingLabel struct {
	ReferenceNo        string      `json:"reference_no"`       // 参考号
	OrderCode          string      `json:"order_code"`         // 订单号
	OrderAddressType   string      `json:"order_address_type"` // 地址类型（Commercial 商业 Residential 住宅）
	OrderStatus        int         `json:"order_status"`       // 订单状态（1 已提交、2 已预报）
	OrderSubStatus     string      `json:"order_sub_status"`
	OrderWaitingStatus string      `json:"order_waiting_status"`
	SyncServiceStatus  string      `json:"sync_service_status"`
	LogisticsErr       string      `json:"logistics_err"` // 错误信息
	Labels             []Label     `json:"labels"`        // 面单信息
	MergeLabel         string      `json:"merge_label"`   // 合并面单 URL
	Fee                []Fee       `json:"fee"`           // 总费用信息
	FeeDetail          []FeeDetail `json:"fee_detail"`    // 参考号
}

// Label 面单
type Label struct {
	TrackingNumber  string      `json:"tracking_number"`  // 物流单号
	TrackingNumber2 null.String `json:"tracking_number2"` // 物流单号（UPS MI时产品返回USPS单号，其余产品不返回）
	LabelUrl        string      `json:"label_url"`        // 面单链接
	FileType        string      `json:"file_type"`        // 面单类型
	Context         any         `json:"context"`          // 上下文？
}
