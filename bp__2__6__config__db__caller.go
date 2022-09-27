package bp

import "fmt"

type T_BP__config__db__caller struct {
	is_exist_error__caller_sql bool `json:"-"`

	Arrpt_group []*T_BP__config__db__caller__group `json:"group"`
}

type T_BP__config__db__caller__group struct {
	S_group_name string                             `json:"group_name"`
	Arrpt_query  []*T_BP__config__db__caller__query `json:"querys"`
}

type T_BP__config__db__caller__query struct {
	S_query_name string `json:"name"`
	S_comment    string `json:"comment,omitempty"`

	S_sql                   string                                                `json:"sql"`
	Arrpt_select_field_type []*T_BP__config__db__caller__query__select_field_type `json:"select_field_type,omitempty"`
	IS_insert__multi        bool                                                  `json:"insert__multi,omitempty"`
	IS_update__null_ignore  bool                                                  `json:"update__null_ignore,omitempty"`

	S_tmp_err__caller__query  string `json:"-"`
	S_tmp_err__caller__parser string `json:"-"`
}

// select 만 field type이 있는 이유
// select query 는 bp.json 의 schema type 을 통해 타입을 지정할 수 없기 때문에
// 직접 쿼리를 select 를 하고 결과를 추출해 타입에 넣음
// snum, uint 등의 custom type 은 여기서 처리
type T_BP__config__db__caller__query__select_field_type struct {
	S_name           string `json:"name"`
	S_field_type__bp string `json:"type"`
}

//------------------------------------------------------------------------------------------------//
// call

func (t *T_BP__config__db__caller) Init() {
	t.Arrpt_group = make([]*T_BP__config__db__caller__group, 0, 10)
}

func (t *T_BP__config__db__caller) Add_group(_s_group_name string) *T_BP__config__db__caller__group {
	pt_group := &T_BP__config__db__caller__group{}
	pt_group.Init(_s_group_name)
	t.Arrpt_group = append(t.Arrpt_group, pt_group)

	return pt_group
}

func (t *T_BP__config__db__caller) Print__sql_errors() {
	for _, pt_group := range t.Arrpt_group {
		for _, pt_query := range pt_group.Arrpt_query {
			// 에러 없으면 출력 X
			if pt_query.S_tmp_err__caller__query == "" && pt_query.S_tmp_err__caller__parser == "" {
				continue
			}

			// 에러 존재 지점
			fmt.Printf("------------------------------------------------------------------\n")
			fmt.Printf("group_name\n\t%s\n", pt_group.S_group_name)
			fmt.Printf("query_name\n\t%s\n", pt_query.S_query_name)
			fmt.Printf("sql\n\t%s\n", pt_query.S_sql)
			if pt_query.S_tmp_err__caller__parser != "" {
				fmt.Printf("parser error\n\t%s\n", pt_query.S_tmp_err__caller__parser)
			}
			if pt_query.S_tmp_err__caller__query != "" {
				fmt.Printf("query error\n\t%s\n", pt_query.S_tmp_err__caller__query)
			}
		}
	}
}

//------------------------------------------------------------------------------------------------//
// group

func (t *T_BP__config__db__caller__group) Init(_s_group_name string) {
	t.S_group_name = _s_group_name
	t.Arrpt_query = make([]*T_BP__config__db__caller__query, 0, 10)
}

func (t *T_BP__config__db__caller__group) Add_querytion(_s_query_name, _s_sql string) *T_BP__config__db__caller__query {
	pt_query := &T_BP__config__db__caller__query{}
	pt_query.Init(_s_query_name, _s_sql)

	t.Arrpt_query = append(t.Arrpt_query, pt_query)
	return pt_query
}

//------------------------------------------------------------------------------------------------//
// query

func (t *T_BP__config__db__caller__query) Init(_s_query_name, _s_sql string) {
	t.S_query_name = _s_query_name
	t.S_sql = _s_sql
	t.Arrpt_select_field_type = make([]*T_BP__config__db__caller__query__select_field_type, 0, 10)
}

func (t *T_BP__config__db__caller__query) Add_select_field_type(_s_field_name string, _s_field_type__bp string) {
	if t.Arrpt_select_field_type == nil {
		t.Arrpt_select_field_type = make([]*T_BP__config__db__caller__query__select_field_type, 0, 10)
	}
	pt_field_type := &T_BP__config__db__caller__query__select_field_type{}
	pt_field_type.Set(_s_field_name, _s_field_type__bp)
	t.Arrpt_select_field_type = append(t.Arrpt_select_field_type, pt_field_type)
}

func (t *T_BP__config__db__caller__query) Get_select_field_type(_s_field_name string) (s_field_type__bp string) {
	for _, pt := range t.Arrpt_select_field_type {
		if pt.S_name == _s_field_name {
			return pt.S_field_type__bp
		}
	}

	return ""
}

//------------------------------------------------------------------------------------------------//
// select field type

func (t *T_BP__config__db__caller__query__select_field_type) Set(_s_field_name string, _s_field_type__bp string) {
	t.S_name = _s_field_name
	t.S_field_type__bp = _s_field_type__bp
}
