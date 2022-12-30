package plan

import "github.com/kevinjqiu/external-dns-docker/endpoint"

type Plan struct {
	Current []*endpoint.Endpoint
	Desired []*endpoint.Endpoint
	Missing []*endpoint.Endpoint
	Changes *Changes
}

type Changes struct {
	Create []*endpoint.Endpoint
	Update []*endpoint.Endpoint
	Delete []*endpoint.Endpoint
}
