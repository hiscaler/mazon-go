package entity

type RateCalcResult struct {
	SmCode          string `json:"sm_code"`           // 物流产品
	AddressTypeText string `json:"address_type_text"` // 地址类型描述
	AddressType     string `json:"address_type"`      // 地址类型
	CurrencyCode    string `json:"currency_code"`     // 币种
	ShippingCharge  string `json:"shipping_charge"`   // 基础运费
	TotalCharge     string `json:"total_charge"`      // 总金额
	ChargeDetail    []struct {
		FtCode      string `json:"ft_code"`       // 费用英文名称
		FeeTypeCode string `json:"fee_type_code"` // 费用编码
		ChargeDesc  string `json:"charge_desc"`   // 费用描述
		Amount      string `json:"amount"`        // 费用描述
	} `json:"charge_detail"` // 费用明细
}
