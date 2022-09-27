package bp

type T_SQL__select_as struct {
	S_field_name string
	S_as         string
	S_table_name string
}

func (t *T_SQL__select_as) get_field_name__in_sql() (s_field_name string) {
	s_field_name = ""
	if t.S_table_name != "" {
		s_field_name = t.S_table_name + DEF_s_delimeter__between__table_name__to__field_name
	}

	if t.S_as != "" {
		s_field_name += t.S_as
	} else {
		s_field_name += t.S_field_name
	}

	return s_field_name
}

//--------------------------------------------------------------------------------------------------------//

type T_SQL__table_as struct {
	S_table_name string
	S_as         string
}

func (t *T_SQL__table_as) get_table_name__in_sql() (s_table_name string) {
	if t.S_as != "" {
		s_table_name = t.S_as
	} else {
		s_table_name = t.S_table_name
	}

	return s_table_name
}

//--------------------------------------------------------------------------------------------------------//
type T_SQL__field_value struct {
	S_field_name string
	S_table_name string
	BT_val       []byte
}
