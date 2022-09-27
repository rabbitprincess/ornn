// Code generated - DO NOT EDIT.
// This file is a generated and any changes will be lost.

package bp_db

import (
	"fmt"
	. "module/db"
	"strings"
)

type C_DB struct {
	C_select C_Select
	C_insert C_Insert
	C_update C_Update
	C_delete C_Delete
}

func (t *C_DB) Init(
	_pc_db_job *C_DB_job,
) {
	t.C_select.Init(_pc_db_job)
	t.C_insert.Init(_pc_db_job)
	t.C_update.Init(_pc_db_job)
	t.C_delete.Init(_pc_db_job)
}

func (t *C_Select) Init(
	_pc_db_job *C_DB_job,
) {
	t.pc_db_job = _pc_db_job
}

type C_Select struct {
	pc_db_job *C_DB_job
}

type T_Select__one struct {
	U8_seq uint64
}

func (t *C_Select) One() (
	arrpt_one []*T_Select__one,
	err error,
) {
	arri_arg := make([]interface{}, 0, 0)
	arri_arg = append(arri_arg, I_to_arri()...)
	
	s_sql := fmt.Sprintf(
		"select u8_seq from bp_gen_test_schema",
	)
	pc_ret := t.pc_db_job.Query(
		s_sql,
		arri_arg...,
	)
	defer pc_ret.Close()
	
	arrpt_one = make([]*T_Select__one, 0, 100)
	for {
		pt_struct := &T_Select__one{}
		is_end, err := pc_ret.Row_next(pt_struct)
		if err != nil {
			return nil, err
		}
		if is_end == true {
			break
		}
		arrpt_one = append(arrpt_one, pt_struct)
	}
	
	return arrpt_one, nil
}

type T_Select__all struct {
	U8_seq   uint64
	Is_bool  int8
	S_str    string
	N1_num   int8
	N2_num   int16
	N4_num   int32
	N8_num   int64
	U1_num   uint8
	U2_num   uint16
	U4_num   uint32
	U8_num   uint64
	Bt_bin   []byte
	Sn_snum  T_Col__snum
	F_num    float32
	D_num    float64
	Js_str   string
	Dt_time  int32
	Dtn_time T_Col__dtn
}

func (t *C_Select) All() (
	arrpt_all []*T_Select__all,
	err error,
) {
	arri_arg := make([]interface{}, 0, 0)
	arri_arg = append(arri_arg, I_to_arri()...)
	
	s_sql := fmt.Sprintf(
		"select * from bp_gen_test_schema",
	)
	pc_ret := t.pc_db_job.Query(
		s_sql,
		arri_arg...,
	)
	defer pc_ret.Close()
	
	arrpt_all = make([]*T_Select__all, 0, 100)
	for {
		pt_struct := &T_Select__all{}
		is_end, err := pc_ret.Row_next(pt_struct)
		if err != nil {
			return nil, err
		}
		if is_end == true {
			break
		}
		arrpt_all = append(arrpt_all, pt_struct)
	}
	
	return arrpt_all, nil
}

type T_Select__tpl__ret_single struct {
	U8_seq   uint64
	Is_bool  int8
	S_str    string
	N1_num   int8
	N2_num   int16
	N4_num   int32
	N8_num   int64
	Bt_bin   []byte
	Sn_snum  T_Col__snum
	F_num    float32
	D_num    float64
	Js_str   string
	Dt_time  int32
	Dtn_time T_Col__dtn
}

func (t *C_Select) Tpl__ret_single(
	_tpl__t1_select string,
) (
	pt_tpl__ret_single *T_Select__tpl__ret_single,
	err error,
) {
	arri_arg := make([]interface{}, 0, 0)
	arri_arg = append(arri_arg, I_to_arri()...)
	
	s_sql := fmt.Sprintf(
		"select u8_seq,is_bool,s_str,n1_num,n2_num,n4_num,n8_num,bt_bin,sn_snum,f_num,d_num,js_str,dt_time,dtn_time from %s where u8_seq=1 or is_bool=0 limit 1",
		_tpl__t1_select,
	)
	pc_ret := t.pc_db_job.Query(
		s_sql,
		arri_arg...,
	)
	defer pc_ret.Close()
	
	for {
		pt_struct := &T_Select__tpl__ret_single{}
		is_end, err := pc_ret.Row_next(pt_struct)
		if err != nil {
			return nil, err
		}
		if is_end == true {
			break
		}
		pt_tpl__ret_single = pt_struct
		break
	}
	
	return pt_tpl__ret_single, nil
}

