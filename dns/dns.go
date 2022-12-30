package dns

import (
	"context"
)

type Provider interface {
	Records(ctx context.Context) ([]*Record, error)
	ApplyPlan(ctx context.Context, plan *Plan) error
	NewRecord(ctx context.Context, baseName, recordType, value string, ttl int64) (*Record, error)
}
