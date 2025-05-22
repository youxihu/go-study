package biz

import (
	"context"
	"go-study/app/database/tabtaba/internal/data"
	"go-study/app/database/tabtaba/internal/data/ent"
	"go-study/app/database/tabtaba/internal/entity"
	ctx2 "go-study/app/database/tabtaba/pkg/ctx"
)

type FinanceUseCase struct {
	repo data.FinanceRepo
	env  string
}

func NewFinanceUseCase(env string) *FinanceUseCase {
	repo := data.NewFinanceRepo(env)
	return &FinanceUseCase{repo: repo, env: env}
}
func (uc *FinanceUseCase) ImportFinance(ctx context.Context, req []*entity.Project) ([]*ent.Finance, error) {

	// 提取常量
	var (
		defaultTransactNo       = "000000000000000000000000"
		defaultOrderNo          = "000000000000000000000000"
		defaultPayModeType int8 = 3
		defaultPayInfo          = ""
		defaultChangeType  int8 = 1
	)

	// 初始化切片容量
	finances := make([]*ent.Finance, 0, len(req))
	// 组装数据
	for _, v := range req {
		finances = append(finances, &ent.Finance{
			TenantID:       ctx2.GetTenantId(ctx),
			ProjectID:      v.ProjectId,
			TransactNo:     defaultTransactNo,
			OrderNo:        defaultOrderNo,
			AccountID:      v.AccountId,
			TheirAccountID: v.WorkerAccountId,
			Summary:        v.Summary,
			PayModeType:    defaultPayModeType,
			PayInfo:        defaultPayInfo,
			TransactAt:     v.TransactAt,
			Remark:         v.Remark,
			Amount:         v.Amount,
			ChangeType:     defaultChangeType,
			CreatedAt:      v.CreateTime,
			UpdatedAt:      v.CreateTime,
		})
	}
	return uc.repo.BulkCreateFinance(ctx, finances)
}
