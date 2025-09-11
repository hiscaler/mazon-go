package entity

type Token struct {
	AccessToken string `json:"access_token"` // 授权 Token
	UserInfo    struct {
		ID           int    `json:"u_id"`            // 用户 ID
		Account      string `json:"u_account"`       // 用户账号
		CustomerCode string `json:"u_customer_code"` // 用户客户代码
	} `json:"user_info"` // 用户信息
}