type T_Select__where_arg struct {
	U8_seq  uint64
	Is_bool int8
}

func (t *C_Select) Where_arg(
	_arg__seq interface{},
	_arg__bool interface{},
) (
	pt_where_arg *T_Select__where_arg,
	err error,
) {
	arri_arg := make([]interface{}, 0, 2)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__seq,
		_arg__bool,
	)...)
	
	s_sql := fmt.Sprintf(
		"select u8_seq,is_bool from bp_gen_test_schema where u8_seq=? and is_bool=? limit 1 offset 4",
	)
	pc_ret := t.pc_db_job.Query(
		s_sql,
		arri_arg...,
	)
	defer pc_ret.Close()
	
	for {
		pt_struct := &T_Select__where_arg{}
		is_end, err := pc_ret.Row_next(pt_struct)
		if err != nil {
			return nil, err
		}
		if is_end == true {
			break
		}
		pt_where_arg = pt_struct
		break
	}
	
	return pt_where_arg, nil
}

type T_Select__where_arg__snum struct {
	Sn_snum T_Col__snum
}

func (t *C_Select) Where_arg__snum(
	_tpl__table_name string,
) (
	arrpt_where_arg__snum []*T_Select__where_arg__snum,
	err error,
) {
	arri_arg := make([]interface{}, 0, 0)
	arri_arg = append(arri_arg, I_to_arri()...)
	
	s_sql := fmt.Sprintf(
		"select sn_snum from %s order by sn_snum",
		_tpl__table_name,
	)
	pc_ret := t.pc_db_job.Query(
		s_sql,
		arri_arg...,
	)
	defer pc_ret.Close()
	
	arrpt_where_arg__snum = make([]*T_Select__where_arg__snum, 0, 100)
	for {
		pt_struct := &T_Select__where_arg__snum{}
		is_end, err := pc_ret.Row_next(pt_struct)
		if err != nil {
			return nil, err
		}
		if is_end == true {
			break
		}
		arrpt_where_arg__snum = append(arrpt_where_arg__snum, pt_struct)
	}
	
	return arrpt_where_arg__snum, nil
}

type T_Select__tpl__n__arg struct {
	U8_seq  uint64
	Is_bool int8
}

func (t *C_Select) Tpl__n__arg(
	_tpl__table_name string,
	_arg__seq interface{},
	_arg__bool interface{},
) (
	arrpt_tpl__n__arg []*T_Select__tpl__n__arg,
	err error,
) {
	arri_arg := make([]interface{}, 0, 2)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__seq,
		_arg__bool,
	)...)
	
	s_sql := fmt.Sprintf(
		"select u8_seq,is_bool from %s where u8_seq>? or is_bool=?",
		_tpl__table_name,
	)
	pc_ret := t.pc_db_job.Query(
		s_sql,
		arri_arg...,
	)
	defer pc_ret.Close()
	
	arrpt_tpl__n__arg = make([]*T_Select__tpl__n__arg, 0, 100)
	for {
		pt_struct := &T_Select__tpl__n__arg{}
		is_end, err := pc_ret.Row_next(pt_struct)
		if err != nil {
			return nil, err
		}
		if is_end == true {
			break
		}
		arrpt_tpl__n__arg = append(arrpt_tpl__n__arg, pt_struct)
	}
	
	return arrpt_tpl__n__arg, nil
}

type T_Select__where_arg__group_by struct {
	U8_seq   uint64
	Is_bool  int8
	S_str    string
	N1_num   int8
	N2_num   int16
	N4_num   int32
	N8_num   int64
	U1_num   uint8
	U2_num   uint16
	U4_num   uint32
	U8_num   uint64
	Bt_bin   []byte
	Sn_snum  T_Col__snum
	F_num    float32
	D_num    float64
	Js_str   string
	Dt_time  int32
	Dtn_time T_Col__dtn
}

