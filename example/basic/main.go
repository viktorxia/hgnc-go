/* You may try: go run example/basic/main.go */

package main

import (
	"fmt"
	"log"

	h "github.com/viktorxia/hgnc-go"
)

func main() {

	// ------------------------------------------------------------------------------------------
	// Init: Load HGNC database
	// Download hgnc_complete_set.txt from: https://www.genenames.org/download/archive/
	hgnc, err := h.LoadTsv("data/hgnc_complete_set.txt.gz", true)
	if err != nil {
		log.Fatalf("Failed to load HGNC data: %v", err)
	}
	fmt.Println("=== HGNC Database Loaded Successfully ===\n")

	// ------------------------------------------------------------------------------------------
	/*
			[High-level APIs]

			* High-level APIs are thin wrappers of Fetch & Lookup methods that could
			satisfy most of the requirements in a medical project. If you need more
			flexibility, you can use Fetch & Lookup methods directly.

			* Gene-related APIs usually accept multiple expression of a gene, such as
			HGNC ID, symbol, symbol alias, previous symbol, Ensembl gene ID (ENSG),
		    UCSC gene ID (uc) and Entrez/NCBI gene ID (integer). Classification of the
			naming system is automatically inferred from the input string.

			* If input gene symbol string is a previous symbol or an alias, it will be
			automatically converted to the current symbol. You don't need to worry about
			using a old symbol as long as it is recorded in the hgnc complete dataset.
	*/

	fmt.Println("=== High-level API Examples ===")

	// Test coding gene detection
	genes := []string{"BRCA1", "HGNC:1100", "7157", "LINC00115"}
	for _, gene := range genes {
		isCoding := hgnc.IsCodingGene(gene)
		fmt.Printf("%-12s -> Coding gene: %v\n", gene, isCoding)
	}
	fmt.Println()

	// Gene ID conversions
	fmt.Println("=== Gene ID Conversions ===")
	if entrezID, ok := hgnc.SymbolToEntrezID("TP53"); ok {
		fmt.Printf("TP53 -> Entrez ID: %s\n", entrezID)
	}

	if symbol, ok := hgnc.EntrezIDToSymbol("7157"); ok {
		fmt.Printf("7157 -> Symbol: %s\n", symbol)
	}

	if ensg, ok := hgnc.SymbolToEnsg("BRCA1"); ok {
		fmt.Printf("BRCA1 -> Ensembl ID: %s\n", ensg)
	}

	// MANE Select transcripts
	if mane, ok := hgnc.GetManeSelect("TP53"); ok {
		fmt.Printf("TP53 -> MANE Select: %s\n", mane)
	}
	if mane, ok := hgnc.GetManeSelectENST("TP53"); ok {
		fmt.Printf("TP53 -> MANE Select (ENST): %s\n", mane)
	}
	if mane, ok := hgnc.GetManeSelectRefseq("TP53"); ok {
		fmt.Printf("TP53 -> MANE Select (RefSeq): %s\n", mane)
	}
	fmt.Println()

	// ------------------------------------------------------------------------------------------
	/*
		[Custom query: Fetch & Lookup methods]

		* Fetch & Lookup methods are the core methods of HGNC struct.
		* Fetch method returns a list of Records that match the query string.
		* Lookup method returns a list of strings that match the query string.

		* Both methods use cache to improve performance, please refer to field.go
		to see all cached fields.
	*/

	fmt.Println("=== Fetch Method Examples ===")

	// Fetch complete records
	records := hgnc.Fetch("BRCA1", h.FIELD_SYMBOL)
	fmt.Printf("Found %d records for BRCA1\n", len(records))

	for i, record := range records {
		fmt.Printf("Record %d:\n", i+1)

		// Convert to map for easy access
		data := record.ToMap()
		fmt.Printf("  HGNC ID: %s\n", data[h.FIELD_HGNC_ID])
		fmt.Printf("  Symbol: %s\n", data[h.FIELD_SYMBOL])
		fmt.Printf("  Name: %s\n", data[h.FIELD_NAME])

		// Or use direct accessors
		fmt.Printf("  Locus Group: %s\n", record.Get(h.FIELD_LOCUS_GROUP))
		fmt.Printf("  Status: %s\n", record.Get(h.FIELD_STATUS))

		// Convenient accessor methods
		fmt.Printf("  HGNC ID (direct): %s\n", record.HgncID())
		fmt.Printf("  Symbol (direct): %s\n", record.Symbol())
		fmt.Printf("  Entrez ID: %s\n", record.EntrezID())
		break // Show only first record for brevity
	}
	fmt.Println()

	fmt.Println("=== Lookup Method Examples ===")

	// Lookup specific fields
	maneSelects := hgnc.Lookup("BRCA1", h.FIELD_SYMBOL, h.FIELD_MANE_SELECT)
	fmt.Printf("BRCA1 MANE Select transcripts: %v\n", maneSelects)

	entrezIDs := hgnc.Lookup("BRCA1", h.FIELD_SYMBOL, h.FIELD_ENTREZ_ID)
	fmt.Printf("BRCA1 Entrez IDs: %v\n", entrezIDs)

	aliases := hgnc.Lookup("BRCA1", h.FIELD_SYMBOL, h.FIELD_ALIAS_SYMBOL)
	fmt.Printf("BRCA1 aliases: %v\n", aliases)

	// Multiple genes lookup
	fmt.Println("\n=== Batch Processing Example ===")
	testGenes := []string{"TP53", "BRCA1", "EGFR", "NOTEXIST"}

	for _, gene := range testGenes {
		locusGroups := hgnc.Lookup(gene, h.FIELD_SYMBOL, h.FIELD_LOCUS_GROUP)
		if len(locusGroups) > 0 {
			fmt.Printf("%-10s -> Locus Group: %s\n", gene, locusGroups[0])
		} else {
			fmt.Printf("%-10s -> Not found\n", gene)
		}
	}
	fmt.Println()

	// ------------------------------------------------------------------------------------------

	// Symbol normalization control
	fmt.Println("=== Symbol Normalization ===")
	fmt.Println("test GBA -> GBA1 gene symbol normalization")

	hgnc.SetAutoNormSymbol(true)
	fmt.Println("Auto-normalization enabled (default)")
	_count := 0
	for _, record := range hgnc.Fetch("GBA", h.FIELD_SYMBOL) {
		_count++
		fmt.Printf("Found record. Symbol: %s, HGNC ID: %s\n", record.Symbol(), record.HgncID())
	}
	fmt.Printf("Found %d records for GBA\n", _count)

	hgnc.SetAutoNormSymbol(false)
	fmt.Println("Auto-normalization disabled")
	_count = 0
	for _, record := range hgnc.Fetch("GBA", h.FIELD_SYMBOL) {
		fmt.Printf("Found record. Symbol: %s, HGNC ID: %s\n", record.Symbol(), record.HgncID())
	}
	fmt.Printf("Found %d records for GBA\n", _count)

	// ------------------------------------------------------------------------------------------
	fmt.Println("\n=== Processing Complete ===")
}
