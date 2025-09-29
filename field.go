package hgnc_go

type Field string

const (

	// ---------------- indexed

	FIELD_HGNC_ID          Field = "hgnc_id"          // #1
	FIELD_SYMBOL           Field = "symbol"           // #2
	FIELD_ENTREZ_ID        Field = "entrez_id"        // #19
	FIELD_ENSEMBL_GENE_ID  Field = "ensembl_gene_id"  // #20
	FIELD_UCSC_ID          Field = "ucsc_id"          // #22
	FIELD_REFSEQ_ACCESSION Field = "refseq_accession" // #24
	FIELD_OMIM_ID          Field = "omim_id"          // #32

	// ---------------- others

	FIELD_NAME                     Field = "name"                     // #3
	FIELD_LOCUS_GROUP              Field = "locus_group"              // #4
	FIELD_LOCUS_TYPE               Field = "locus_type"               // #5
	FIELD_STATUS                   Field = "status"                   // #6
	FIELD_LOCATION                 Field = "location"                 // #7
	FIELD_LOCATION_SORTABLE        Field = "location_sortable"        // #8
	FIELD_ALIAS_SYMBOL             Field = "alias_symbol"             // #9
	FIELD_ALIAS_NAME               Field = "alias_name"               // #10
	FIELD_PREV_SYMBOL              Field = "prev_symbol"              // #11
	FIELD_PREV_NAME                Field = "prev_name"                // #12
	FIELD_GENE_FAMILY              Field = "gene_family"              // #13
	FIELD_GENE_FAMILY_ID           Field = "gene_family_id"           // #14
	FIELD_DATE_APPROVED_RESERVED   Field = "date_approved_reserved"   // #15
	FIELD_DATE_SYMBOL_CHANGED      Field = "date_symbol_changed"      // #16
	FIELD_DATE_NAME_CHANGED        Field = "date_name_changed"        // #17
	FIELD_DATE_MODIFIED            Field = "date_modified"            // #18
	FIELD_VEGA_ID                  Field = "vega_id"                  // #21
	FIELD_ENA                      Field = "ena"                      // #23
	FIELD_CCDS_ID                  Field = "ccds_id"                  // #25
	FIELD_UNIPROT_IDS              Field = "uniprot_ids"              // #26
	FIELD_PUBMED_ID                Field = "pubmed_id"                // #27
	FIELD_MGD_ID                   Field = "mgd_id"                   // #28
	FIELD_RGD_ID                   Field = "rgd_id"                   // #29
	FIELD_LSDB                     Field = "lsdb"                     // #30
	FIELD_COSMIC                   Field = "cosmic"                   // #31
	FIELD_MIRBASE                  Field = "mirbase"                  // #33
	FIELD_HOMEODB                  Field = "homeodb"                  // #34
	FIELD_SNORNABASE               Field = "snornabase"               // #35
	FIELD_BIOPARADIGMS_SLC         Field = "bioparadigms_slc"         // #36
	FIELD_ORPHANET                 Field = "orphanet"                 // #37
	FIELD_PSEUDOGENE_ORG           Field = "pseudogene.org"           // #38
	FIELD_HORDE_ID                 Field = "horde_id"                 // #39
	FIELD_MEROPS                   Field = "merops"                   // #40
	FIELD_IMGT                     Field = "imgt"                     // #41
	FIELD_IUPHAR                   Field = "iuphar"                   // #42
	FIELD_KZNF_GENE_CATALOG        Field = "kznf_gene_catalog"        // #43
	FIELD_MAMIT_TRNADB             Field = "mamit-trnadb"             // #44
	FIELD_CD                       Field = "cd"                       // #45
	FIELD_LNCRNADB                 Field = "lncrnadb"                 // #46
	FIELD_ENZYME_ID                Field = "enzyme_id"                // #47
	FIELD_INTERMEDIATE_FILAMENT_DB Field = "intermediate_filament_db" // #48
	FIELD_AGR                      Field = "agr"                      // #49
	FIELD_MANE_SELECT              Field = "mane_select"              // #50
)

var indexedFields = []Field{
	FIELD_HGNC_ID,
	FIELD_SYMBOL,
	FIELD_ENTREZ_ID,
	FIELD_ENSEMBL_GENE_ID,
	FIELD_UCSC_ID,
	FIELD_REFSEQ_ACCESSION,
	FIELD_OMIM_ID,
}

func GetAllIndexedFieldNames() []string {
	result := make([]string, len(indexedFields))
	for i, f := range indexedFields {
		result[i] = string(f)
	}
	return result
}

