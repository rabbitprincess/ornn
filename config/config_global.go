package config

type Global struct {
	DoNotEdit    string `json:"do_not_edit"`
	PackageName  string `json:"package_name"`
	ClassName    string `json:"class_name"`
	InstanceType string `json:"instance_type"`
	InstanceName string `json:"instance_name"`
	StructName   string `json:"struct_name"`
	TplPrefix    string `json:"tpl_prefix"`
	ArgPrefix    string `json:"arg_prefix"`

	Import    []*Import `json:"import"`
	FieldType []string  `json:"field_type"`
}

// TODO: 코브라로 변경 예정
func (t *Global) InitDefault() {
	t.DoNotEdit = "// Code generated - DO NOT EDIT.\n// This file is a generated and any changes will be lost.\n"
	t.PackageName = "gen"
	t.ClassName = "Schema"
	t.InstanceType = "*Job"
	t.InstanceName = "Job"
	t.StructName = "t"
	t.TplPrefix = "tpl_"
	t.ArgPrefix = "arg_"
}

func (t *Global) ConvFieldType(fldType string) string {
	for _, fld := range t.FieldType {
		if fld == fldType {
			return fldType
		}
	}
	return ""
}

type Import struct {
	Alias string `json:"alias"`
	Path  string `json:"path"`
}
