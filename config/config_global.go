package config

type Global struct {
	DoNotEdit   string `json:"do_not_edit"`
	PackageName string `json:"package_name"`
	ClassName   string `json:"class_name"`
	TplPrefix   string `json:"tpl_prefix"`
	ArgPrefix   string `json:"arg_prefix"`

	Import []*Import `json:"import"`
}

// TODO: 코브라로 변경 예정
func (t *Global) InitDefault() {
	t.DoNotEdit = "// Code generated - DO NOT EDIT.\n// This file is a generated and any changes will be lost.\n"
	t.PackageName = "gen"
	t.ClassName = "Schema"
	t.TplPrefix = "tpl_"
	t.ArgPrefix = "arg_"

	t.Import = []*Import{
		{Alias: "", Path: "fmt"},
		{Alias: ".", Path: "github.com/gokch/ornn/db"},
	}
}

type Import struct {
	Alias string `json:"alias"`
	Path  string `json:"path"`
}
