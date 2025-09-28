package hgnc_go

import (
	"encoding/json"
	"io"
)

// Record represents a single row of data from the HGNC data file.
type Record struct {
	data map[Field]string
}

// ToMap returns the internal map representation of the Record.
func (r *Record) ToMap() map[Field]string {
	return r.data
}

// ToStrMap returns the internal map representation of the Record, with Field keys
func (r *Record) ToStrMap() map[string]string {
	result := make(map[string]string)
	for k, v := range r.data {
		result[string(k)] = v
	}
	return result
}

// Dump writes the Record to the given writer as JSON.
func (r *Record) Dump(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(r.data)
}

// Dumps returns the Record as a JSON string.
func (r *Record) Dumps() (string, error) {
	jsonBytes, err := json.Marshal(r.data)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

// Get returns the value of the given field in the Record.
func (r *Record) Get(field Field) string {
	return r.data[field]
}

// -------------------------------------------------
// Accessors for each field in the Record struct:

func (r *Record) HgncID() string {
	return r.data[FIELD_HGNC_ID]
}

func (r *Record) Symbol() string {
	return r.data[FIELD_SYMBOL]
}

func (r *Record) EntrezID() string {
	return r.data[FIELD_ENTREZ_ID]
}

func (r *Record) EnsemblGeneID() string {
	return r.data[FIELD_ENSEMBL_GENE_ID]
}

func (r *Record) UcscID() string {
	return r.data[FIELD_UCSC_ID]
}

func (r *Record) RefseqAccession() string {
	return r.data[FIELD_REFSEQ_ACCESSION]
}

func (r *Record) OmimID() string {
	return r.data[FIELD_OMIM_ID]
}

func (r *Record) Name() string {
	return r.data[FIELD_NAME]
}

func (r *Record) LocusGroup() string {
	return r.data[FIELD_LOCUS_GROUP]
}

func (r *Record) LocusType() string {
	return r.data[FIELD_LOCUS_TYPE]
}

func (r *Record) Status() string {
	return r.data[FIELD_STATUS]
}

func (r *Record) Location() string {
	return r.data[FIELD_LOCATION]
}

func (r *Record) LocationSortable() string {
	return r.data[FIELD_LOCATION_SORTABLE]
}

func (r *Record) AliasSymbol() string {
	return r.data[FIELD_ALIAS_SYMBOL]
}

func (r *Record) AliasName() string {
	return r.data[FIELD_ALIAS_NAME]
}

func (r *Record) PrevSymbol() string {
	return r.data[FIELD_PREV_SYMBOL]
}

func (r *Record) PrevName() string {
	return r.data[FIELD_PREV_NAME]
}

func (r *Record) GeneFamily() string {
	return r.data[FIELD_GENE_FAMILY]
}

func (r *Record) GeneFamilyID() string {
	return r.data[FIELD_GENE_FAMILY_ID]
}

func (r *Record) DateApprovedReserved() string {
	return r.data[FIELD_DATE_APPROVED_RESERVED]
}

func (r *Record) DateSymbolChanged() string {
	return r.data[FIELD_DATE_SYMBOL_CHANGED]
}

func (r *Record) DateNameChanged() string {
	return r.data[FIELD_DATE_NAME_CHANGED]
}

func (r *Record) DateModified() string {
	return r.data[FIELD_DATE_MODIFIED]
}

func (r *Record) VegaID() string {
	return r.data[FIELD_VEGA_ID]
}

func (r *Record) ENA() string {
	return r.data[FIELD_ENA]
}

func (r *Record) CcdsID() string {
	return r.data[FIELD_CCDS_ID]
}

func (r *Record) UniprotIDs() string {
	return r.data[FIELD_UNIPROT_IDS]
}

func (r *Record) PubmedID() string {
	return r.data[FIELD_PUBMED_ID]
}

func (r *Record) MgdID() string {
	return r.data[FIELD_MGD_ID]
}

func (r *Record) RgdID() string {
	return r.data[FIELD_RGD_ID]
}

func (r *Record) LSDB() string {
	return r.data[FIELD_LSDB]
}

func (r *Record) Cosmic() string {
	return r.data[FIELD_COSMIC]
}

func (r *Record) Mirbase() string {
	return r.data[FIELD_MIRBASE]
}

func (r *Record) HomeoDB() string {
	return r.data[FIELD_HOMEODB]
}

func (r *Record) SnoRNABase() string {
	return r.data[FIELD_SNORNABASE]
}

func (r *Record) BioparadigmsSLC() string {
	return r.data[FIELD_BIOPARADIGMS_SLC]
}

func (r *Record) Orphanet() string {
	return r.data[FIELD_ORPHANET]
}

func (r *Record) PseudogeneOrg() string {
	return r.data[FIELD_PSEUDOGENE_ORG]
}

func (r *Record) HordeID() string {
	return r.data[FIELD_HORDE_ID]
}

func (r *Record) MEROPS() string {
	return r.data[FIELD_MEROPS]
}

func (r *Record) IMGT() string {
	return r.data[FIELD_IMGT]
}

func (r *Record) IUPHAR() string {
	return r.data[FIELD_IUPHAR]
}

func (r *Record) KZNFGeneCatalog() string {
	return r.data[FIELD_KZNF_GENE_CATALOG]
}

func (r *Record) MamitTRNADB() string {
	return r.data[FIELD_MAMIT_TRNADB]
}

func (r *Record) CD() string {
	return r.data[FIELD_CD]
}

func (r *Record) LncRNADB() string {
	return r.data[FIELD_LNCRNADB]
}

func (r *Record) EnzymeID() string {
	return r.data[FIELD_ENZYME_ID]
}

func (r *Record) IntermediateFilamentDB() string {
	return r.data[FIELD_INTERMEDIATE_FILAMENT_DB]
}

func (r *Record) AGR() string {
	return r.data[FIELD_AGR]
}

func (r *Record) ManeSelect() string {
	return r.data[FIELD_MANE_SELECT]
}
