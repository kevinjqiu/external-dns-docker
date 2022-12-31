package controller

import (
	"context"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/kevinjqiu/external-dns-docker/dns"
	"github.com/sirupsen/logrus"
	"time"
)

var ttl int64 = 600 // TODO: get that from container label

var logger = logrus.WithField("component", "controller")

func sanitizeHostName(hostname string) string {
	var ret = hostname

	if ret[0] == '/' {
		ret = ret[1:]
	}

	return ret
}

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

	for _, container := range containers {
		for _, name := range container.Names {
			for _, network := range container.NetworkSettings.Networks {
				record, err := s.dnsProvider.NewRecord(context.Background(), sanitizeHostName(name), "A", network.IPAddress, ttl)
				if err != nil {
					logger.WithError(err).Warnf("unable to create record: %v", err)
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

func (s *Controller) refresh() error {
	plan, err := s.generatePlan()
	if err != nil {
		return err
	}

	if err := s.dnsProvider.ApplyPlan(context.Background(), plan); err != nil {
		logger.WithError(err).Warnf("unable to apply changes: %v", err)
	}

	return nil
}

func (s *Controller) Run() error {
	if err := s.refresh(); err != nil {
		logger.WithError(err).Warnf("unable to refresh domain names: %v", err)
	}

	messageChan, errChan := s.dockerClient.Events(context.Background(), types.EventsOptions{
		Since: "",
		Until: "",
		Filters: filters.NewArgs(
			filters.KeyValuePair{Key: "event", Value: "start"},
			filters.KeyValuePair{Key: "event", Value: "stop"},
			filters.KeyValuePair{Key: "event", Value: "die"},
			filters.KeyValuePair{Key: "event", Value: "destroy"},
			filters.KeyValuePair{Key: "event", Value: "oom"},
			filters.KeyValuePair{Key: "event", Value: "rename"},
		),
	})

	ticker := time.NewTicker(30 * time.Second)

	for {
		select {
		case message := <-messageChan:
			// TODO: update a single DNS name when containers start/stop
			logger.Info("Got message: %v", message)

		case <-ticker.C:
			if err := s.refresh(); err != nil {
				logger.WithError(err).Warnf("unable to refresh domain names: %v", err)
			}

		case err := <-errChan:
			return err
		}
	}
}

func NewController(dockerClient *client.Client, dnsProvider dns.Provider) *Controller {
	return &Controller{
		dockerClient:  dockerClient,
		dnsProvider:   dnsProvider,
		labelEnabled:  "external-dns-docker/enabled",
		labelHostname: "external-dns-docker/hostname",
	}
}
