package service

import (
	"context"
	"errors"
	"go-study/app/database/tabtaba/internal/biz"
	"go-study/app/database/tabtaba/internal/entity"
)

type Service struct {
	financeUc *biz.FinanceUseCase
	projectUc *biz.ProjectUseCase
}

func NewService(env string) *Service {
	financeUc := biz.NewFinanceUseCase(env)
	projectUc := biz.NewProjectUseCase(env)
	return &Service{
		financeUc: financeUc,
		projectUc: projectUc,
	}
}

func (s *Service) ExportOfflineProject(ctx context.Context, req []*entity.ExcelParseData) error {
	// 调用 ImportProject 方法
	projects, err := s.projectUc.ImportProject(ctx, req)
	if err != nil {
		return err
	}
	if len(projects) == 0 {
		return errors.New("can not import project")
	}
	if _, err := s.financeUc.ImportFinance(ctx, projects); err != nil {
		return err
	}
	return nil
}
