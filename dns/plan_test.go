package dns

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func record(name, value, recordType, id string) *Record {
	record := &Record{
		Name:  name,
		Value: value,
		Type:  recordType,
		ProviderMetadata: map[string]interface{}{
			"ID": id,
		},
	}

	if id != "" {
		record.ProviderMetadata = map[string]interface{}{
			"ID": id,
		}
	}
	return record
}

func A(name, value, id string) *Record {
	return record(name, value, "A", id)
}

func CNAME(name, value, id string) *Record {
	return record(name, value, "CNAME", id)
}

func requireNumberOfChanges(t *testing.T, changes *Changes, numCreate, numDelete int) {
	require.Equal(t, numCreate, len(changes.Create))
	require.Equal(t, numDelete, len(changes.Delete))
}

func TestPlanChanges_NewRecords(t *testing.T) {
	plan := Plan{
		Current: []*Record{
			A("foo.docker.mydomain.com", "192.168.100.100", "ID_1"),
		},
		Desired: []*Record{
			A("foo.docker.mydomain.com", "192.168.100.100", ""),
			A("bar.docker.mydomain.com", "192.168.100.101", ""),
			A("quux.docker.mydomain.com", "192.168.100.102", ""),
		},
	}

	changes := plan.Changes()

	requireNumberOfChanges(t, changes, 2, 0)

	require.Equal(t, "bar.docker.mydomain.com", changes.Create[0].Name)
	require.Equal(t, "quux.docker.mydomain.com", changes.Create[1].Name)
}

func TestPlanChanges_MultipleRecords(t *testing.T) {
	plan := Plan{
		Current: []*Record{
			A("foo.docker.mydomain.com", "192.168.100.100", "ID_1"),
		},
		Desired: []*Record{
			A("foo.docker.mydomain.com", "192.168.100.100", ""),
			A("foo.docker.mydomain.com", "192.168.100.101", ""),
		},
	}

	changes := plan.Changes()

	requireNumberOfChanges(t, changes, 1, 0)

	require.Equal(t, "foo.docker.mydomain.com", changes.Create[0].Name)
}

func TestPlanChanges_MultipleRecordsWithDifferentTypes(t *testing.T) {
	plan := Plan{
		Current: []*Record{
			A("foo.docker.mydomain.com", "192.168.100.100", "ID_1"),
			CNAME("foo.docker.mydomain.com", "woohoo.google.com", "ID_2"),
		},
		Desired: []*Record{
			A("foo.docker.mydomain.com", "192.168.100.100", ""),
			A("foo.docker.mydomain.com", "192.168.100.101", ""),
			CNAME("foo.docker.mydomain.com", "woohoo.google.com", ""),
		},
	}

	changes := plan.Changes()

	requireNumberOfChanges(t, changes, 1, 0)
	require.Equal(t, "foo.docker.mydomain.com", changes.Create[0].Name)
}

func TestPlanChanges_DeleteRecord(t *testing.T) {
	plan := Plan{
		Current: []*Record{
			A("foo.docker.mydomain.com", "192.168.100.100", "ID_1"),
			A("bar.docker.mydomain.com", "10.100.100.110", "ID_2"),
		},
		Desired: []*Record{
			A("bar.docker.mydomain.com", "10.100.100.110", ""),
		},
	}

	changes := plan.Changes()

	requireNumberOfChanges(t, changes, 0, 1)
	require.Equal(t, "foo.docker.mydomain.com", changes.Delete[0].Name)
	require.Equal(t, "ID_1", changes.Delete[0].ProviderMetadata["ID"])
}

func TestPlanChanges_DeleteMultipleRecordsOfTheSameHostname(t *testing.T) {
	plan := Plan{
		Current: []*Record{
			A("foo.docker.mydomain.com", "192.168.100.100", "ID_1"),
			A("foo.docker.mydomain.com", "192.168.100.101", "ID_11"),
			A("bar.docker.mydomain.com", "10.100.100.110", "ID_2"),
		},
		Desired: []*Record{
			A("bar.docker.mydomain.com", "10.100.100.110", ""),
		},
	}

	changes := plan.Changes()

	requireNumberOfChanges(t, changes, 0, 2)
	require.Equal(t, "foo.docker.mydomain.com", changes.Delete[0].Name)
	require.Equal(t, "ID_1", changes.Delete[0].ProviderMetadata["ID"])

	require.Equal(t, "foo.docker.mydomain.com", changes.Delete[1].Name)
	require.Equal(t, "ID_11", changes.Delete[1].ProviderMetadata["ID"])
}

