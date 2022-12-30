package svc

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type ExternalDNSService struct {
	dockerClient  *client.Client
	labelEnabled  string
	labelHostname string
	dnsProviders  []DNSProvider
}

func (s *ExternalDNSService) getContainers() ([]types.Container, error) {
	opts := types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.KeyValuePair{Key: "status", Value: "running"},
			filters.KeyValuePair{Key: "label", Value: s.labelEnabled},
		),
	}

	return s.dockerClient.ContainerList(context.Background(), opts)
}

func (s *ExternalDNSService) Run() {
	// Upon start, gather a list of eligible containers
	containers, err := s.getContainers()
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		fmt.Println(container)
	}
}

func NewExternalDNSService(dockerClient *client.Client, dnsProviders []DNSProvider) *ExternalDNSService {
	return &ExternalDNSService{
		dockerClient:  dockerClient,
		dnsProviders:  dnsProviders,
		labelEnabled:  "external-dns-docker/enabled",
		labelHostname: "external-dns-docker/hostname",
	}
}
