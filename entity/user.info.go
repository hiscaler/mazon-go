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
	ShipperCode          string `json:"shipper_code"`
	ShipperName          string `json:"shipper_name"`
	ShipperCompany       string `json:"shipper_company,omitempty"`
	ShipperAddress1      string `json:"shipper_address1"`
	ShipperAddress2      string `json:"shipper_address2,omitempty"`
	ShipperCountry       string `json:"shipper_country"`
	ShipperStateProvince string `json:"shipper_state_province"`
	ShipperCity          string `json:"shipper_city"`
	ShipperPostalCode    string `json:"shipper_postal_code"`
	ShipperTelPhone      string `json:"shipper_tel_phone"`
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
