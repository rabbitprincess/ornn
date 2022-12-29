package go_orm_gen

import (
	"fmt"
	"strings"
)

type T_BP__config__db__schema struct {
	Arrpt_table []*T_BP__config__db__schema__table `json:"tables"`
}

type T_BP__config__db__schema__table struct {
	S_name      string                             `json:"table_name"`
	Arrpt_field []*T_BP__config__db__schema__field `json:"fields"`
	Arrpt_index []*T_BP__config__db__schema__index `json:"indexs"`
}

type T_BP__config__db__schema__field struct {
	S_name     string `json:"name"`
	S_type__DB string `json:"type_db"`
	S_type__BP string `json:"type_bp"`

	S_comment string `json:"comment,omitempty"`
}

type T_BP__config__db__schema__index struct {
	S_name    string   `json:"name"`
	S_type    string   `json:"type"`
	Arrs_keys []string `json:"keys"`

	S_comment string `json:"comment,omitempty"`
}

//------------------------------------------------------------------------------------------------//
// call

func (t *T_BP__config__db__schema) init() {
	t.Arrpt_table = make([]*T_BP__config__db__schema__table, 0, 10)
}

func (t *T_BP__config__db__schema) add_table(_pt_table *T_BP__config__db__schema__table) {
	t.Arrpt_table = append(t.Arrpt_table, _pt_table)
}

func (t *T_BP__config__db__schema) get_table(_s_table_name string) *T_BP__config__db__schema__table {
	for _, pt := range t.Arrpt_table {
		if _s_table_name == pt.S_name {
			return pt
		}
	}
	return nil
}

func (t *T_BP__config__db__schema) update_schema_table(_pt_schema *T_BP__config__db__schema, _s_table_name___prefix string) error {
	arrpt_table__new := make([]*T_BP__config__db__schema__table, 0, len(_pt_schema.Arrpt_table))

	for _, _pt := range _pt_schema.Arrpt_table {
		var is_exist bool
		// 1. 중복 테이블은 업데이트
		for _, pt := range t.Arrpt_table {
			if pt.S_name == _pt.S_name {
				// 중복 인덱스 업데이트
				pt.update_table_index(_pt)
				// 중복 필드 업데이트
				pt.update_table_field(_pt)
				arrpt_table__new = append(arrpt_table__new, pt)
				is_exist = true
				break
			}
		}
		// 2. 새로운 테이블 추가
		if is_exist == false {
			arrpt_table__new = append(arrpt_table__new, _pt)
		}
		// 3. 기존 테이블은 추가하지 않음 ( 삭제 )
	}

	// prefix 가 있을 시 후처리 - 해당 prefix 를 가지고 있는 테이블만 생성
	if _s_table_name___prefix != "" {
		arrpt_table__new__with_prefix_only := make([]*T_BP__config__db__schema__table, 0, len(_pt_schema.Arrpt_table))
		arrs_table_name__prefix := strings.Split(_s_table_name___prefix, ",")
		for _, pt_table__new := range arrpt_table__new {
			for _, s_table_name__prefix := range arrs_table_name__prefix {
				// table_name__prefix 중 하나랑 매칭될 시
				if strings.HasPrefix(pt_table__new.S_name, s_table_name__prefix) == true {
					arrpt_table__new__with_prefix_only = append(arrpt_table__new__with_prefix_only, pt_table__new)
					break
				}
			}
		}
		arrpt_table__new = arrpt_table__new__with_prefix_only
	}

	t.Arrpt_table = arrpt_table__new
	return nil
}

func (t *T_BP__config__db__schema) get_table_list__have__field_name(_s_field_name string, _arrs_target_table []string) (arrs_table_name__match []string, err error) {
	arrs_table_name__match = make([]string, 0, 10)

	for _, s_table_name := range _arrs_target_table {
		pt_table := t.get_table(s_table_name)
		if pt_table == nil {
			return nil, fmt.Errorf("wrong table_name in sql query, table_name is not exist in schema | table_name : %s", s_table_name)
		}
		if pt_table.get_field(_s_field_name) != nil {
			arrs_table_name__match = append(arrs_table_name__match, s_table_name)
		}
	}
	return arrs_table_name__match, nil
}

//------------------------------------------------------------------------------------------------//
// group

func (t *T_BP__config__db__schema__table) init(_s_table_name string) {
	t.S_name = _s_table_name
	t.Arrpt_field = make([]*T_BP__config__db__schema__field, 0, 10)
	t.Arrpt_index = make([]*T_BP__config__db__schema__index, 0, 10)
}

func (t *T_BP__config__db__schema__table) add_field(_pt_field *T_BP__config__db__schema__field) {
	t.Arrpt_field = append(t.Arrpt_field, _pt_field)
}

func (t *T_BP__config__db__schema__table) get_field(_s_field_name string) *T_BP__config__db__schema__field {
	for _, pt_field := range t.Arrpt_field {
		if _s_field_name == pt_field.S_name {
			return pt_field
		}
	}
	return nil
}

func (t *T_BP__config__db__schema__table) update_table_field(_pt_table *T_BP__config__db__schema__table) error {
	arrpt_field__new := make([]*T_BP__config__db__schema__field, 0, len(_pt_table.Arrpt_field))
	for _, _pt := range _pt_table.Arrpt_field {
		var is_exist bool
		for _, pt := range t.Arrpt_field {
			// 1. 중복 필드는 업데이트
			if pt.S_name == _pt.S_name {
				if pt.S_type__DB == _pt.S_type__DB && pt.S_type__BP != _pt.S_type__BP {
					// type__db 는 같은데 type__bp 만 다르면 업데이트 하지 않음
					arrpt_field__new = append(arrpt_field__new, pt)
				} else {
					// 나머지 케이스에서는 업데이트
					arrpt_field__new = append(arrpt_field__new, _pt)
				}
				is_exist = true
				break
			}
		}
		// 2. 새로운 필드는 추가
		if is_exist == false {
			arrpt_field__new = append(arrpt_field__new, _pt)
		}
		// 3. 기존 필드는 추가하지 않음
	}
	t.Arrpt_field = arrpt_field__new
	return nil
}

func (t *T_BP__config__db__schema__table) update_table_index(_pt_table *T_BP__config__db__schema__table) error {
	t.Arrpt_index = _pt_table.Arrpt_index
	return nil
}

func (t *T_BP__config__db__schema__table) add_index(_pt_index *T_BP__config__db__schema__index) {
	t.Arrpt_index = append(t.Arrpt_index, _pt_index)
}

//------------------------------------------------------------------------------------------------//
// func

func (t *T_BP__config__db__schema__field) set(_s_name string, _s_type__db string, _s_type__bp string) {
	t.S_name = _s_name
	t.S_type__DB = _s_type__db
	t.S_type__BP = _s_type__bp
}

//------------------------------------------------------------------------------------------------//
// idx

func (t *T_BP__config__db__schema__index) set(_s_name string, _s_type string, _arrs_keys []string) {
	t.S_name = _s_name
	t.S_type = _s_type
	t.Arrs_keys = _arrs_keys
}
