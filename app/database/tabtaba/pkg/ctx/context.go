package ctx

import "context"

const TenantKey = "tenant"
const EnvKey = "env"

func NewCtx(tenant, env string) context.Context {
	ctx := context.TODO()
	ctx = withTenant(ctx, tenant)
	ctx = withEnv(ctx, env)
	return ctx
}

func withTenant(ctx context.Context, tenant string) context.Context {
	return context.WithValue(ctx, TenantKey, tenant)
}

func withEnv(ctx context.Context, env string) context.Context {
	return context.WithValue(ctx, EnvKey, env)
}

func GetTenantId(ctx context.Context) (reply int32) {
	reply = 1
	tenant := ctx.Value(TenantKey)
	if tenant == nil {
		return
	}
	if tenant.(string) == "bbx" {
		return
	}
	if tenant.(string) == "bbz" {
		return 2
	}
	return
}

func GetEnv(ctx context.Context) string {
	return ctx.Value(EnvKey).(string)
}

func GetWorkerAccountId(ctx context.Context) (reply int32) {
	reply = 1990
	tenantId := GetTenantId(ctx)
	switch tenantId {
	case 1:
		return
	case 2:
		reply = 3410
	}
	return
}
