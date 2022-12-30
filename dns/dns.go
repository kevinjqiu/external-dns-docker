package dns

import (
	"context"
	"github.com/kevinjqiu/external-dns-docker/endpoint"
	"github.com/kevinjqiu/external-dns-docker/plan"
)

type Provider interface {
	Records(ctx context.Context) ([]*endpoint.Endpoint, error)
	ApplyPlan(ctx context.Context, plan plan.Plan) error
}