func (t *C_Select) Where_arg__group_by(
	_arg__u8_seq interface{},
	_arg__is_bool interface{},
	_arg__u8_seq__group_by interface{},
) (
	arrpt_where_arg__group_by []*T_Select__where_arg__group_by,
	err error,
) {
	arri_arg := make([]interface{}, 0, 3)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__u8_seq,
		_arg__is_bool,
		_arg__u8_seq__group_by,
	)...)
	
	s_sql := fmt.Sprintf(
		"select * from bp_gen_test_schema where u8_seq=? or is_bool=? group by ?",
	)
	pc_ret := t.pc_db_job.Query(
		s_sql,
		arri_arg...,
	)
	defer pc_ret.Close()
	
	arrpt_where_arg__group_by = make([]*T_Select__where_arg__group_by, 0, 100)
	for {
		pt_struct := &T_Select__where_arg__group_by{}
		is_end, err := pc_ret.Row_next(pt_struct)
		if err != nil {
			return nil, err
		}
		if is_end == true {
			break
		}
		arrpt_where_arg__group_by = append(arrpt_where_arg__group_by, pt_struct)
	}
	
	return arrpt_where_arg__group_by, nil
}

func (t *C_Insert) Init(
	_pc_db_job *C_DB_job,
) {
	t.pc_db_job = _pc_db_job
}

type C_Insert struct {
	pc_db_job *C_DB_job
}

func (t *C_Insert) Arg__field_name_all(
	_arg__s_str *string,
	_arg__bt_bin *[]byte,
	_arg__sn_snum *T_Col__snum,
	_arg__js_str *string,
) (
	nn_last_insert_id int64,
	err error,
) {
	arri_arg := make([]interface{}, 0, 4)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__s_str,
		_arg__bt_bin,
		_arg__sn_snum,
		_arg__js_str,
	)...)
	
	s_sql := fmt.Sprintf(
		"insert into bp_gen_test_schema (u8_seq,is_bool,s_str,n1_num,n2_num,n4_num,n8_num,u1_num,u2_num,u4_num,u8_num,bt_bin,sn_snum,f_num,d_num,js_str,dt_time,dtn_time) values (0, 0, ?, 0, 1, 2, 3, 0, 1, 2, 3, ?, ?, 1.0, 1.233, ?, 176653424, 43223545)",
	)
	
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.LastInsertId()
}

func (t *C_Insert) Tpl__field_name_all__arg(
	_tpl__bp_gen_test_schema string,
	_tpl__u8_seq string,
	_arg__u8_seq *uint64,
	_arg__is_bool *int8,
	_arg__s_str *string,
	_arg__n1_num *int8,
	_arg__n2_num *int16,
	_arg__n4_num *int32,
	_arg__n8_num *int64,
	_arg__bt_bin *[]byte,
	_arg__sn_snum *T_Col__snum,
	_arg__f_num *float32,
	_arg__d_num *float64,
	_arg__js_str *string,
	_arg__dt_time *uint32,
	_arg__dtn_time *T_Col__dtn,
) (
	nn_last_insert_id int64,
	err error,
) {
	arri_arg := make([]interface{}, 0, 14)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__u8_seq,
		_arg__is_bool,
		_arg__s_str,
		_arg__n1_num,
		_arg__n2_num,
		_arg__n4_num,
		_arg__n8_num,
		_arg__bt_bin,
		_arg__sn_snum,
		_arg__f_num,
		_arg__d_num,
		_arg__js_str,
		_arg__dt_time,
		_arg__dtn_time,
	)...)
	
	s_sql := fmt.Sprintf(
		"insert into %s (%s,is_bool,s_str,n1_num,n2_num,n4_num,n8_num,bt_bin,sn_snum,f_num,d_num,js_str,dt_time,dtn_time) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		_tpl__bp_gen_test_schema,
		_tpl__u8_seq,
	)
	
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.LastInsertId()
}

func (t *C_Insert) Arg__field_star(
	_arg__u8_seq *uint64,
	_arg__is_bool *int8,
	_arg__s_str *string,
	_arg__n1_num *int8,
	_arg__n2_num *int16,
	_arg__n4_num *int32,
	_arg__n8_num *int64,
	_arg__u1_num *uint8,
	_arg__u2_num *uint16,
	_arg__u4_num *uint32,
	_arg__u8_num *uint64,
	_arg__bt_bin *[]byte,
	_arg__sn_snum *T_Col__snum,
	_arg__f_num *float32,
	_arg__d_num *float64,
	_arg__js_str *string,
	_arg__dt_time *uint32,
	_arg__dtn_time *T_Col__dtn,
) (
	nn_last_insert_id int64,
	err error,
) {
	arri_arg := make([]interface{}, 0, 18)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__u8_seq,
		_arg__is_bool,
		_arg__s_str,
		_arg__n1_num,
		_arg__n2_num,
		_arg__n4_num,
		_arg__n8_num,
		_arg__u1_num,
		_arg__u2_num,
		_arg__u4_num,
		_arg__u8_num,
		_arg__bt_bin,
		_arg__sn_snum,
		_arg__f_num,
		_arg__d_num,
		_arg__js_str,
		_arg__dt_time,
		_arg__dtn_time,
	)...)
	
	s_sql := fmt.Sprintf(
		"insert into bp_gen_test_schema values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	)
	
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.LastInsertId()
}

