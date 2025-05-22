package main

import (
	"context"
	"go-study/app/database/account/internal/biz"
	"go-study/app/database/account/internal/data"
	"log"
)

func main() {
	client := data.NewEntClient() // 使用 data 包创建客户端

	ctx := context.Background()

	// 执行迁移
	if err := biz.MigrateAccountToWallet(ctx, client); err != nil {
		log.Fatalf("数据迁移失败: %v", err)
	}

	log.Println("数据迁移成功")
}
