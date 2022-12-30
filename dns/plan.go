package dns

type Plan struct {
	Current []*Record
	Desired []*Record
	Missing []*Record
	Changes *Changes
}

type Changes struct {
	Create []*Record
	Update []*Record
	Delete []*Record
}
