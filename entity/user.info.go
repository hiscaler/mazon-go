package entity

// UserInfo 用户信息
type UserInfo struct {
	Code    string   `json:"code"`
	Balance string   `json:"balance"`
	SmCode  []string `json:"sm_code"`
	Address []string `json:"address"`
}
