package config

import "log"

type GenCode struct {
	Imports []*GenImport `json:"imports"`
	Fields  []*GenField  `json:"fields"`
}

func (t *GenCode) ConvFieldType(genType string) string {
	for _, pt_field_type := range t.Fields {
		if pt_field_type.GenType == genType {
			return genType
		}
	}
	log.Fatalf("invalid field type - %s", genType)
	return ""
}

type GenImport struct {
	Alias string `json:"alias"`
	Path  string `json:"path"`
}

type GenField struct {
	GenType string `json:"gen_type"`
}
