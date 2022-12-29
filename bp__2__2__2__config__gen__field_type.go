package go_orm_gen

type T_BP__config__gen__field_type struct {
	S_type_name__bp string                                    `json:"type_name_bp"`
	Arrpt_lang      []*T_BP__config__gen__field_type__by_lang `json:"lang_type"`
}

type T_BP__config__gen__field_type__by_lang struct {
	TD_s_lang_name    LangType `json:"lang_name"`
	S_type_name__lang string   `json:"type_name"`
}

func (t *T_BP__config__gen__field_type) Get_conv_field_type__by_lang(_td_s_lang_name LangType) string {
	for _, pt_lang_type := range t.Arrpt_lang {
		if pt_lang_type.TD_s_lang_name == _td_s_lang_name {
			return pt_lang_type.S_type_name__lang
		}
	}
	return ""
}
