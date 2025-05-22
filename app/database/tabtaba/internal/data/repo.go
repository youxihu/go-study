package data

import (
	"context"
	"go-study/app/database/tabtaba/internal/data/ent"
)

type FinanceRepo interface {
	BulkCreateFinance(ctx context.Context, finance []*ent.Finance) ([]*ent.Finance, error)
}
type ProjectRepo interface {
	BulkCreateProject(ctx context.Context, project []*ent.Project) ([]*ent.Project, error)
	BulkCreateProjectAction(ctx context.Context, req []*ent.ProjectActionRecord) ([]*ent.ProjectActionRecord, error)
}
