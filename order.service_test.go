package mazon

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_orderService_Create(t *testing.T) {
	req := CreateOrderRequest{
		ReferenceNO:      "TEST-FROM-SDK",
		SMCode:           "USPS GA13",
		OAFirstname:      "ZZZ",
		OACompany:        "SZZZ",
		OATelephone:      "0731-12345678",
		OACountry:        "US",
		OAState:          "CA",
		OACity:           "Ontario",
		OAPostcode:       "91761",
		OAStreetAddress1: "2078 E Francis Street",
		BoxList: []OrderBox{
			{
				Height:       1,
				Length:       1,
				Width:        1,
				ActualWeight: 1,
				Sku:          "MKG001",
				CnName:       "个性化定制马克杯",
				EngName:      "Personalized custom mugs",
			},
		},
		IsMoreBox:   1,
		ShipperCode: "S0004",
		Remark:      "测试订单，无需发货",
	}
	res, err := client.Services.Order.Create(ctx, req)
	assert.Nil(t, err)
	assert.NotEmpty(t, res.OrderCode)
}

func Test_orderService_Query(t *testing.T) {
	res, err := client.Services.Order.Query(ctx, OrderQueryRequest{
		OrderCode: "EPB00120250912114236000021",
	})
	assert.Nil(t, err)
	assert.NotEmpty(t, res)
}

func Test_orderService_Cancel(t *testing.T) {
	status, err := client.Services.Order.Cancel(ctx, CancelOrderRequest{
		OrderCode: "EPB00120250912114236000021",
	})
	assert.Nil(t, err)
	if err == nil {
		assert.Contains(t, status, []int{5, 6})
	}
}
