package endpoint

type Endpoint struct {
	DNSName    string
	Targets    string
	RecordType string
	RecordTTL  int64
}