func (t *C_Insert) Arg__field_one__use_set_keyword(
	_arg__u8_seq *uint64,
) (
	nn_last_insert_id int64,
	err error,
) {
	arri_arg := make([]interface{}, 0, 1)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__u8_seq,
	)...)
	
	s_sql := fmt.Sprintf(
		"insert into bp_gen_test_schema set u8_seq = ?",
	)
	
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.LastInsertId()
}

func (t *C_Insert) All_field__multi_insert(
	_arg__u8_seq []*uint64,
	_arg__is_bool []*int8,
	_arg__s_str []*string,
	_arg__n1_num []*int8,
	_arg__n2_num []*int16,
	_arg__n4_num []*int32,
	_arg__n8_num []*int64,
	_arg__u1_num []*uint8,
	_arg__u2_num []*uint16,
	_arg__u4_num []*uint32,
	_arg__u8_num []*uint64,
	_arg__bt_bin []*[]byte,
	_arg__sn_snum []*T_Col__snum,
	_arg__f_num []*float32,
	_arg__d_num []*float64,
	_arg__js_str []*string,
	_arg__dt_time []*uint32,
	_arg__dtn_time []*T_Col__dtn,
) (
	nn_last_insert_id int64,
	err error,
) {
	n_len_arg := len(_arg__u8_seq)
	if n_len_arg == 0 {
		return 0, fmt.Errorf("arg len is zero")
	}
	if n_len_arg != len(_arg__u8_seq) || n_len_arg != len(_arg__is_bool) || n_len_arg != len(_arg__s_str) || n_len_arg != len(_arg__n1_num) || n_len_arg != len(_arg__n2_num) || n_len_arg != len(_arg__n4_num) || n_len_arg != len(_arg__n8_num) || n_len_arg != len(_arg__u1_num) || n_len_arg != len(_arg__u2_num) || n_len_arg != len(_arg__u4_num) || n_len_arg != len(_arg__u8_num) || n_len_arg != len(_arg__bt_bin) || n_len_arg != len(_arg__sn_snum) || n_len_arg != len(_arg__f_num) || n_len_arg != len(_arg__d_num) || n_len_arg != len(_arg__js_str) || n_len_arg != len(_arg__dt_time) || n_len_arg != len(_arg__dtn_time) {
		return 0, fmt.Errorf("arg len is not same")
	}
	
	arri_arg := make([]interface{}, 0, n_len_arg*18)
	for i := 0; i < n_len_arg; i++ {
		arri_arg = append(arri_arg, I_to_arri(
			_arg__u8_seq[i],
			_arg__is_bool[i],
			_arg__s_str[i],
			_arg__n1_num[i],
			_arg__n2_num[i],
			_arg__n4_num[i],
			_arg__n8_num[i],
			_arg__u1_num[i],
			_arg__u2_num[i],
			_arg__u4_num[i],
			_arg__u8_num[i],
			_arg__bt_bin[i],
			_arg__sn_snum[i],
			_arg__f_num[i],
			_arg__d_num[i],
			_arg__js_str[i],
			_arg__dt_time[i],
			_arg__dtn_time[i],
		)...)
	}
	
	s_sql := fmt.Sprintf(
		"insert into bp_gen_test_schema values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)%s",
		strings.Repeat(", (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", n_len_arg-1),
	)
	
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.LastInsertId()
}

