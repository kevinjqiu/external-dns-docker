package cloudflare

import (
	"context"
	cloudflare_sdk "github.com/cloudflare/cloudflare-go"
	"github.com/kevinjqiu/external-dns-docker/dns"
	"os"
	"strings"
)

type CloudflareProvider struct {
	api      *cloudflare_sdk.API
	zoneName string
	suffix   string
}

func newDnsRecord(record cloudflare_sdk.DNSRecord) *dns.Record {
	return &dns.Record{
		Name:  record.Name,
		Value: record.Content,
		Type:  record.Type,
		TTL:   int64(record.TTL),
	}
}

func (c *CloudflareProvider) fullSuffix() string {
	var retval string
	if c.suffix == "" {
		retval = "." + c.zoneName
	} else {
		retval = "." + c.suffix + "." + c.zoneName
	}

	return retval
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
	return nil
}

func (c *CloudflareProvider) NewRecord(ctx context.Context, baseName, recordType, value string, ttl int64) (*dns.Record, error) {
	return &dns.Record{
		Name:  baseName + "." + c.fullSuffix(),
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
	return &CloudflareProvider{
		api:      api,
		zoneName: zoneName,
		suffix:   suffix,
	}, nil
}
