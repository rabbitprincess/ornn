package go_orm_gen

import "log"

type T_BP__config__gen struct {
	Arrpt_lang_type  []*T_BP__config__gen__lang_type  `json:"lang"`
	Arrpt_field_type []*T_BP__config__gen__field_type `json:"field_type"`
}

func (t *T_BP__config__gen) Conv_field_type__bp_to_lang(_s_field_type__bp string, _td_s_lang_name LangType) string {
	for _, pt_field_type := range t.Arrpt_field_type {
		if pt_field_type.S_type_name__bp != _s_field_type__bp {
			continue
		}
		return pt_field_type.Get_conv_field_type__by_lang(_td_s_lang_name)
	}
	log.Fatalf("Conv_field_type - %s", _s_field_type__bp)
	return ""
}
