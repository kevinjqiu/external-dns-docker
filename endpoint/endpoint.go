package endpoint

type DNSEndpoint struct {
	DNSName    string
	Targets    string
	RecordType string
	RecordTTL  int64
}
