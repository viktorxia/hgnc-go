# HGNC-Go

A Go library for querying HGNC (Human Gene Nomenclature Committee) gene nomenclature data.



## 1. Installation

```bash
go get github.com/viktorxia/hgnc-go
```



## 2. Quick Start

### 2.1 Download HGNC Data

Source data  provided by the [HUGO Gene Nomenclature Committee](https://www.genenames.org/).

Download `hgnc_complete_set.txt` or `hgnc_complete_set.txt.gz` from [HGNC Downloads](https://www.genenames.org/download/archive/).

***\*Note\****: HGNC complete set field names and structure may change over time, which could cause compatibility issues. Users may need to modify the code accordingly based on their specific requirements and data versions. Please comply with HGNC data usage policies and requirements to avoid any economic and legal disputes. Users are responsible for ensuring their use of HGNC data follows the official terms and conditions.



### 2.2 Basic Usage

```go
package main

import (
    "fmt"
    "log"
    
    h "github.com/viktorxia/hgnc-go"
)

func main() {
    // Load HGNC database
    hgnc, err := h.LoadTsv("data/hgnc_complete_set.txt.gz", true)
    if err != nil {
        log.Fatal(err)
    }
    
    // Check if gene is protein-coding
    isCoding := hgnc.IsCodingGene("BRCA1")
    fmt.Printf("BRCA1 is coding gene: %v\n", isCoding)
    
    // Convert gene symbol to Entrez ID
    if entrezID, ok := hgnc.SymbolToEntrezID("TP53"); ok {
        fmt.Printf("TP53 Entrez ID: %s\n", entrezID)
    }
}
```

**💡 See `example/basic/main.go` for a comprehensive example with most features demonstrated.**



## 3. High-Level APIs

These APIs provide convenient methods for common gene queries and automatically handle multiple gene ID formats.

The library automatically detects and handles multiple gene ID formats. You can choose one of the following gene representations:

- **HGNC ID**: `HGNC:1100`
- **Gene Symbol**: `BRCA1`, `TP53` (including aliases and previous symbols)
- **Entrez ID**: `672` (numbers only)
- **Ensembl Gene ID**: `ENSG00000012048`
- **UCSC ID**: `uc002ict.4`



### 3.1 Coding Gene?

```go
// Check if gene is protein-coding
isCoding := hgnc.IsCodingGene("BRCA1")        // Accepts multiple ID formats
isCoding := hgnc.IsCodingGene("HGNC:1100")    // Same result
isCoding := hgnc.IsCodingGene("672")          // Same result (Entrez ID)
```



### 3.2 ID Conversion

hgnc-go support convertion between the following systems: 

- EntrezID
- gene symbol (HGNC)
- UCSC ID
- Ensembl gene ID

```go
// Symbol conversions
entrezID, ok := hgnc.SymbolToEntrezID("BRCA1")  // Symbol -> Entrez ID
ensg, ok := hgnc.SymbolToEnsg("TP53")           // Symbol -> Ensembl Gene ID
ucscID, ok := hgnc.SymbolToUcscID("EGFR")       // Symbol -> UCSC ID

// To Symbol conversions
symbol, ok := hgnc.EntrezIDToSymbol("7157")     // Entrez ID -> Symbol
symbol, ok := hgnc.EnsgToSymbol("ENSG00000141510")  // Ensembl ID -> Symbol
symbol, ok := hgnc.UcscIDToSymbol("uc002ict.4")     // UCSC ID -> Symbol
```



### 3.3 MANE Select Transcripts

```go
// Get complete MANE Select info (ENST|RefSeq format)
mane, ok := hgnc.GetManeSelect("EGFR")

// Get only ENST transcript ID
enst, ok := hgnc.GetManeSelectENST("EGFR") 

// Get only RefSeq transcript ID
refseq, ok := hgnc.GetManeSelectRefseq("EGFR")
```



### 3.4 RefSeq Accessions

```go
// Get RefSeq accessions for a gene
refseqAccs, ok := hgnc.GeneRefseqAccs("BRCA1")  // Accepts multiple ID formats
```



### 3.5 Symbol Normalization

The library automatically converts aliases and previous symbols of genes to the HGNC symbols:

```go
// Enable auto-normalization (default)
hgnc.SetAutoNormSymbol(true)

// Disable auto-normalization
hgnc.SetAutoNormSymbol(false)
```



## 4. Flexible Query Methods

For more flexibility, use the `Fetch` and `Lookup` methods. Think of them as Unix commands:

- **Fetch**: Like `grep` - returns complete matching records
- **Lookup**: Like `grep | cut` - returns specific field values from matching records

### 4.1 Fetch Record(s)

Returns complete gene records (`[]*Records`) matching the query:

```go
// Find all records where symbol equals "BRCA1"
records := hgnc.Fetch("BRCA1", h.FIELD_SYMBOL)

for _, record := range records {
    fmt.Printf("HGNC ID: %s\n", record.HgncID())
    fmt.Printf("Symbol: %s\n", record.Symbol())
    fmt.Printf("Name: %s\n", record.Get(hgnc.FIELD_NAME))
    fmt.Printf("Entrez ID: %s\n", record.EntrezID())
}
```

### 4.2 Lookup Value(s)

Returns specific field values (`[]string`) for matching records:

```go
// Find records where symbol="BRCA1", then return their Entrez IDs
entrezIDs := hgnc.Lookup("BRCA1", h.FIELD_SYMBOL, h.FIELD_ENTREZ_ID)

// Find records where symbol="TP53", then return their aliases
aliases := hgnc.Lookup("TP53", h.FIELD_SYMBOL, h.FIELD_ALIAS_SYMBOL)

// Find records where Entrez ID="7157", then return their symbols
symbols := hgnc.Lookup("7157", h.FIELD_ENTREZ_ID, h.FIELD_SYMBOL)
```



## 5. Field & Performance Guide

**Indexed fields are 1,000-10,000x faster** than non-indexed fields!

Lookup time comparison in `example/cache_vs_nocache.go`

| Cached | Operations | Goroutines | Time     | Operations/sec |
| ------ | ---------- | ---------- | -------- | -------------- |
| Y      | 1000       | 1          | 443.61µs | 2,254,217      |
| N      | 1000       | 1          | 5.078s   | 197            |
| Y      | 5000       | 2          | 2.80ms   | 1,786,357      |
| Y      | 5000       | 2          | 13.075s  | 382            |
| N      | 10000      | 4          | 1.99ms   | 5,020,763      |
| N      | 10000      | 4          | 14.412s  | 694            |

Indexed Fields (Fast)

- `FIELD_HGNC_ID` - HGNC ID
- `FIELD_SYMBOL` - Gene symbol
- `FIELD_ENTREZ_ID` - Entrez Gene ID
- `FIELD_ENSEMBL_GENE_ID` - Ensembl Gene ID
- `FIELD_UCSC_ID` - UCSC ID
- `FIELD_REFSEQ_ACCESSION` - RefSeq accession
- `FIELD_OMIM_ID` - OMIM ID

Other Fields (Slow)

- `FIELD_NAME` - Gene name
- `FIELD_LOCUS_GROUP` - Locus group
- `FIELD_LOCUS_TYPE` - Locus type
- `FIELD_LOCATION` - Chromosomal location
- `FIELD_ALIAS_SYMBOL` - Gene aliases
- `FIELD_PREV_SYMBOL` - Previous symbols
- `FIELD_MANE_SELECT` - MANE Select transcripts

```go
// FAST - O(1) hash lookup
gene := hgnc.Fetch("BRCA1", h.FIELD_SYMBOL)

// SLOW - O(n) linear scan through ~45k records  
gene := hgnc.Fetch("breast cancer 1", h.FIELD_NAME)
```

💡 All fields are defined in `fields.go`

**Test performance yourself:** `go run example/cache_vs_nocache/main.go`





## 6. Examples



### Mixed ID Format Processing

```go
// High-level APIs handle multiple formats automatically
genes := []string{"TP53", "BRCA1", "7157", "HGNC:11998", "ENSG00000141510"}
for _, gene := range genes {
    isCoding := hgnc.IsCodingGene(gene)
    fmt.Printf("%s is coding gene: %v\n", gene, isCoding)
}
```



### MANE Select Processing

```go
gene := "BRCA1"
if mane, ok := hgnc.GetManeSelect(gene); ok {
    fmt.Printf("Full MANE: %s\n", mane)  // e.g., "ENST00000357654.9|NM_007294.4"
    
    if enst, ok := hgnc.GetManeSelectENST(gene); ok {
        fmt.Printf("ENST: %s\n", enst)   // e.g., "ENST00000357654.9"
    }
    
    if refseq, ok := hgnc.GetManeSelectRefseq(gene); ok {
        fmt.Printf("RefSeq: %s\n", refseq)  // e.g., "NM_007294.4"
    }
}
```



### Advanced Queries with Fetch/Lookup

```go
// Find all protein-coding genes on chromosome 17
records := hgnc.Fetch("protein-coding gene", h.FIELD_LOCUS_GROUP)
for _, record := range records {
    if strings.Contains(record.Get(h.FIELD_LOCATION), "17q") {
        fmt.Printf("Gene: %s, Location: %s\n", 
            record.Symbol(), record.Get(h.FIELD_LOCATION))
    }
}

// Get all Ensembl IDs for genes with "BRCA" in their symbol
symbols := hgnc.Lookup("BRCA1", h.FIELD_SYMBOL, h.FIELD_ENSEMBL_GENE_ID)
fmt.Printf("Ensembl IDs: %v\n", symbols)
```



## 7. License

This project is licensed under the GNU General Public License v3.0 - see the LICENSE file for details.

