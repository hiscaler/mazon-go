package entity

// LogisticsLabel 物流面单信息
type LogisticsLabel struct {
	ReferenceNo string `json:"reference_no"` // 参考号
	OrderCode   string `json:"order_code"`   // 订单号, 我司系统唯一单号
	MergeLabel  string `json:"merge_label"`  // 订单下所有物流单号的合并面单地址
	Labels      []struct {
		TrackingNumber string `json:"tracking_number"` // 物流单号
		LabelUrl       string `json:"label_url"`       // 面单地址
		FileType       string `json:"file_type"`       // 面单类型
	} `json:"labels"` // 订单物流信息明细
}
