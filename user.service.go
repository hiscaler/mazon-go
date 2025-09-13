package mazon

import (
	"context"

	"github.com/hiscaler/mazon-go/entity"
)

// 用户信息服务
type userService service

// Information 获取用户信息
func (s userService) Information(ctx context.Context) (info entity.UserInfo, err error) {
	res := struct {
		NormalResponse
		Result entity.UserInfo `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetResult(&res).
		Post("/getUserInfo")
	if err = recheckError(resp, res.NormalResponse, err); err != nil {
		return
	}
	return res.Result, nil
}
