package biz

import (
	"context"
	"go-study/app/database/tabtaba/internal/data"
	"go-study/app/database/tabtaba/internal/data/ent"
	"go-study/app/database/tabtaba/internal/entity"
	ctx2 "go-study/app/database/tabtaba/pkg/ctx"
	"time"
)

type ProjectUseCase struct {
	repo data.ProjectRepo
	env  string
}

func NewProjectUseCase(env string) *ProjectUseCase {
	repo := data.NewProjectRepo(env)
	return &ProjectUseCase{repo: repo, env: env}
}

func (uc *ProjectUseCase) ImportProject(ctx context.Context, excelData []*entity.ExcelParseData) ([]*entity.Project, error) {
	// 提取常量
	var (
		defaultAccountPhone            = "13506536336"
		defaultProjectStatus     int32 = 200
		defaultFiles                   = "[]"
		defaultActionPlat        int8  = 2
		defaultAction                  = "完成验收"
		defaultProjectStatusText       = "已完成"
	)
	// 初始化切片容量
	projects := make([]*ent.Project, 0, len(excelData))
	// 组装数据
	for _, v := range excelData {
		projects = append(projects, &ent.Project{
			TenantID:          ctx2.GetTenantId(ctx),
			AccountID:         v.AccountId,
			AccountPhone:      defaultAccountPhone,
			WorkerAccountID:   v.WorkerAccountId,
			ProjectType:       v.ProjectType,
			Title:             v.Title,
			ExpectDeliverTime: v.ExpectDeliverTime,
			TenderFiles:       defaultFiles,
			DeliverBidFiles:   defaultFiles,
			DeliverBidTime:    v.ExpectDeliverTime,
			Remark:            v.Remark,
			ProjectStatus:     defaultProjectStatus,
			Amount:            v.Amount,
			ESignNotify:       defaultFiles,
			CreatedAt:         v.CreateTime,
			PayMoney:          v.Amount,
			UpdatedAt:         v.CreateTime,
		})
	}
	// 批量插入项目
	insertedProject, err := uc.repo.BulkCreateProject(ctx, projects)
	if err != nil {
		return nil, err
	}
	projectActions := make([]*ent.ProjectActionRecord, 0, len(insertedProject))
	for _, v := range insertedProject {
		projectActions = append(projectActions, &ent.ProjectActionRecord{
			TenantID:          ctx2.GetTenantId(ctx),
			ProjectID:         v.ID,
			ActionPlat:        defaultActionPlat,
			Action:            defaultAction,
			AccountID:         v.AccountID,
			ProjectStatus:     defaultProjectStatus,
			ProjectStatusText: defaultProjectStatusText,
			CreatedAt:         v.CreatedAt,
			UpdatedAt:         v.CreatedAt,
		})
	}
	// 批量插入动作
	if _, err := uc.repo.BulkCreateProjectAction(ctx, projectActions); err != nil {
		return nil, err
	}
	// 返回数据
	var reply = make([]*entity.Project, 0)
	for _, v := range insertedProject {
		var (
			transactAt time.Time
			summary    string
		)
		for _, parseData := range excelData {
			if parseData.Title == v.Title {
				transactAt = parseData.TransactAt
				summary = parseData.Summary
			}
		}
		reply = append(reply, &entity.Project{
			AccountId:         v.AccountID,
			ProjectId:         v.ID,
			Amount:            v.Amount,
			CreateTime:        v.CreatedAt,
			ExpectDeliverTime: v.DeliverBidTime,
			Remark:            v.Remark,
			TransactAt:        transactAt,
			Title:             v.Title,
			WorkerAccountId:   v.WorkerAccountID,
			Summary:           summary,
		})
	}
	return reply, nil
}