func (t *C_Insert) Tpl__all_field__multi_insert(
	_tpl__t1_update string,
	_arg__u8_seq []*uint64,
	_arg__is_bool []*int8,
	_arg__s_str []*string,
	_arg__n1_num []*int8,
	_arg__n2_num []*int16,
	_arg__n4_num []*int32,
	_arg__n8_num []*int64,
	_arg__u1_num []*uint8,
	_arg__u2_num []*uint16,
	_arg__u4_num []*uint32,
	_arg__u8_num []*uint64,
	_arg__bt_bin []*[]byte,
	_arg__sn_snum []*T_Col__snum,
	_arg__f_num []*float32,
	_arg__d_num []*float64,
	_arg__js_str []*string,
	_arg__dt_time []*uint32,
	_arg__dtn_time []*T_Col__dtn,
) (
	nn_last_insert_id int64,
	err error,
) {
	n_len_arg := len(_arg__u8_seq)
	if n_len_arg == 0 {
		return 0, fmt.Errorf("arg len is zero")
	}
	if n_len_arg != len(_arg__u8_seq) || n_len_arg != len(_arg__is_bool) || n_len_arg != len(_arg__s_str) || n_len_arg != len(_arg__n1_num) || n_len_arg != len(_arg__n2_num) || n_len_arg != len(_arg__n4_num) || n_len_arg != len(_arg__n8_num) || n_len_arg != len(_arg__u1_num) || n_len_arg != len(_arg__u2_num) || n_len_arg != len(_arg__u4_num) || n_len_arg != len(_arg__u8_num) || n_len_arg != len(_arg__bt_bin) || n_len_arg != len(_arg__sn_snum) || n_len_arg != len(_arg__f_num) || n_len_arg != len(_arg__d_num) || n_len_arg != len(_arg__js_str) || n_len_arg != len(_arg__dt_time) || n_len_arg != len(_arg__dtn_time) {
		return 0, fmt.Errorf("arg len is not same")
	}
	
	arri_arg := make([]interface{}, 0, n_len_arg*18)
	for i := 0; i < n_len_arg; i++ {
		arri_arg = append(arri_arg, I_to_arri(
			_arg__u8_seq[i],
			_arg__is_bool[i],
			_arg__s_str[i],
			_arg__n1_num[i],
			_arg__n2_num[i],
			_arg__n4_num[i],
			_arg__n8_num[i],
			_arg__u1_num[i],
			_arg__u2_num[i],
			_arg__u4_num[i],
			_arg__u8_num[i],
			_arg__bt_bin[i],
			_arg__sn_snum[i],
			_arg__f_num[i],
			_arg__d_num[i],
			_arg__js_str[i],
			_arg__dt_time[i],
			_arg__dtn_time[i],
		)...)
	}
	
	s_sql := fmt.Sprintf(
		"insert into %s values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)%s",
		_tpl__t1_update,
		strings.Repeat(", (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", n_len_arg-1),
	)
	
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.LastInsertId()
}

func (t *C_Insert) Tpl__all_field__seq_null__multi_insert(
	_tpl__t1_update string,
	_arg__is_bool []*int8,
	_arg__s_str []*string,
	_arg__n1_num []*int8,
	_arg__n2_num []*int16,
	_arg__n4_num []*int32,
	_arg__n8_num []*int64,
	_arg__u1_num []*uint8,
	_arg__u2_num []*uint16,
	_arg__u4_num []*uint32,
	_arg__u8_num []*uint64,
	_arg__bt_bin []*[]byte,
	_arg__sn_snum []*T_Col__snum,
	_arg__f_num []*float32,
	_arg__d_num []*float64,
	_arg__js_str []*string,
	_arg__dt_time []*uint32,
	_arg__dtn_time []*T_Col__dtn,
) (
	nn_last_insert_id int64,
	err error,
) {
	n_len_arg := len(_arg__is_bool)
	if n_len_arg == 0 {
		return 0, fmt.Errorf("arg len is zero")
	}
	if n_len_arg != len(_arg__is_bool) || n_len_arg != len(_arg__s_str) || n_len_arg != len(_arg__n1_num) || n_len_arg != len(_arg__n2_num) || n_len_arg != len(_arg__n4_num) || n_len_arg != len(_arg__n8_num) || n_len_arg != len(_arg__u1_num) || n_len_arg != len(_arg__u2_num) || n_len_arg != len(_arg__u4_num) || n_len_arg != len(_arg__u8_num) || n_len_arg != len(_arg__bt_bin) || n_len_arg != len(_arg__sn_snum) || n_len_arg != len(_arg__f_num) || n_len_arg != len(_arg__d_num) || n_len_arg != len(_arg__js_str) || n_len_arg != len(_arg__dt_time) || n_len_arg != len(_arg__dtn_time) {
		return 0, fmt.Errorf("arg len is not same")
	}
	
	arri_arg := make([]interface{}, 0, n_len_arg*17)
	for i := 0; i < n_len_arg; i++ {
		arri_arg = append(arri_arg, I_to_arri(
			_arg__is_bool[i],
			_arg__s_str[i],
			_arg__n1_num[i],
			_arg__n2_num[i],
			_arg__n4_num[i],
			_arg__n8_num[i],
			_arg__u1_num[i],
			_arg__u2_num[i],
			_arg__u4_num[i],
			_arg__u8_num[i],
			_arg__bt_bin[i],
			_arg__sn_snum[i],
			_arg__f_num[i],
			_arg__d_num[i],
			_arg__js_str[i],
			_arg__dt_time[i],
			_arg__dtn_time[i],
		)...)
	}
	
	s_sql := fmt.Sprintf(
		"insert into %s values (NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)%s",
		_tpl__t1_update,
		strings.Repeat(", (NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", n_len_arg-1),
	)
	
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.LastInsertId()
}

