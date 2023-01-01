package config

type Global struct {
	DoNotEdit   string `cty:"do_not_edit"`
	PackageName string `cty:"package_name"`
	ClassName   string `cty:"class_name"`

	Import []*Import `cty:"import"`
}

func (t *Global) InitDefault() {
	t.DoNotEdit = "// Code generated - DO NOT EDIT.\n// This file is a generated and any changes will be lost.\n"
	t.PackageName = "gen"
	t.ClassName = "Schema"

	t.Import = []*Import{
		{Alias: "", Path: "fmt"},
		{Alias: ".", Path: "github.com/gokch/ornn/db"},
	}
}

type Import struct {
	Alias string `cty:"alias"`
	Path  string `cty:"path"`
}
