package bp

type T_BP__config__gen__lang_type struct {
	TD_s_lang_name TD_S_lang_name                          `json:"lang_name"`
	Arrpt_import   []*T_BP__config__gen__lang_type__import `json:"import"`
}

type T_BP__config__gen__lang_type__import struct {
	S_name string `json:"name"`
	S_path string `json:"path"`
}
