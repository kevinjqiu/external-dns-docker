package controller

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/kevinjqiu/external-dns-docker/dns"
	"log"
)

type Controller struct {
	dockerClient  *client.Client
	labelEnabled  string
	labelHostname string
	dnsProvider   dns.Provider
}

func (s *Controller) getEligibleContainers() ([]types.Container, error) {
	opts := types.ContainerListOptions{
		All: true,
		Filters: filters.NewArgs(
			filters.KeyValuePair{Key: "status", Value: "running"},
			filters.KeyValuePair{Key: "label", Value: s.labelEnabled},
		),
	}

	return s.dockerClient.ContainerList(context.Background(), opts)
}

func (s *Controller) getCurrentRecords() ([]*dns.Record, error) {
	records, err := s.dnsProvider.Records(context.Background())
	if err != nil {
		return nil, err
	}

	return records, nil
}

func (s *Controller) getDesiredRecords() ([]*dns.Record, error) {
	containers, err := s.getEligibleContainers()
	if err != nil {
		return nil, err
	}

	records := make([]*dns.Record, 0, len(containers))

	var ttl int64 = 600 // TODO: get that from container label

	for _, container := range containers {
		for _, name := range container.Names {
			for _, network := range container.NetworkSettings.Networks {
				record, err := s.dnsProvider.NewRecord(context.Background(), name, "A", network.IPAddress, ttl)
				if err != nil {
					log.Printf("unable to create record: %v", err)
					continue
				}
				records = append(records, record)
			}
		}
	}

	return records, nil
}

func (s *Controller) generatePlan() (*dns.Plan, error) {
	current, err := s.getCurrentRecords()
	if err != nil {
		return nil, err
	}

	desired, err := s.getDesiredRecords()
	if err != nil {
		return nil, err
	}

	return &dns.Plan{
		Current: current,
		Desired: desired,
	}, nil
}

func (s *Controller) Run() error {
	plan, err := s.generatePlan()
	if err != nil {
		return err
	}

	fmt.Printf("%v\n", plan.Desired)
	fmt.Printf("%v\n", plan.Current)
	// Upon start, gather a list of eligible containers
	// messageChan, errChan := cli.Events(context.Background(), types.EventsOptions{})

	// for {
	// 	select {
	// 	case message := <-messageChan:
	// 		fmt.Println(message)

	// 	case err := <-errChan:
	// 		panic(err)
	// 	}
	// }
	return nil
}

func NewController(dockerClient *client.Client, dnsProvider dns.Provider) *Controller {
	return &Controller{
		dockerClient:  dockerClient,
		dnsProvider:   dnsProvider,
		labelEnabled:  "external-dns-docker/enabled",
		labelHostname: "external-dns-docker/hostname",
	}
}
