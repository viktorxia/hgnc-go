package hgnc_go

import (
	"strconv"
	"strings"
)

/*
In this file, 'gene' could be:
	- symbol (HGNC)
	- entrez id (pure number)
	- hgnc id (starts with 'HGNC:')
	- ensembl gene id (starts with 'ENSG')
	- ucsc id (starts with 'uc')

classifyGeneStringSystem() function can classify the 'gene' and return the field type.
*/

// classifyGeneStringSystem classifies the 'gene' string and returns the field type.
func classifyGeneStringSystem(gene string) Field {
	if strings.HasPrefix(gene, "HGNC:") {
		return FIELD_HGNC_ID
	} else if strings.HasPrefix(gene, "ENSG") {
		return FIELD_ENSEMBL_GENE_ID
	} else if strings.HasPrefix(gene, "uc") {
		return FIELD_UCSC_ID
	} else if _, err := strconv.Atoi(gene); err == nil {
		return FIELD_ENTREZ_ID
	} else {
		return FIELD_SYMBOL
	}
}

// GetManeSelect gets mane select transcript for a gene
func (h *HGNC) GetManeSelect(gene string) (string, bool) {
	field := classifyGeneStringSystem(gene)
	if result := h.Lookup(gene, field, FIELD_MANE_SELECT); len(result) > 0 {
		return result[0], true
	}
	return "", false
}

// GetManeSelectENST gets mane select transcript for a gene and returns only the ENST id
func (h *HGNC) GetManeSelectENST(gene string) (string, bool) {
	result, found := h.GetManeSelect(gene)
	if found {
		return strings.Split(result, "|")[0], true
	}
	return "", false
}

// GetManeSelectRefseq gets mane select transcript for a gene and returns only the RefSeq id
func (h *HGNC) GetManeSelectRefseq(gene string) (string, bool) {
	result, found := h.GetManeSelect(gene)
	if found {
		split := strings.Split(result, "|")
		if len(split) > 1 {
			return split[1], true
		}
	}
	return "", false
}

// IsCodingGene checks if a gene is protein-coding by it's locus group
func (h *HGNC) IsCodingGene(gene string) bool {
	field := classifyGeneStringSystem(gene)
	if result := h.Lookup(gene, field, FIELD_LOCUS_GROUP); len(result) > 0 {
		if strings.HasPrefix(result[0], "protein-coding") {
			return true
		}
	}
	return false
}

// EntrezIDToSymbol converts entrez id to gene symbol
func (h *HGNC) EntrezIDToSymbol(entrezID string) (string, bool) {
	if result := h.Lookup(entrezID, FIELD_ENTREZ_ID, FIELD_SYMBOL); len(result) > 0 {
		return result[0], true
	}
	return "", false
}

// SymbolToEntrezID convert gene symbol to entrez id
func (h *HGNC) SymbolToEntrezID(symbol string) (string, bool) {
	if result := h.Lookup(symbol, FIELD_SYMBOL, FIELD_ENTREZ_ID); len(result) > 0 {
		return result[0], true
	}
	return "", false
}

// EnsgToSymbol converts ensembl gene id to gene symbol
func (h *HGNC) EnsgToSymbol(ensg string) (string, bool) {
	ensg = strings.Split(ensg, ".")[0]
	if result := h.Lookup(ensg, FIELD_ENSEMBL_GENE_ID, FIELD_SYMBOL); len(result) > 0 {
		return result[0], true
	}
	return "", false
}

// SymbolToEnsg converts gene symbol to ensembl gene id
func (h *HGNC) SymbolToEnsg(symbol string) (string, bool) {
	if result := h.Lookup(symbol, FIELD_SYMBOL, FIELD_ENSEMBL_GENE_ID); len(result) > 0 {
		return result[0], true
	}
	return "", false
}

// UcscIDToSymbol converts ucsc id to gene symbol
func (h *HGNC) UcscIDToSymbol(ucscID string) (string, bool) {
	if result := h.Lookup(ucscID, FIELD_UCSC_ID, FIELD_SYMBOL); len(result) > 0 {
		return result[0], true
	}
	return "", false
}

// SymbolToUcscID converts gene symbol to ucsc id
func (h *HGNC) SymbolToUcscID(symbol string) (string, bool) {
	if result := h.Lookup(symbol, FIELD_SYMBOL, FIELD_UCSC_ID); len(result) > 0 {
		return result[0], true
	}
	return "", false
}

// GeneRefseqAccs gets refseq accessions for a gene
func (h *HGNC) GeneRefseqAccs(gene string) (string, bool) {

	field := classifyGeneStringSystem(gene)
	if result := h.Lookup(gene, field, FIELD_REFSEQ_ACCESSION); len(result) > 0 {
		return result[0], true
	}
	return "", false
}
