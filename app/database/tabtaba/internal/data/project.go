package data

import (
	"context"
	"go-study/app/database/tabtaba/internal/data/ent"
	"log"
)

var _ ProjectRepo = (*projectRepo)(nil)

type projectRepo struct {
	db *ent.Client
}

func NewProjectRepo(env string) ProjectRepo {
	var dsn string
	switch env {
	case "dev":
		dsn = DevProjectDsn
	case "prod":
		dsn = ProdProjectDsn
	case "public":
		dsn = FatProjectDsn
	case "rc":
		dsn = UatProjectDsn
	default:
		dsn = DevProjectDsn
	}
	pdb, err := ent.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("failed to connect to finance database: %v", err)
	}
	return &projectRepo{
		db: pdb,
	}
}

func (r *projectRepo) BulkCreateProject(ctx context.Context, req []*ent.Project) ([]*ent.Project, error) {
	bulk := make([]*ent.ProjectCreate, len(req))
	for i, project := range req {
		bulk[i] = r.db.Project.Create().
			SetAccountID(project.AccountID).
			SetAccountPhone(project.AccountPhone).
			SetWorkerAccountID(project.WorkerAccountID).
			SetTenantID(project.TenantID).
			SetProjectType(project.ProjectType).
			SetTitle(project.Title).
			SetRemark(project.Remark).
			SetAmount(project.Amount).
			SetPayMoney(project.PayMoney).
			SetProjectStatus(project.ProjectStatus).
			SetCreatedAt(project.CreatedAt).
			SetUpdatedAt(project.UpdatedAt).
			SetDeliverBidTime(project.DeliverBidTime).
			SetExpectDeliverTime(project.ExpectDeliverTime).
			SetTenderFiles(project.TenderFiles).
			SetDeliverBidFiles(project.DeliverBidFiles).
			SetESignNotify(project.ESignNotify)
	}
	return r.db.Project.CreateBulk(bulk...).Save(ctx)
}

func (r *projectRepo) BulkCreateProjectAction(ctx context.Context, req []*ent.ProjectActionRecord) ([]*ent.ProjectActionRecord, error) {
	bulk := make([]*ent.ProjectActionRecordCreate, len(req))
	for i, projectAction := range req {
		bulk[i] = r.db.ProjectActionRecord.Create().
			SetProjectID(projectAction.ProjectID).
			SetTenantID(projectAction.TenantID).
			SetActionPlat(projectAction.ActionPlat).
			SetAction(projectAction.Action).
			SetAccountID(projectAction.AccountID).
			SetProjectStatus(projectAction.ProjectStatus).
			SetProjectStatusText(projectAction.ProjectStatusText).
			SetCreatedAt(projectAction.CreatedAt).
			SetUpdatedAt(projectAction.UpdatedAt)
	}
	return r.db.ProjectActionRecord.CreateBulk(bulk...).Save(ctx)
}
