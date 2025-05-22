package data

import (
	"context"
	_ "entgo.io/ent/dialect/sql"
	_ "github.com/go-sql-driver/mysql" // 导入 MySQL 驱动
	"go-study/app/database/tabtaba/internal/data/ent"
)

const (
	DevProjectDsn  = "root:P@ssword1@tcp(192.168.0.214:3306)/project"
	DevFianceDsn   = "root:P@ssword1@tcp(192.168.0.214:3306)/finance"
	ProdProjectDsn = "root:root@Pswword@tcp(36.27.222.248:3306)/project"
	ProdFianceDsn  = "root:root@Pswword@tcp(36.27.222.248:3306)/finance"
	FatProjectDsn  = "root:root@Pswword@tcp(116.211.128.111:3306)/project"
	FatFianceDsn   = "root:root@Pswword@tcp(116.211.128.11:3306)/finance"
	UatProjectDsn  = "root:root@Pswword@tcp(115.223.6.206:3306)/project"
	UatFianceDsn   = "root:root@Pswword@tcp(115.223.6.206:3306)/finance"
)

type DB struct {
	pdb *ent.Client
	fdb *ent.Client
}

func (db *DB) GetProjectClient() *ent.Client {
	return db.pdb
}

func (db *DB) GetFinanceClient() *ent.Client {
	return db.fdb
}

// BeginProject 开启 Project 数据库的事务
func (db *DB) BeginProject(ctx context.Context) (*ent.Tx, error) {
	return db.pdb.Tx(ctx)
}

// BeginFinance 开启 Finance 数据库的事务
func (db *DB) BeginFinance(ctx context.Context) (*ent.Tx, error) {
	return db.fdb.Tx(ctx)
}
