package mazon

import (
	"context"
	"errors"
	"strings"

	"github.com/hiscaler/mazon-go/entity"
)

type scanFormService service

// Create 基于多个跟踪号生成 ScanForm， 跟踪号必须同一个发货地址才可生成
// https://www.mazonlabel.com/docs/orderapi/%E7%94%9F%E6%88%90ScanForm%E5%8D%95%E6%8D%AE.html
func (s scanFormService) Create(ctx context.Context, trackingNumbers ...string) (forms []entity.ScanForm, err error) {
	numbers := make([]string, 0, len(trackingNumbers))
	for _, number := range trackingNumbers {
		number = strings.TrimSpace(number)
		if number != "" {
			numbers = append(numbers, number)
		}
	}
	if len(numbers) == 0 {
		return forms, errors.New("无效的跟踪号")
	}

	res := struct {
		NormalResponse
		Result []entity.ScanForm `json:"result"`
	}{}
	resp, err := s.httpClient.R().
		SetContext(ctx).
		SetBody(map[string]string{"tracking_number": strings.Join(numbers, ",")}).
		SetResult(&res).
		Post("/createScanForm")
	if err = recheckError(resp, res.NormalResponse, err); err != nil {
		return forms, err
	}
	return res.Result, nil
}
