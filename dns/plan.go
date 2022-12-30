package dns

type Plan struct {
	Current []*Record
	Desired []*Record
}

func buildMap(list []*Record) map[string][]*Record {
	m := make(map[string][]*Record)

	for _, item := range list {
		_, ok := m[item.Name]
		if !ok {
			m[item.Name] = make([]*Record, 0)
		}
		m[item.Name] = append(m[item.Name], item)
	}

	return m
}

func generateCreateList(currentMap, desiredMap map[string][]*Record) []*Record {
	records := make([]*Record, 0)

	for desiredKey, desiredRecords := range desiredMap {
		_, ok := currentMap[desiredKey]

		if ok {
			continue
		}

		for _, record := range desiredRecords {
			records = append(records, record)
		}
	}

	return records
}

func generateDeleteList(currentMap, desiredMap map[string][]*Record) []*Record {
	records := make([]*Record, 0)

	for currentKey, currentRecords := range currentMap {
		_, ok := desiredMap[currentKey]

		if ok {
			continue
		}

		for _, record := range currentRecords {
			records = append(records, record)
		}
	}

	return records
}

func generateUpdateList(currentMap, desiredMap map[string][]*Record) []*Record {
	records := make([]*Record, 0)

	commonKeys := make(map[string]struct{})

	for currentKey := range currentMap {
		if _, ok := desiredMap[currentKey]; ok {
			commonKeys[currentKey] = struct{}{}
		}
	}

	for desiredKey := range desiredMap {
		if _, ok := currentMap[desiredKey]; ok {
			commonKeys[desiredKey] = struct{}{}
		}
	}

	for commonKey := range commonKeys {
		currentRecords := currentMap[commonKey]
		desiredRecords := desiredMap[commonKey]

	}

	return records
}

func (p *Plan) Changes() *Changes {
	currentMap := buildMap(p.Current)
	desiredMap := buildMap(p.Desired)

	return &Changes{
		Create: generateCreateList(currentMap, desiredMap),
		Delete: generateDeleteList(currentMap, desiredMap),
		Update: generateUpdateList(currentMap, desiredMap),
	}
}

type Changes struct {
	Create []*Record
	Update []*Record
	Delete []*Record
}
