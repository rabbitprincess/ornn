package config

type Global struct {
	DoNotEdit   string `json:"do_not_edit"`
	PackageName string `json:"package_name"`
	ClassName   string `json:"class_name"`

	Import []*Import `json:"import"`
}

func (t *Global) InitDefault() {
	t.DoNotEdit = "// Code generated - DO NOT EDIT.\n// This file is a generated and any changes will be lost.\n"
	t.PackageName = "gen"
	t.ClassName = "Gen"

	t.Import = []*Import{
		{Alias: "", Path: "fmt"},
		{Alias: ".", Path: "github.com/gokch/ornn/db"},
	}
}

type Import struct {
	Alias string `json:"alias"`
	Path  string `json:"path"`
}
