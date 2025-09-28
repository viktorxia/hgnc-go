package hgnc_go

import (
	"bufio"
	"compress/gzip"
	"errors"
	"os"
	"strings"
)

// Cache is a map of field to a slice of integers.
// Each integer represents a index of HGNC.records.
type Cache map[string][]int

type HGNC struct {
	records        []*Record           // all records in HGNC file
	geneSymbolMap  map[string]string   // cache, key = symbol, value = standard HGNC symbol
	stdHgncSymbols map[string]struct{} // cache, key = standard HGNC symbol, value = empty struct{}
	caches         map[Field]Cache     // cache for some important fields
	autoNormSymbol bool                // whether to normalize symbol automatically
}

func (h *HGNC) SetAutoNormSymbol(autoNormSymbol bool) {
	h.autoNormSymbol = autoNormSymbol
}

// LoadTsv is the constructor of HGNC struct.
func LoadTsv(filepath string, gzipped bool) (*HGNC, error) {

	// init
	h := &HGNC{
		records:        make([]*Record, 0),
		geneSymbolMap:  make(map[string]string),
		stdHgncSymbols: make(map[string]struct{}),
		caches:         make(map[Field]Cache),
		autoNormSymbol: true,
	}

	for _, field := range indexedFields {
		// h.caches -> map[Field]Cache
		// h.caches[field] -> cache -> map[string][]int
		// h.caches[field][value] -> []int
		cache := make(Cache)
		h.caches[field] = cache
	}

	// open file
	fh, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer fh.Close()

	var scanner *bufio.Scanner

	if !gzipped {
		scanner = bufio.NewScanner(fh)
	} else {
		gz, err := gzip.NewReader(fh)
		if err != nil {
			return nil, err
		}
		scanner = bufio.NewScanner(gz)
		defer gz.Close()
	}

	// read header line
	if !scanner.Scan() {
		return nil, errors.New("failed reading header line")
	}
	headerLine := scanner.Text()
	headerMap := make(map[string]int)
	for i, field := range strings.Split(headerLine, "\t") {
		headerMap[field] = i
	}

	// collect data
	recordIdx := 0
	for scanner.Scan() {
		line := scanner.Text()
		record := line2Record(line, headerMap)

		// records
		h.records = append(h.records, record)

		// standard symbols
		h.stdHgncSymbols[record.data[FIELD_SYMBOL]] = struct{}{}

		// alias & prev symbols
		aliasSymbolStr := record.data[FIELD_ALIAS_SYMBOL]
		prevSymbolStr := record.data[FIELD_PREV_SYMBOL]
		if aliasSymbolStr != "" {
			for _, alias := range strings.Split(aliasSymbolStr, "|") {
				alias = strings.TrimSpace(alias)
				if alias != "" {
					h.geneSymbolMap[alias] = record.data[FIELD_SYMBOL]
				}
			}
		}
		if prevSymbolStr != "" {
			for _, prevSymbol := range strings.Split(prevSymbolStr, "|") {
				prevSymbol = strings.TrimSpace(prevSymbol)
				if prevSymbol != "" {
					h.geneSymbolMap[prevSymbol] = record.data[FIELD_SYMBOL]
				}
			}
		}

		// caches
		for _, field := range indexedFields {
			value := record.data[field]
			// h.caches -> map[Field]Cache
			// h.caches[field] -> cache -> map[string][]int
			// h.caches[field][value] -> []int
			if value == "" {
				continue
			}
			if _, ok := h.caches[field][value]; !ok {
				h.caches[field][value] = []int{recordIdx}
			} else {
				h.caches[field][value] = append(
					h.caches[field][value], recordIdx,
				)
			}
		}

		recordIdx++
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return h, nil
}

// line2Record converts a line of HGNC file to a Record struct.
func line2Record(line string, headerMap map[string]int) *Record {

	record := new(Record)
	record.data = make(map[Field]string)

	l := strings.Split(line, "\t")

	for fieldName, tsvIdx := range headerMap {
		if tsvIdx < len(l) {
			// !!! some fields are quoted with double quotes,
			// or with spaces at the beginning or end.
			value := strings.Trim(l[tsvIdx], "\"")
			value = strings.TrimSpace(value)
			record.data[Field(fieldName)] = value
		} else {
			record.data[Field(fieldName)] = ""
		}
	}

	return record
}
