package entity

type OrderCreateResult struct {
	OrderCode   string      `json:"order_code"`   // 订单号，后续的操作需要以此订单为依据条件
	LabelStatus int         `json:"label_status"` // 0:预报失败 1：预报中，还没有返回预报结果 2：预报成功
	Fee         []Fee       `json:"fee"`          // 费用信息
	FeeDetail   []FeeDetail `json:"fee_detail"`   // 费用详情
	Labels      []Label     `json:"labels"`       // 为空时代表没有生产物流信息，需要异步获取
	MergeLabel  string      `json:"merge_label"`  // 合并面单
}
