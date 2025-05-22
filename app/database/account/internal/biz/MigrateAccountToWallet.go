package biz

import (
	"context"
	"go-study/app/database/account/internal/data/ent"
	"log"
)

// MigrateAccountToWallet 将旧表 Account 数据迁移到新表 AccountWallet
func MigrateAccountToWallet(ctx context.Context, client *ent.Client) error {
	// 查询旧表 Account 数据
	oldData, err := client.Account.Query().All(ctx)
	if err != nil {
		log.Printf("查询失败: %v", err)
		return err
	}

	log.Printf("待迁移记录数: %d", len(oldData))
	if len(oldData) == 0 {
		log.Println("没有需要迁移的数据")
		return nil
	}

	// 设置批次大小
	batchSize := 200                                           // 可根据数据库负载调整批次大小
	totalBatches := (len(oldData) + batchSize - 1) / batchSize // 计算总批次数

	// 分批处理插入
	for batch := 0; batch < totalBatches; batch++ {
		// 计算当前批次的起始和结束索引
		start := batch * batchSize
		end := start + batchSize
		if end > len(oldData) {
			end = len(oldData)
		}

		// 生成当前批次的数据
		bulk := make([]*ent.AccountWalletCreate, end-start)
		for i, data := range oldData[start:end] {
			// 获取用户身份字段 user_identity
			userIdentity := data.UserIdentity // 直接从数据对象中获取 user_identity 字段

			openInstalment := 2

			if userIdentity == 3 {
				openInstalment = 1
			}

			// 构造批量插入数据
			bulk[i] = client.AccountWallet.Create().
				SetID(data.ID).
				SetAccountID(data.ID).
				SetTenantID(data.TenantID).
				SetAccountIntegral(data.AccountIntegral).
				SetAccountBalance(data.AccountBalance).
				SetFreezeBalance(data.FreezeBalance).
				SetFreezeIntegral(data.FreezeIntegral).
				SetCreatedAt(data.CreatedAt).
				SetUpdatedAt(data.UpdatedAt).
				SetOpenInstallment(int8(openInstalment)).
				SetStatus(1).
				SetState(1)
		}

		// 执行批量插入操作
		err = client.AccountWallet.CreateBulk(bulk...). // 执行批量操作
			OnConflict(). // 处理冲突
			UpdateNewValues(). // 更新新值
			Exec(ctx) // 使用 Exec 执行操作
		if err != nil {
			log.Printf("第 %d 批次迁移失败: %v", batch+1, err)
			return err
		}

		log.Printf("第 %d 批次迁移完成，迁移记录数：%d", batch+1, end-start)
	}

	log.Println("数据迁移完成")
	return nil
}
