package mazon

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_rateService_Calc(t *testing.T) {
	req := RateCalcRequest{
		ReferenceNO:      "TEST-FROM-SDK",
		SMCode:           "USPS GA13",
		OAFirstname:      "ZEB2",
		OACompany:        "SILBER BLITZ",
		OATelephone:      "0731-12345678",
		OACountry:        "US",
		OAState:          "CA",
		OACity:           "Ontario",
		OAPostcode:       "91761",
		OAStreetAddress1: "2078 E Francis Street",
		BoxList: []RateCalcOrderBox{
			{
				Height:       1,
				Length:       1,
				Width:        1,
				ActualWeight: 1,
			},
		},
		IsMoreBox:   1,
		ShipperCode: "S0004",
		Remark:      "测试订单，请不要发货",
	}
	res, _ := client.Services.Rate.Calc(ctx, req)
	//assert.Nil(t, err)
	//assert.Equal(t, req.SMCode, res.SmCode)
	fmt.Println(fmt.Sprintf("%#v", res))
	b, _ := json.Marshal(&res)
	fmt.Println(string(b))
	assert.Equal(t, res.AddressType, 2)
}
