package entity

import "time"

type Project struct {
	AccountId         int32     // 账户ID
	ProjectId         int32     // 项目ID
	Amount            float64   // 金额
	CreateTime        time.Time // 创建时间
	ExpectDeliverTime time.Time // 预期交付时间
	Remark            string    // 备注
	TransactAt        time.Time // 交易时间
	Title             string    // 标题
	WorkerAccountId   int32
	Summary           string
}