func TestPlanChanges_DeleteRecordsOfDifferentTypes(t *testing.T) {
	plan := Plan{
		Current: []*Record{
			A("foo.docker.mydomain.com", "192.168.100.100", "ID_1"),
			A("foo.docker.mydomain.com", "192.168.100.101", "ID_11"),
			CNAME("foo.docker.mydomain.com", "woohoo.google.com", "ID_12"),
			A("bar.docker.mydomain.com", "10.100.100.110", "ID_2"),
		},
		Desired: []*Record{
			A("bar.docker.mydomain.com", "10.100.100.110", ""),
		},
	}

	changes := plan.Changes()

	requireNumberOfChanges(t, changes, 0, 3)
	require.Equal(t, "foo.docker.mydomain.com", changes.Delete[0].Name)
	require.Equal(t, "ID_1", changes.Delete[0].ProviderMetadata["ID"])

	require.Equal(t, "foo.docker.mydomain.com", changes.Delete[1].Name)
	require.Equal(t, "ID_11", changes.Delete[1].ProviderMetadata["ID"])

	require.Equal(t, "foo.docker.mydomain.com", changes.Delete[2].Name)
	require.Equal(t, "ID_12", changes.Delete[2].ProviderMetadata["ID"])
}

func TestPlanChanges_UpdateSingle(t *testing.T) {
	plan := Plan{
		Current: []*Record{
			A("foo.docker.mydomain.com", "192.168.100.100", "ID_1"),
		},
		Desired: []*Record{
			A("foo.docker.mydomain.com", "10.100.100.115", ""),
		},
	}

	changes := plan.Changes()

	requireNumberOfChanges(t, changes, 1, 1)
	require.Equal(t, "foo.docker.mydomain.com", changes.Delete[0].Name)
	require.Equal(t, "ID_1", changes.Delete[0].ProviderMetadata["ID"])

	require.Equal(t, "foo.docker.mydomain.com", changes.Create[0].Name)
	require.Equal(t, "10.100.100.115", changes.Create[0].Value)
}

func TestPlanChanges_UpdateMultiple(t *testing.T) {
	plan := Plan{
		Current: []*Record{
			A("foo.docker.mydomain.com", "192.168.100.100", "ID_1"),
			A("foo.docker.mydomain.com", "192.168.100.101", "ID_2"),
			A("bar.docker.mydomain.com", "192.168.101.100", "ID_3"),
			A("quux.docker.mydomain.com", "192.168.102.100", "ID_4"),
		},
		Desired: []*Record{
			A("foo.docker.mydomain.com", "10.100.100.115", ""),
			A("bar.docker.mydomain.com", "10.100.100.116", ""),
			A("quux.docker.mydomain.com", "192.168.102.100", ""),
		},
	}

	changes := plan.Changes()

	requireNumberOfChanges(t, changes, 2, 3)

	require.Equal(t, "foo.docker.mydomain.com", changes.Delete[0].Name)
	require.Equal(t, "ID_1", changes.Delete[0].ProviderMetadata["ID"])
	require.Equal(t, "foo.docker.mydomain.com", changes.Delete[1].Name)
	require.Equal(t, "ID_2", changes.Delete[1].ProviderMetadata["ID"])
	require.Equal(t, "bar.docker.mydomain.com", changes.Delete[2].Name)
	require.Equal(t, "ID_3", changes.Delete[2].ProviderMetadata["ID"])

	require.Equal(t, "foo.docker.mydomain.com", changes.Create[0].Name)
	require.Equal(t, "bar.docker.mydomain.com", changes.Create[1].Name)
}
