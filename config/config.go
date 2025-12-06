package config

type Config struct {
	Debug         bool   `json:"debug"`          // 是否启用调试模式
	Timeout       int    `json:"timeout"`        // HTTP 超时设定（单位：秒）
	AppKey        string `json:"app_key"`        // App Key
	AppToken      string `json:"app_token"`      // App Token
	TokenDuration int    `json:"token_duration"` // Token 生效时长（单位：小时）
}