func (t *C_Insert) Tpl__snum(
	_tpl__t1_update string,
	_arg__sn_snum *T_Col__snum,
) (
	nn_last_insert_id int64,
	err error,
) {
	arri_arg := make([]interface{}, 0, 1)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__sn_snum,
	)...)
	
	s_sql := fmt.Sprintf(
		"insert into %s (sn_snum) values (?)",
		_tpl__t1_update,
	)
	
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.LastInsertId()
}

func (t *C_Insert) Tpl__snum__multi_insert(
	_tpl__t1_update string,
	_arg__sn_snum []*T_Col__snum,
) (
	nn_last_insert_id int64,
	err error,
) {
	n_len_arg := len(_arg__sn_snum)
	if n_len_arg == 0 {
		return 0, fmt.Errorf("arg len is zero")
	}
	if n_len_arg != len(_arg__sn_snum) {
		return 0, fmt.Errorf("arg len is not same")
	}
	
	arri_arg := make([]interface{}, 0, n_len_arg*1)
	for i := 0; i < n_len_arg; i++ {
		arri_arg = append(arri_arg, I_to_arri(
			_arg__sn_snum[i],
		)...)
	}
	
	s_sql := fmt.Sprintf(
		"insert into %s (sn_snum) values (?)%s",
		_tpl__t1_update,
		strings.Repeat(", (?)", n_len_arg-1),
	)
	
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.LastInsertId()
}

func (t *C_Update) Init(
	_pc_db_job *C_DB_job,
) {
	t.pc_db_job = _pc_db_job
}

type C_Update struct {
	pc_db_job *C_DB_job
}

func (t *C_Update) Arg(
	_arg__is_bool *int8,
	_arg__s_str *string,
	_arg__u8_seq interface{},
) (
	nn_row_affected int64,
	err error,
) {
	s_sql := fmt.Sprintf(
		"update bp_gen_test_schema set is_bool = ?, s_str = ? where u8_seq = ?",
	)
	arri_arg := make([]interface{}, 0, 3)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__is_bool,
		_arg__s_str,
		_arg__u8_seq,
	)...)
	
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.RowsAffected()
}

func (t *C_Update) Tpl__arg(
	_tpl__t1_update string,
	_arg__is_bool *int8,
	_arg__s_str *string,
	_arg__u8_seq interface{},
) (
	nn_row_affected int64,
	err error,
) {
	s_sql := fmt.Sprintf(
		"update %s set is_bool = ?, s_str = ? where u8_seq = ?",
		_tpl__t1_update,
	)
	arri_arg := make([]interface{}, 0, 3)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__is_bool,
		_arg__s_str,
		_arg__u8_seq,
	)...)
	
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.RowsAffected()
}

