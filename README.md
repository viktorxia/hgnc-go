# HGNC-Go

A Go library for querying HGNC (Human Gene Nomenclature Committee) gene nomenclature data.



## Installation

```bash
go get github.com/viktorxia/hgnc-go
```



## Quick Start

### 1. Download HGNC Data

Source data  provided by the [HUGO Gene Nomenclature Committee](https://www.genenames.org/).

Download `hgnc_complete_set.txt` or `hgnc_complete_set.txt.gz` from [HGNC Downloads](https://www.genenames.org/download/archive/).

***\*Note\****: HGNC complete set field names and structure may change over time, which could cause compatibility issues. Users may need to modify the code accordingly based on their specific requirements and data versions. Please comply with HGNC data usage policies and requirements to avoid any economic and legal disputes. Users are responsible for ensuring their use of HGNC data follows the official terms and conditions.



### 2. Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    hgnc "github.com/viktorxia/hgnc-go"
)

func main() {
    // Load HGNC database
    db, err := hgnc.LoadTsv("data/hgnc_complete_set.txt.gz", true)
    if err != nil {
        log.Fatal(err)
    }
    
    // Check if gene is protein-coding
    isCoding := db.IsCodingGene("BRCA1")
    fmt.Printf("BRCA1 is coding gene: %v\n", isCoding)
    
    // Convert gene symbol to Entrez ID
    if entrezID, ok := db.SymbolToEntrezID("TP53"); ok {
        fmt.Printf("TP53 Entrez ID: %s\n", entrezID)
    }
}
```

**💡 See `example/basic/main.go` for a comprehensive example with most features demonstrated.**



## High-Level APIs

These APIs provide convenient methods for common gene queries and automatically handle multiple gene ID formats.

### Gene ID Formats Supported

The library automatically detects and handles multiple gene ID formats:

- **HGNC ID**: `HGNC:1100`
- **Gene Symbol**: `BRCA1`, `TP53` (including aliases and previous symbols)
- **Entrez ID**: `672` (numbers only)
- **Ensembl Gene ID**: `ENSG00000012048`
- **UCSC ID**: `uc002ict.4`

### Gene Classification

```go
// Check if gene is protein-coding
isCoding := db.IsCodingGene("BRCA1")        // Accepts multiple ID formats
isCoding := db.IsCodingGene("HGNC:1100")    // Same result
isCoding := db.IsCodingGene("672")          // Same result (Entrez ID)
```

### ID Conversions

```go
// Symbol conversions
entrezID, ok := db.SymbolToEntrezID("BRCA1")  // Symbol -> Entrez ID
ensg, ok := db.SymbolToEnsg("TP53")           // Symbol -> Ensembl Gene ID
ucscID, ok := db.SymbolToUcscID("EGFR")       // Symbol -> UCSC ID

// To Symbol conversions
symbol, ok := db.EntrezIDToSymbol("7157")     // Entrez ID -> Symbol
symbol, ok := db.EnsgToSymbol("ENSG00000141510")  // Ensembl ID -> Symbol
symbol, ok := db.UcscIDToSymbol("uc002ict.4")     // UCSC ID -> Symbol
```

### MANE Select Transcripts

```go
// Get complete MANE Select info (ENST|RefSeq format)
mane, ok := db.GetManeSelect("EGFR")

// Get only ENST transcript ID
enst, ok := db.GetManeSelectENST("EGFR") 

// Get only RefSeq transcript ID
refseq, ok := db.GetManeSelectRefseq("EGFR")
```

### RefSeq Accessions

```go
// Get RefSeq accessions for a gene
refseqAccs, ok := db.GeneRefseqAccs("BRCA1")  // Accepts multiple ID formats
```

### Symbol Normalization

The library automatically converts gene aliases and previous symbols to current HGNC symbols:

```go
// Enable auto-normalization (default)
db.SetAutoNormSymbol(true)

// Disable auto-normalization
db.SetAutoNormSymbol(false)
```



## Flexible Query Methods (Advanced Usage)

For more flexibility, use the `Fetch` and `Lookup` methods. Think of them as Unix commands:

- **Fetch**: Like `grep` - returns complete matching records
- **Lookup**: Like `grep | cut` - returns specific field values from matching records

### Fetch Method

Returns complete gene records (`[]*Records`) matching the query:

```go
// Find all records where symbol equals "BRCA1"
records := db.Fetch("BRCA1", hgnc.FIELD_SYMBOL)

for _, record := range records {
    fmt.Printf("HGNC ID: %s\n", record.HgncID())
    fmt.Printf("Symbol: %s\n", record.Symbol())
    fmt.Printf("Name: %s\n", record.Get(hgnc.FIELD_NAME))
    fmt.Printf("Entrez ID: %s\n", record.EntrezID())
}
```

### Lookup Method

Returns specific field values (`[]string`) for matching records:

```go
// Find records where symbol="BRCA1", then return their Entrez IDs
entrezIDs := db.Lookup("BRCA1", hgnc.FIELD_SYMBOL, hgnc.FIELD_ENTREZ_ID)

// Find records where symbol="TP53", then return their aliases
aliases := db.Lookup("TP53", hgnc.FIELD_SYMBOL, hgnc.FIELD_ALIAS_SYMBOL)

// Find records where Entrez ID="7157", then return their symbols
symbols := db.Lookup("7157", hgnc.FIELD_ENTREZ_ID, hgnc.FIELD_SYMBOL)
```



## Available Fields

### Indexed Fields (cached for performance)

- `FIELD_HGNC_ID` - HGNC ID
- `FIELD_SYMBOL` - Gene symbol
- `FIELD_ENTREZ_ID` - Entrez Gene ID
- `FIELD_ENSEMBL_GENE_ID` - Ensembl Gene ID
- `FIELD_UCSC_ID` - UCSC ID
- `FIELD_REFSEQ_ACCESSION` - RefSeq accession
- `FIELD_OMIM_ID` - OMIM ID

### Other Fields

- `FIELD_NAME` - Gene name
- `FIELD_LOCUS_GROUP` - Locus group
- `FIELD_LOCUS_TYPE` - Locus type
- `FIELD_LOCATION` - Chromosomal location
- `FIELD_ALIAS_SYMBOL` - Gene aliases
- `FIELD_PREV_SYMBOL` - Previous symbols
- `FIELD_MANE_SELECT` - MANE Select transcripts

💡 All fields are defined in `fields.go`



## Examples



### Mixed ID Format Processing

```go
// High-level APIs handle multiple formats automatically
genes := []string{"TP53", "BRCA1", "7157", "HGNC:11998", "ENSG00000141510"}
for _, gene := range genes {
    isCoding := db.IsCodingGene(gene)
    fmt.Printf("%s is coding gene: %v\n", gene, isCoding)
}
```



### MANE Select Processing

```go
gene := "BRCA1"
if mane, ok := db.GetManeSelect(gene); ok {
    fmt.Printf("Full MANE: %s\n", mane)  // e.g., "ENST00000357654.9|NM_007294.4"
    
    if enst, ok := db.GetManeSelectENST(gene); ok {
        fmt.Printf("ENST: %s\n", enst)   // e.g., "ENST00000357654.9"
    }
    
    if refseq, ok := db.GetManeSelectRefseq(gene); ok {
        fmt.Printf("RefSeq: %s\n", refseq)  // e.g., "NM_007294.4"
    }
}
```



### Advanced Queries with Core Methods

```go
// Find all protein-coding genes on chromosome 17
records := db.Fetch("protein-coding gene", hgnc.FIELD_LOCUS_GROUP)
for _, record := range records {
    if strings.Contains(record.Get(hgnc.FIELD_LOCATION), "17q") {
        fmt.Printf("Gene: %s, Location: %s\n", 
            record.Symbol(), record.Get(hgnc.FIELD_LOCATION))
    }
}

// Get all Ensembl IDs for genes with "BRCA" in their symbol
symbols := db.Lookup("BRCA1", hgnc.FIELD_SYMBOL, hgnc.FIELD_ENSEMBL_GENE_ID)
fmt.Printf("Ensembl IDs: %v\n", symbols)
```



## License

This project is licensed under the GNU General Public License v3.0 - see the LICENSE file for details.

