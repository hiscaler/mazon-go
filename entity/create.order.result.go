package entity

type CreateOrderResult struct {
	OrderCode string `json:"order_code"` // 订单号，后续的操作需要以此订单为依据条件
	Fee       []struct {
		FtCode       string `json:"ft_code"`       // 费用英文名称
		Amount       string `json:"amount"`        // 金额
		CurrencyCode string `json:"currency_code"` // 币种
		FtName       string `json:"ft_name"`       // 费用中文名称
	} `json:"fee"` // 费用信息
	FeeDetail []struct {
		FtCode         string `json:"ft_code"`         // 费用英文名称
		Amount         string `json:"amount"`          // 金额
		TrackingNumber string `json:"tracking_number"` // 物流单号（为空时代表没有生成物流单号，需要异步获取）
		BoxCode        string `json:"box_code"`        // 箱号
		CurrencyCode   string `json:"currency_code"`   // 币种
		FtName         string `json:"ft_name"`         // 费用中文名称
	} `json:"fee_detail"` // 费用详情
	Labels struct {
		TrackingNumber  string `json:"tracking_number"`  // 物流单号
		TrackingNumber2 string `json:"tracking_number2"` // 物流单号（UPS MI时产品返回USPS单号，其余产品不返回）
		LabelUrl        string `json:"label_url"`        // 面单URL链接
		FileType        string `json:"file_type"`        // 面单类型
	} `json:"labels"`                      // 为空时代表没有生产物流信息，需要异步获取
	MergeLabel string `json:"merge_label"` //合并面单
}
