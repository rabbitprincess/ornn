package bp

func (t *T_BP__gen__data) add_group(_pt *T_BP__gen__data__group) {
	if t.arrpt_group == nil {
		t.arrpt_group = make([]*T_BP__gen__data__group, 0, 10)
	}
	t.arrpt_group = append(t.arrpt_group, _pt)
}

//------------------------------------------------------------------------------------------------------------//
type T_BP__gen__data__group struct {
	s_group_name string
	arrpt_query  []*T_BP__gen__data__query
}

func (t *T_BP__gen__data__group) Set(_s_group_name string) {
	if t.arrpt_query == nil {
		t.arrpt_query = make([]*T_BP__gen__data__query, 0, 10)
	}
	t.s_group_name = _s_group_name
}

func (t *T_BP__gen__data__group) add_query(_pt *T_BP__gen__data__query) {
	t.arrpt_query = append(t.arrpt_query, _pt)
}

//------------------------------------------------------------------------------------------------------------//

type TD_N1_query_type int8

const (
	TD_N1_query_type__select TD_N1_query_type = iota + 1
	TD_N1_query_type__insert
	TD_N1_query_type__update
	TD_N1_query_type__delete
)

type T_BP__gen__data__query struct {
	TD_n1_query_type TD_N1_query_type
	s_group_name     string
	s_query_name     string
	s_query          string

	t_tpl T_BP__gen__data____struct
	t_arg T_BP__gen__data____struct
	t_ret T_BP__gen__data____struct

	is_select__single      bool
	is_insert__multi       bool
	is_update__null_ignore bool
}

//------------------------------------------------------------------------------------------------------------//
type T_Pair_data struct {
	s_key   string
	s_value string
}

type T_BP__gen__data____struct struct {
	arrpt_pair []*T_Pair_data
}

func (t *T_BP__gen__data____struct) set_key(_arrs_sql_field_name []string) {
	for _, s_field_name := range _arrs_sql_field_name {
		t.set_key_value(s_field_name, "")
	}
}

func (t *T_BP__gen__data____struct) set_key_value(_s_key string, _s_value__new string) {
	if t.arrpt_pair == nil {
		t.arrpt_pair = make([]*T_Pair_data, 0, 10)
	}

	for _, pt_field_type := range t.arrpt_pair {
		if pt_field_type.s_key == _s_key {
			pt_field_type.s_value = _s_value__new
			return
		}
	}

	pt_field_type := &T_Pair_data{}
	pt_field_type.s_key = _s_key
	pt_field_type.s_value = _s_value__new

	t.arrpt_pair = append(t.arrpt_pair, pt_field_type)
}
