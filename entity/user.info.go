package entity

// UserInfo 用户信息
type UserInfo struct {
	Code    string           `json:"code"`
	Balance string           `json:"balance"`
	SmCode  []string         `json:"sm_code"`
	Address []ShipperAddress `json:"address"`
}

// ShipperAddress 发件人信息, 发件人信息必须与我司备案信息完全一致, 发件人信息与发件人编码同时存在时，以发件人信息为准
type ShipperAddress struct {
	ShipperCode          string `json:"shipper_code"`               // 发件人编码, 发件人信息与发件人编码同时存在时，以发件人信息为准
	ShipperName          string `json:"shipper_name"`               // 发件人姓名,长度最短3位数,最长35位数
	ShipperCompany       string `json:"shipper_company,omitempty"`  // 发件人公司， 长度最长35位数
	ShipperTelPhone      string `json:"shipper_tel_phone"`          // 发件人电话， 10-15位之间
	ShipperCountry       string `json:"shipper_country"`            // 发件人国家， 要求固定值US
	ShipperStateProvince string `json:"shipper_state_province"`     // 发件人州, 只能为大写二字编码
	ShipperCity          string `json:"shipper_city"`               // 发件人城市, 长度最短1位数，最长30位数
	ShipperPostalCode    string `json:"shipper_postal_code"`        // 发件人邮编, 长度最短5位数,最长10位数
	ShipperAddress1      string `json:"shipper_address1"`           // 发件人地址1,长度不得超过35位
	ShipperAddress2      string `json:"shipper_address2,omitempty"` // 发件人地址2,长度不得超过35位
}

// ReturnAddress 退件信息
type ReturnAddress struct {
	StreetAddress    string `json:"street_address"`      // 街道地址
	SecondaryAddress string `json:"secondary_address"`   // 街道地址2
	ZipCodeAndPlus4  string `json:"zip_code_and_plus4"`  // 退件信息-邮政编码-扩展邮编,如果有扩展邮编则用 - 拼接在后面，比如：75115-2500
	City             string `json:"city"`                // 城市
	State            string `json:"state"`               // 州
	FirstName        string `json:"first_name"`          // 联系人名字
	LastName         string `json:"last_name,omitempty"` // 联系人姓氏
	Phone            string `json:"phone,omitempty"`     // 联系人电话
}
