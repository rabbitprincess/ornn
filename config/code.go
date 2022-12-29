package config

import "log"

type Code struct {
	Import    []*Import `json:"import"`
	FieldType []string  `json:"field_type"`
}

func (t *Code) ConvFieldType(fldType string) string {
	for _, fld := range t.FieldType {
		if fld == fldType {
			return fldType
		}
	}
	log.Fatalf("invalid field type - %s", fldType)
	return ""
}

type Import struct {
	Alias string `json:"alias"`
	Path  string `json:"path"`
}
