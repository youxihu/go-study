package data

import (
	"entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql"
	"go-study/app/database/account/internal/data/ent"
	"time"
)

// NewEntClient 创建数据库ent客户端
func NewEntClient() *ent.Client {

	drv, err := sql.Open(
		"mysql",
		"root:P@ssword1@tcp(192.168.0.214:3306)/bbx_account?charset=utf8mb4&parseTime=True&loc=Local",
	)
	if err != nil {
		panic(err)
	}
	drv.DB().SetMaxIdleConns(10)
	drv.DB().SetMaxOpenConns(100)
	drv.DB().SetConnMaxLifetime(time.Hour)
	client := ent.NewClient(ent.Driver(drv))
	return client
}
