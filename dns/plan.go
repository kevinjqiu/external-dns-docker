package dns

import "fmt"

type Plan struct {
	Current []*Record
	Desired []*Record
}

func buildMap(list []*Record) map[string][]*Record {
	m := make(map[string][]*Record)

	for _, item := range list {
		key := fmt.Sprintf("%v:%v", item.Name, item.Type)
		_, ok := m[key]
		if !ok {
			m[key] = make([]*Record, 0)
		}
		m[key] = append(m[key], item)
	}

	return m
}

func recordSetByValue(records []*Record) map[string]*Record {
	ret := make(map[string]*Record)

	for _, record := range records {
		ret[record.Value] = record
	}

	return ret
}

func generateCreateList(currentMap, desiredMap map[string][]*Record) []*Record {
	records := make([]*Record, 0)

	for desiredKey, desiredRecords := range desiredMap {
		currentRecords, ok := currentMap[desiredKey]

		var valueSet = make(map[string]*Record)

		if ok {
			valueSet = recordSetByValue(currentRecords)
		}

		for _, record := range desiredRecords {
			if _, ok := valueSet[record.Value]; !ok {
				records = append(records, record)
			}
		}
	}

	return records
}

func generateDeleteList(currentMap, desiredMap map[string][]*Record) []*Record {
	records := make([]*Record, 0)

	for currentKey, currentRecords := range currentMap {
		desiredRecords, ok := desiredMap[currentKey]

		var valueSet = make(map[string]*Record)

		if ok {
			valueSet = recordSetByValue(desiredRecords)
		}

		for _, record := range currentRecords {
			if _, ok := valueSet[record.Value]; !ok {
				records = append(records, record)
			}
		}
	}

	return records
}

func (p *Plan) Changes() *Changes {
	currentMap := buildMap(p.Current)
	desiredMap := buildMap(p.Desired)

	return &Changes{
		Create: generateCreateList(currentMap, desiredMap),
		Delete: generateDeleteList(currentMap, desiredMap),
	}
}

type Changes struct {
	Create []*Record
	Delete []*Record
}
