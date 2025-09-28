package hgnc_go

// Fetch retrieves records from HGNC based on the given value and query field.
// (similar to grep command in Unix)
func (h *HGNC) Fetch(value string, query Field) []*Record {

	if h == nil {
		panic("HGNC is nil")
	}

	if value == "" {
		return make([]*Record, 0)
	}

	if query == FIELD_SYMBOL {
		value = h.normalizeSymbol(value)
	}

	if _, ok1 := h.caches[query]; ok1 {
		// cached
		// h.caches[query][value] is a slice of indexes of h.records, type: []int
		if indexes, ok2 := h.caches[query][value]; ok2 {
			results := make([]*Record, 0, len(indexes))
			for _, index := range indexes {
				results = append(results, h.records[index])
			}
			return results
		}
		return make([]*Record, 0)
	} else {
		var results []*Record
		for _, record := range h.records {
			if record.data[query] == value {
				results = append(results, record)
			}
		}
		return results
	}
}

// Lookup retrieves values of target field for records in HGNC based on the given value and query field.
// (similar to grep + cut command in Unix)
func (h *HGNC) Lookup(value string, query, target Field) []string {

	if h == nil {
		panic("HGNC is nil")
	}

	if value == "" {
		return make([]string, 0)
	}

	if query == FIELD_SYMBOL {
		value = h.normalizeSymbol(value)
	}

	if _, ok1 := h.caches[query]; ok1 {
		// cached
		// hgnc.caches -> map[Field]Cache
		// hgnc.caches[field] -> cache -> map[string][]int
		// hgnc.caches[field][value] -> []int
		if indexes, ok2 := h.caches[query][value]; ok2 {
			results := make([]string, 0, len(indexes))
			for _, index := range indexes {
				results = append(results, h.records[index].data[target])
			}
			return results
		}
		return make([]string, 0)
	} else {
		// no cache
		var results []string
		for _, record := range h.records {
			if record.data[query] == value {
				results = append(results, record.data[target])
			}
		}
		return results
	}
}
