package entity

import "time"

type ExcelParseData struct {
	AccountId   int32   // 账户ID
	Amount      float64 // 金额
	CreateTime, // 创建时间
	ExpectDeliverTime time.Time // 预期交付时间
	Remark          string    //备注
	TransactAt      time.Time // 交易时间
	Title           string    //标题
	ProjectType     int8      //发布项目等级
	Summary         string    //摘要说明
	WorkerAccountId int32     // 接单人ID
}