var fieldDesc = map[Field]string{
	FIELD_HGNC_ID:                  "HGNC ID. A unique ID created by the HGNC for every approved symbol.",
	FIELD_SYMBOL:                   "The HGNC approved gene symbol. Equates to the \"APPROVED SYMBOL\" field within the gene symbol report.",
	FIELD_NAME:                     "HGNC approved name for the gene. Equates to the \"APPROVED NAME\" field within the gene symbol report.",
	FIELD_LOCUS_GROUP:              "A group name for a set of related locus types as defined by the HGNC (e.g. non-coding RNA).",
	FIELD_LOCUS_TYPE:               "The locus type as defined by the HGNC (e.g. RNA, transfer).",
	FIELD_STATUS:                   "Status of the symbol report, which can be either \"Approved\" or \"Entry Withdrawn\".",
	FIELD_LOCATION:                 "Cytogenetic location of the gene (e.g. 2q34).",
	FIELD_LOCATION_SORTABLE:        "Same as \"location\" but single digit chromosomes are prefixed with a 0 enabling them to be sorted in correct numerical order (e.g. 02q34).",
	FIELD_ALIAS_SYMBOL:             "Other symbols used to refer to this gene as seen in the \"SYNONYMS\" field in the symbol report.",
	FIELD_ALIAS_NAME:               "Other names used to refer to this gene as seen in the \"SYNONYMS\" field in the gene symbol report.",
	FIELD_PREV_SYMBOL:              "Symbols previously approved by the HGNC for this gene. Equates to the \"PREVIOUS SYMBOLS & NAMES\" field within the gene symbol report.",
	FIELD_PREV_NAME:                "Gene names previously approved by the HGNC for this gene. Equates to the \"PREVIOUS SYMBOLS & NAMES\" field within the gene symbol report.",
	FIELD_GENE_FAMILY:              "Name given to a gene family or group the gene has been assigned to. Equates to the \"GENE FAMILY\" field within the gene symbol report.",
	FIELD_GENE_FAMILY_ID:           "ID used to designate a gene family or group the gene has been assigned to.",
	FIELD_DATE_APPROVED_RESERVED:   "The date the entry was first approved.",
	FIELD_DATE_SYMBOL_CHANGED:      "The date the gene symbol was last changed.",
	FIELD_DATE_NAME_CHANGED:        "The date the gene name was last changed.",
	FIELD_DATE_MODIFIED:            "Date the entry was last modified.",
	FIELD_ENTREZ_ID:                "Entrez gene ID. Found within the \"GENE RESOURCES\" section of the gene symbol report.",
	FIELD_ENSEMBL_GENE_ID:          "Ensembl gene ID. Found within the \"GENE RESOURCES\" section of the gene symbol report.",
	FIELD_VEGA_ID:                  "Vega gene ID. Found within the \"GENE RESOURCES\" section of the gene symbol report.",
	FIELD_UCSC_ID:                  "UCSC gene ID. Found within the \"GENE RESOURCES\" section of the gene symbol report.",
	FIELD_ENA:                      "International Nucleotide Sequence Database Collaboration (GenBank, ENA and DDBJ) accession number(s). Found within the \"NUCLEOTIDE SEQUENCES\" section of the gene symbol report.",
	FIELD_REFSEQ_ACCESSION:         "RefSeq nucleotide accession(s). Found within the \"NUCLEOTIDE SEQUENCES\" section of the gene symbol report.",
	FIELD_CCDS_ID:                  "Consensus CDS ID. Found within the \"NUCLEOTIDE SEQUENCES\" section of the gene symbol report.",
	FIELD_UNIPROT_IDS:              "UniProt protein accession. Found within the \"PROTEIN RESOURCES\" section of the gene symbol report.",
	FIELD_PUBMED_ID:                "Pubmed and Europe Pubmed Central PMID(s).",
	FIELD_MGD_ID:                   "Mouse genome informatics database ID. Found within the \"HOMOLOGS\" section of the gene symbol report.",
	FIELD_RGD_ID:                   "Rat genome database gene ID. Found within the \"HOMOLOGS\" section of the gene symbol report.",
	FIELD_LSDB:                     "The name of the Locus Specific Mutation Database and URL for the gene separated by a | character",
	FIELD_COSMIC:                   "Symbol used within the Catalogue of somatic mutations in cancer for the gene.",
	FIELD_OMIM_ID:                  "Online Mendelian Inheritance in Man (OMIM) ID",
	FIELD_MIRBASE:                  "miRBase ID",
	FIELD_HOMEODB:                  "Homeobox Database ID",
	FIELD_SNORNABASE:               "snoRNABase ID",
	FIELD_BIOPARADIGMS_SLC:         "Symbol used to link to the SLC tables database at bioparadigms.org for the gene",
	FIELD_ORPHANET:                 "Orphanet ID",
	FIELD_PSEUDOGENE_ORG:           "Pseudogene.org",
	FIELD_HORDE_ID:                 "Symbol used within HORDE for the gene",
	FIELD_MEROPS:                   "ID used to link to the MEROPS peptidase database",
	FIELD_IMGT:                     "Symbol used within international ImMunoGeneTics information system",
	FIELD_IUPHAR:                   "The objectId used to link to the IUPHAR/BPS Guide to PHARMACOLOGY database. To link to IUPHAR/BPS Guide to PHARMACOLOGY database only use the number (only use 1 from the result objectId:1)",
	FIELD_KZNF_GENE_CATALOG:        "ID used to link to the Human KZNF Gene Catalog",
	FIELD_MAMIT_TRNADB:             "ID to link to the Mamit-tRNA database",
	FIELD_CD:                       "Symbol used within the Human Cell Differentiation Molecule database for the gene",
	FIELD_LNCRNADB:                 "lncRNA Database ID - Resource now defunct",
	FIELD_ENZYME_ID:                "ENZYME EC accession number",
	FIELD_INTERMEDIATE_FILAMENT_DB: "ID used to link to the Human Intermediate Filament Database",
	FIELD_AGR:                      "The HGNC ID that the Alliance of Genome Resources (AGR) have linked to their record of the gene. Use the HGNC ID to link to the AGR.",
	FIELD_MANE_SELECT:              "NCBI and Ensembl transcript IDs/acessions including the version number for one high-quality representative transcript per protein-coding gene that is well-supported by experimental data and represents the biology of the gene. The IDs are delimited by |.",
}

func (h *HGNC) GetFieldDesc(field Field) string {
	if desc, ok := fieldDesc[field]; ok {
		return desc
	}
	return ""
}
