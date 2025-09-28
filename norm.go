package hgnc_go

// normalizeSymbol converts alias/previous symbols to standard HGNC symbols.
func (h *HGNC) normalizeSymbol(symbol string) string {

	if !h.autoNormSymbol {
		return symbol
	}

	if _, ok := h.stdHgncSymbols[symbol]; ok {
		return symbol
	}
	if stdSymbol, ok := h.geneSymbolMap[symbol]; ok {
		return stdSymbol
	}
	return symbol
}
