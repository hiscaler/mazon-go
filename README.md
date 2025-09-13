美正系统
=======

## 开发文档

[https://www.mazonlabel.com/docs/](https://www.mazonlabel.com/docs/)

## 安装

```shell
go get -u github.com/hiscaler/mazon-go
```

## 使用

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/hiscaler/mazon-go/config"
)

func main() {
	b, err := os.ReadFile("./config/config.json")
	if err != nil {
		panic(fmt.Sprintf("Read config error: %s", err.Error()))
	}
	var cfg config.Config
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		panic(fmt.Sprintf("Parse config file error: %s", err.Error()))
	}

	client := NewClient(ctx, cfg)
	err := client.getAccessToken(ctx)
	if err != nil {
		fmt.Println(err)
	}

	req := CreateOrderRequest{
		ReferenceNO:      "business-id",
		SMCode:           "USPS",
		OAFirstname:      "John",
		OACompany:        "Apple Inc.",
		OATelephone:      "12345678",
		OACountry:        "US",
		OAState:          "CA",
		OACity:           "Ontario",
		OAPostcode:       "91761",
		OAStreetAddress1: "2025 D Francis Street",
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
		Remark:      "Test it!",
	}
	res, err := client.Services.Order.Create(context.Context(), req)
}

```

## Services
 - Order
   - Order.Create 创建订单 
   - Order.Query 根据查询条件筛选符合条件的订单列表数据
   - Order.Cancel 取消订单
 - Rate
   - Calc 运费计算 
 - ScanForm
   - ScanForm.Create 基于多个跟踪号生成 ScanForm 
 - ShippingLabel
   - ShippingLabel.Detail 获取面单
   - ShippingLabel.Query 根据物流单号获取面单信息
 - User
   - Information 获取用户信息