func (t *C_Update) Tpl__arg__without_null(
	_tpl__t1_update string,
	_arg__u1_num *uint8,
	_arg__is_bool *int8,
	_arg__s_str *string,
	_arg__u8_seq interface{},
) (
	nn_row_affected int64,
	err error,
) {
	s_sql := fmt.Sprintf(
		"update %s set u1_num = ?, is_bool = ?, s_str = ? where u8_seq = ?",
		_tpl__t1_update,
	)
	
	arri_arg := make([]interface{}, 0, 4)
	arrs_sets__removed := make([]string, 0, 4)
	if _arg__u1_num == nil {
		arrs_sets__removed = append(arrs_sets__removed, "u1_num")
	} else {
		arri_arg = append(arri_arg, _arg__u1_num)
	}
	if _arg__is_bool == nil {
		arrs_sets__removed = append(arrs_sets__removed, "is_bool")
	} else {
		arri_arg = append(arri_arg, _arg__is_bool)
	}
	if _arg__s_str == nil {
		arrs_sets__removed = append(arrs_sets__removed, "s_str")
	} else {
		arri_arg = append(arri_arg, _arg__s_str)
	}
	if _arg__u8_seq == nil {
		arrs_sets__removed = append(arrs_sets__removed, "u8_seq")
	} else {
		arri_arg = append(arri_arg, _arg__u8_seq)
	}
	
	if len(arrs_sets__removed) != 0 {
		s_sql, _ = SQL_remove__update_set_field__null(s_sql, arrs_sets__removed)
	}
	
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.RowsAffected()
}

func (t *C_Update) Tpl__add__1(
	_tpl__t1_update string,
	_arg__u8_seq interface{},
) (
	nn_row_affected int64,
	err error,
) {
	s_sql := fmt.Sprintf(
		"update %s set u8_num = u8_num + 1 where u8_seq = ?",
		_tpl__t1_update,
	)
	arri_arg := make([]interface{}, 0, 1)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__u8_seq,
	)...)
	
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.RowsAffected()
}

func (t *C_Update) Tpl__arg__add__left(
	_tpl__t1_update string,
	_arg__u8_num *uint64,
	_arg__u8_seq interface{},
) (
	nn_row_affected int64,
	err error,
) {
	s_sql := fmt.Sprintf(
		"update %s set u8_num = ? + u8_num where u8_seq = ?",
		_tpl__t1_update,
	)
	arri_arg := make([]interface{}, 0, 2)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__u8_num,
		_arg__u8_seq,
	)...)
	
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.RowsAffected()
}

func (t *C_Update) Tpl__arg__add__right(
	_tpl__t1_update string,
	_arg__u8_num *uint64,
	_arg__u8_seq interface{},
) (
	nn_row_affected int64,
	err error,
) {
	s_sql := fmt.Sprintf(
		"update %s set u8_num = u8_num + ? where u8_seq = ?",
		_tpl__t1_update,
	)
	arri_arg := make([]interface{}, 0, 2)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__u8_num,
		_arg__u8_seq,
	)...)
	
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.RowsAffected()
}

func (t *C_Update) Tpl__arg__add__both(
	_tpl__t1_update string,
	_arg__u8_num *uint64,
	_arg__u4_num *uint32,
	_arg__u8_seq interface{},
) (
	nn_row_affected int64,
	err error,
) {
	s_sql := fmt.Sprintf(
		"update %s set u8_num = u8_num + ?, u4_num = ? + u8_num where u8_seq = ?",
		_tpl__t1_update,
	)
	arri_arg := make([]interface{}, 0, 3)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__u8_num,
		_arg__u4_num,
		_arg__u8_seq,
	)...)
	
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.RowsAffected()
}

func (t *C_Delete) Init(
	_pc_db_job *C_DB_job,
) {
	t.pc_db_job = _pc_db_job
}

type C_Delete struct {
	pc_db_job *C_DB_job
}

func (t *C_Delete) Arg(
	_arg__u8_seq interface{},
) (
	nn_row_affected int64,
	err error,
) {
	arri_arg := make([]interface{}, 0, 1)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__u8_seq,
	)...)
	
	s_sql := fmt.Sprintf(
		"delete from bp_gen_test_schema where u8_seq = ?",
	)
			
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.RowsAffected()
}

func (t *C_Delete) Tpl__arg(
	_tpl__t1_delete string,
	_arg__u8_seq interface{},
) (
	nn_row_affected int64,
	err error,
) {
	arri_arg := make([]interface{}, 0, 1)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__u8_seq,
	)...)
	
	s_sql := fmt.Sprintf(
		"delete from %s where u8_seq = ?",
		_tpl__t1_delete,
	)
			
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.RowsAffected()
}

func (t *C_Delete) Tpl__all(
	_tpl__t1_delete string,
) (
	nn_row_affected int64,
	err error,
) {
	arri_arg := make([]interface{}, 0, 0)
	arri_arg = append(arri_arg, I_to_arri()...)
	
	s_sql := fmt.Sprintf(
		"delete from %s",
		_tpl__t1_delete,
	)
			
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.RowsAffected()
}

