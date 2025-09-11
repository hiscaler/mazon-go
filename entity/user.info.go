package entity

// UserInfo 用户信息
type UserInfo struct {
	Code    string    `json:"code"`
	Balance string    `json:"balance"`
	SmCode  []string  `json:"sm_code"`
	Address []Address `json:"address"`
}

type Address struct {
	ShipperCode          string `json:"shipper_code"`
	ShipperName          string `json:"shipper_name"`
	ShipperCompany       string `json:"shipper_company"`
	ShipperAddress1      string `json:"shipper_address1"`
	ShipperAddress2      string `json:"shipper_address2"`
	ShipperStateProvince string `json:"shipper_state_province"`
	ShipperCity          string `json:"shipper_city"`
	ShipperPostalCode    string `json:"shipper_postal_code"`
	ShipperTelphone      string `json:"shipper_telphone"`
}
