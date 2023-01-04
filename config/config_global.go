package config

type Global struct {
	DoNotEdit   string `json:"do_not_edit"`
	PackageName string `json:"package_name"`
	ClassName   string `json:"class_name"`

	Import []*Import `json:"import"`
}

type Import struct {
	Alias string `json:"alias"`
	Path  string `json:"path"`
}
