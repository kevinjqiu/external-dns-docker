package cloudflare

import (
	"context"
	"fmt"
	cloudflare_sdk "github.com/cloudflare/cloudflare-go"
	"github.com/kevinjqiu/external-dns-docker/dns"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

var logger = logrus.WithField("dns_provider", "cloudflare")

type CloudflareProvider struct {
	api      *cloudflare_sdk.API
	zoneName string
	zoneID   string
	suffix   string
}

func newDnsRecord(record cloudflare_sdk.DNSRecord) *dns.Record {
	return &dns.Record{
		Name:  record.Name,
		Value: record.Content,
		Type:  record.Type,
		TTL:   int64(record.TTL),
		ProviderMetadata: map[string]interface{}{
			"ID": record.ID,
		},
	}
}

func (c *CloudflareProvider) newCloudflareRecord(record *dns.Record) cloudflare_sdk.DNSRecord {
	return cloudflare_sdk.DNSRecord{
		Type:     record.Type,
		Name:     record.Name,
		Content:  record.Value,
		ZoneID:   c.zoneID,
		ZoneName: c.zoneName,
		TTL:      int(record.TTL),
	}
}

func (c *CloudflareProvider) fullSuffix() string {
	if c.suffix == "" {
		return fmt.Sprintf(".%s", c.zoneName)
	}
	return fmt.Sprintf(".%s.%s", c.suffix, c.zoneName)
}

func (c *CloudflareProvider) Records(ctx context.Context) ([]*dns.Record, error) {
	rr := cloudflare_sdk.DNSRecord{}

	zoneID, err := c.api.ZoneIDByName(c.zoneName)
	if err != nil {
		return nil, err
	}

	cfRecords, err := c.api.DNSRecords(ctx, zoneID, rr)
	if err != nil {
		return nil, err
	}

	records := make([]*dns.Record, 0, len(cfRecords))

	fullSuffix := c.fullSuffix()

	for _, record := range cfRecords {
		if strings.HasSuffix(record.Name, fullSuffix) {
			records = append(records, newDnsRecord(record))
		}
	}

	return records, nil
}

func (c *CloudflareProvider) ApplyPlan(ctx context.Context, plan *dns.Plan) error {
	changes := plan.Changes()

	if len(changes.Create) == 0 {
		logger.Info("No new records to add")
	} else {
		for _, record := range changes.Create {
			logger.Info("Adding record %v", record)
			response, err := c.api.CreateDNSRecord(ctx, c.zoneID, c.newCloudflareRecord(record))
			if err != nil {
				logger.WithError(err).Errorf("cannot create dns record %v: %v", record, err)
			}
			logger.Infof("%v", response)
		}
	}

	if len(changes.Delete) == 0 {
		logger.Info("No records to delete")
	} else {
		for _, record := range changes.Delete {
			logger.Info("Deleting record %v", record)
			recordID, ok := record.ProviderMetadata["ID"].(string)
			if !ok {
				logger.Errorf("unable to get the cloudflare record ID for %v", record)
				continue
			}

			err := c.api.DeleteDNSRecord(ctx, c.zoneID, recordID)
			if err != nil {
				logger.WithError(err).Errorf("cannot delete dns record %v: %v", record, err)
			}
		}
	}
	return nil
}

func (c *CloudflareProvider) NewRecord(ctx context.Context, baseName, recordType, value string, ttl int64) (*dns.Record, error) {
	return &dns.Record{
		Name:  fmt.Sprintf("%s%s", baseName, c.fullSuffix()),
		Value: value,
		Type:  recordType,
		TTL:   ttl,
	}, nil
}

func NewCloudflareProvider(zoneName string, suffix string) (*CloudflareProvider, error) {
	apiToken := os.Getenv("CLOUDFLARE_API_TOKEN")
	api, err := cloudflare_sdk.NewWithAPIToken(apiToken)

	if err != nil {
		return nil, err
	}

	zoneID, err := api.ZoneIDByName(zoneName)
	if err != nil {
		return nil, err
	}

	return &CloudflareProvider{
		api:      api,
		zoneName: zoneName,
		zoneID:   zoneID,
		suffix:   suffix,
	}, nil
}
