package data

import (
	"context"
	"go-study/app/database/tabtaba/internal/data/ent"
	"log"
)

var _ FinanceRepo = (*financeRepo)(nil)

type financeRepo struct {
	db *ent.Client
}

func NewFinanceRepo(env string) FinanceRepo {
	var dsn string
	switch env {
	case "dev":
		dsn = DevFianceDsn
	case "prod":
		dsn = ProdFianceDsn
	case "public":
		dsn = FatFianceDsn
	case "rc":
		dsn = UatFianceDsn
	default:
		dsn = DevFianceDsn
	}
	fdb, err := ent.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to finance database: %v", err)
	}
	return &financeRepo{
		db: fdb,
	}
}

func (r *financeRepo) BulkCreateFinance(ctx context.Context, req []*ent.Finance) ([]*ent.Finance, error) {
	bulk := make([]*ent.FinanceCreate, len(req))
	for i, finance := range req {
		bulk[i] = r.db.Finance.Create().
			SetTenantID(finance.TenantID).
			SetProjectID(finance.ProjectID).
			SetTenantID(finance.TenantID).
			SetTransactNo(finance.TransactNo).
			SetOrderNo(finance.OrderNo).
			SetAccountID(finance.AccountID).
			SetTheirAccountID(finance.TheirAccountID).
			SetSummary(finance.Summary).
			SetPayModeType(finance.PayModeType).
			SetPayInfo(finance.PayInfo).
			SetTransactAt(finance.TransactAt).
			SetRemark(finance.Remark).
			SetAmount(finance.Amount).
			SetChangeType(finance.ChangeType).
			SetCreatedAt(finance.CreatedAt).
			SetUpdatedAt(finance.UpdatedAt)
	}
	return r.db.Finance.CreateBulk(bulk...).Save(ctx)
}
