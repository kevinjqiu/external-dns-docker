package dns

type Record struct {
	Name  string
	Value string
	Type  string
	TTL   int64
}
