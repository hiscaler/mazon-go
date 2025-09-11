package entity

// Fee 费用信息
type Fee struct {
	FtCode       string `json:"ft_code"`       // 费用英文名称
	Amount       string `json:"amount"`        // 金额
	CurrencyCode string `json:"currency_code"` // 币种
	FtName       string `json:"ft_name"`       // 费用中文名称
}

type FeeDetail struct {
	Fee
	TrackingNumber string `json:"tracking_number"` // 物流单号（为空时代表没有生成物流单号，需要异步获取）
	BoxCode        string `json:"box_code"`        // 箱号
}
