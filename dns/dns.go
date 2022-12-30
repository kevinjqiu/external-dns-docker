package dns

import (
	"context"
	"github.com/kevinjqiu/external-dns-docker/endpoint"
)

type Provider interface {
	Records(ctx context.Context) ([]*endpoint.DNSEndpoint, error)
}
