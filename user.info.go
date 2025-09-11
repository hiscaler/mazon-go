package areship

import (
	"context"

	"github.com/hiscaler/areship-go/entity"
)

// 用户信息服务
type userInfoService service

// Detail 获取用户信息
func (s userInfoService) Detail(ctx context.Context) (info entity.UserInfo, err error) {
	res := struct {
		NormalResponse
		Result entity.UserInfo `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		Post("/getUserInfo")
	if err = recheckError(resp, res.NormalResponse, err); err != nil {
		return
	}
	return res.Result, nil
}
