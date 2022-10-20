package config

type Config struct {
	Debug    bool   // 是否启用调试模式
	Timeout  int    // HTTP 超时设定（单位：秒）
	AppKey   string // App Key
	AppToken string // App Token
}
