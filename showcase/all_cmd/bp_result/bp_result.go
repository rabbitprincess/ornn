// Code generated - DO NOT EDIT.
// This file is a generated and any changes will be lost.

package bp_db

import (
	"fmt"
	. "module/db"
	"strings"
)

type C_DB struct {
	C_common C_Common
	C_select C_Select
	C_insert C_Insert
	C_update C_Update
	C_delete C_Delete
}

func (t *C_DB) Init(
	_pc_db_job *C_DB_job,
) {
	t.C_common.Init(_pc_db_job)
	t.C_select.Init(_pc_db_job)
	t.C_insert.Init(_pc_db_job)
	t.C_update.Init(_pc_db_job)
	t.C_delete.Init(_pc_db_job)
}

func (t *C_Common) Init(
	_pc_db_job *C_DB_job,
) {
	t.pc_db_job = _pc_db_job
}

type C_Common struct {
	pc_db_job *C_DB_job
}

type T_Common__arg struct {
	U8_seq  int64
	Is_bool int8
	S_str   string
	Bt_bin  []byte
	N1_num  int8
	N2_num  int16
	N4_num  int32
	N8_num  int64
	U1_num  int8
	U2_num  int16
	U4_num  int32
	U8_num  int64
	F4_num  float32
	F8_num  float64
	Bt_snum []byte
}

func (t *C_Common) Arg(
	_arg__u8_seq interface{},
) (
	arrpt_arg []*T_Common__arg,
	err error,
) {
	arri_arg := make([]interface{}, 0, 1)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__u8_seq,
	)...)
	
	s_sql := fmt.Sprintf(
		"select * from test_table where u8_seq = ?",
	)
	pc_ret := t.pc_db_job.Query(
		s_sql,
		arri_arg...,
	)
	defer pc_ret.Close()
	
	arrpt_arg = make([]*T_Common__arg, 0, 100)
	for {
		pt_struct := &T_Common__arg{}
		is_end, err := pc_ret.Row_next(pt_struct)
		if err != nil {
			return nil, err
		}
		if is_end == true {
			break
		}
		arrpt_arg = append(arrpt_arg, pt_struct)
	}
	
	return arrpt_arg, nil
}

type T_Common__arg__limit_one struct {
	U8_seq  int64
	Is_bool int8
	S_str   string
	Bt_bin  []byte
	N1_num  int8
	N2_num  int16
	N4_num  int32
	N8_num  int64
	U1_num  int8
	U2_num  int16
	U4_num  int32
	U8_num  int64
	F4_num  float32
	F8_num  float64
	Bt_snum []byte
}

func (t *C_Common) Arg__limit_one(
	_arg__u8_seq interface{},
) (
	pt_arg__limit_one *T_Common__arg__limit_one,
	err error,
) {
	arri_arg := make([]interface{}, 0, 1)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__u8_seq,
	)...)
	
	s_sql := fmt.Sprintf(
		"select * from test_table where u8_seq = ? limit 1",
	)
	pc_ret := t.pc_db_job.Query(
		s_sql,
		arri_arg...,
	)
	defer pc_ret.Close()
	
	for {
		pt_struct := &T_Common__arg__limit_one{}
		is_end, err := pc_ret.Row_next(pt_struct)
		if err != nil {
			return nil, err
		}
		if is_end == true {
			break
		}
		pt_arg__limit_one = pt_struct
		break
	}
	
	return pt_arg__limit_one, nil
}

func (t *C_Common) Tpl(
	_tpl__table_name string,
	_arg__s_str *string,
	_arg__n1_num *int8,
) (
	nn_last_insert_id int64,
	err error,
) {
	arri_arg := make([]interface{}, 0, 2)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__s_str,
		_arg__n1_num,
	)...)
	
	s_sql := fmt.Sprintf(
		"insert into %s (s_str, n1_num)values(?, ?)",
		_tpl__table_name,
	)
	
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.LastInsertId()
}

func (t *C_Select) Init(
	_pc_db_job *C_DB_job,
) {
	t.pc_db_job = _pc_db_job
}

type C_Select struct {
	pc_db_job *C_DB_job
}

type T_Select__group_by__limit__offset struct {
	U8_seq  int64
	Is_bool int8
	S_str   string
	Bt_bin  []byte
	N1_num  int8
	N2_num  int16
	N4_num  int32
	N8_num  int64
	U1_num  int8
	U2_num  int16
	U4_num  int32
	U8_num  int64
	F4_num  float32
	F8_num  float64
	Bt_snum []byte
}

func (t *C_Select) Group_by__limit__offset(
	_arg__u8_seq interface{},
	_arg__group_by interface{},
	_arg__limit interface{},
	_arg__offset interface{},
) (
	arrpt_group_by__limit__offset []*T_Select__group_by__limit__offset,
	err error,
) {
	arri_arg := make([]interface{}, 0, 4)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__u8_seq,
		_arg__group_by,
		_arg__limit,
		_arg__offset,
	)...)
	
	s_sql := fmt.Sprintf(
		"select * from test_table where u8_seq = ? group by ? limit ? offset ?",
	)
	pc_ret := t.pc_db_job.Query(
		s_sql,
		arri_arg...,
	)
	defer pc_ret.Close()
	
	arrpt_group_by__limit__offset = make([]*T_Select__group_by__limit__offset, 0, 100)
	for {
		pt_struct := &T_Select__group_by__limit__offset{}
		is_end, err := pc_ret.Row_next(pt_struct)
		if err != nil {
			return nil, err
		}
		if is_end == true {
			break
		}
		arrpt_group_by__limit__offset = append(arrpt_group_by__limit__offset, pt_struct)
	}
	
	return arrpt_group_by__limit__offset, nil
}

type T_Select__tpl struct {
	U8_seq  int64
	Is_bool int8
	S_str   string
	Bt_bin  []byte
	N1_num  int8
	N2_num  int16
	N4_num  int32
	N8_num  int64
	U1_num  int8
	U2_num  int16
	U4_num  int32
	U8_num  int64
	F4_num  float32
	F8_num  float64
	Bt_snum []byte
}

func (t *C_Select) Tpl(
	_tpl__table_name string,
) (
	arrpt_tpl []*T_Select__tpl,
	err error,
) {
	arri_arg := make([]interface{}, 0, 0)
	arri_arg = append(arri_arg, I_to_arri()...)
	
	s_sql := fmt.Sprintf(
		"select * from %s",
		_tpl__table_name,
	)
	pc_ret := t.pc_db_job.Query(
		s_sql,
		arri_arg...,
	)
	defer pc_ret.Close()
	
	arrpt_tpl = make([]*T_Select__tpl, 0, 100)
	for {
		pt_struct := &T_Select__tpl{}
		is_end, err := pc_ret.Row_next(pt_struct)
		if err != nil {
			return nil, err
		}
		if is_end == true {
			break
		}
		arrpt_tpl = append(arrpt_tpl, pt_struct)
	}
	
	return arrpt_tpl, nil
}

type T_Select__custom_field struct {
	U8_seq  uint64
	Is_bool bool
	S_str   string
	Bt_bin  []byte
	N1_num  int8
	N2_num  int16
	N4_num  int32
	N8_num  int64
	U1_num  uint8
	U2_num  uint16
	U4_num  uint32
	U8_num  uint64
	F4_num  float32
	F8_num  float64
	Bt_snum T_Col__snum
}

func (t *C_Select) Custom_field() (
	arrpt_custom_field []*T_Select__custom_field,
	err error,
) {
	arri_arg := make([]interface{}, 0, 0)
	arri_arg = append(arri_arg, I_to_arri()...)
	
	s_sql := fmt.Sprintf(
		"select * from test_table",
	)
	pc_ret := t.pc_db_job.Query(
		s_sql,
		arri_arg...,
	)
	defer pc_ret.Close()
	
	arrpt_custom_field = make([]*T_Select__custom_field, 0, 100)
	for {
		pt_struct := &T_Select__custom_field{}
		is_end, err := pc_ret.Row_next(pt_struct)
		if err != nil {
			return nil, err
		}
		if is_end == true {
			break
		}
		arrpt_custom_field = append(arrpt_custom_field, pt_struct)
	}
	
	return arrpt_custom_field, nil
}

func (t *C_Insert) Init(
	_pc_db_job *C_DB_job,
) {
	t.pc_db_job = _pc_db_job
}

type C_Insert struct {
	pc_db_job *C_DB_job
}

func (t *C_Insert) Auto_increment(
	_arg__is_bool *bool,
	_arg__s_str *string,
	_arg__bt_bin *[]byte,
	_arg__n1_num *int8,
	_arg__n2_num *int16,
	_arg__n4_num *int32,
	_arg__n8_num *int64,
	_arg__u1_num *uint8,
	_arg__u2_num *uint16,
	_arg__u4_num *uint32,
	_arg__u8_num *uint64,
	_arg__f4_num *float32,
	_arg__f8_num *float64,
	_arg__bt_snum *T_Col__snum,
) (
	nn_last_insert_id int64,
	err error,
) {
	arri_arg := make([]interface{}, 0, 14)
	arri_arg = append(arri_arg, I_to_arri(
		_arg__is_bool,
		_arg__s_str,
		_arg__bt_bin,
		_arg__n1_num,
		_arg__n2_num,
		_arg__n4_num,
		_arg__n8_num,
		_arg__u1_num,
		_arg__u2_num,
		_arg__u4_num,
		_arg__u8_num,
		_arg__f4_num,
		_arg__f8_num,
		_arg__bt_snum,
	)...)
	
	s_sql := fmt.Sprintf(
		"insert into test_table values(NULL, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
	)
	
	pc_exec := t.pc_db_job.Exec(
		s_sql,
		arri_arg...,
	)
	
	return pc_exec.LastInsertId()
}

func (t *C_Insert) Multi_insert(
	_arg__is_bool []*bool,
	_arg__s_str []*string,
	_arg__n1_num []*int8,
) (
	nn_last_insert_id int64,
	err error,
) {
	n_len_arg := len(_arg__is_bool)
	if n_len_arg == 0 {
		return 0, fmt.Errorf("arg len is zero")
	}
	if n_len_arg != len(_arg__is_bool) || n_len_arg != len(_arg__s_str) || n_len_arg != len(_arg__n1_num) {
		return 0, fmt.Errorf("arg len is not same")
	}
	
	arri_arg := make([]interface{}, 0, n_len_arg*3)
	for i := 0; i < n_len_arg; i++ {
		arri_arg = append(arri_arg, I_to_arri(
			_arg__is_bool[i],
			_arg__s_str[i],
			_arg__n1_num[i],
		)...)
	}
	
	s_sql := fmt.Sprintf(
		"insert into test_table (u8_seq, is_bool, s_str, n1_num) values(NULL, ?, ?, ?)%s",
		strings.Repeat(", (NULL, ?, ?, ?)", n_len_arg-1),
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

func (t *C_Update) Add(
	_arg__u8_num *uint64,
	_arg__u8_seq interface{},
) (
	nn_row_affected int64,
	err error,
) {
	s_sql := fmt.Sprintf(
		"update test_table set u8_num = ? + u8_num where u8_seq = ?",
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

func (t *C_Update) Null_ignore(
	_arg__is_bool *bool,
	_arg__s_str *string,
	_arg__u1_num *uint8,
	_arg__u8_seq interface{},
) (
	nn_row_affected int64,
	err error,
) {
	s_sql := fmt.Sprintf(
		"update test_table set is_bool = ?, s_str = ?, u1_num = ? where u8_seq = ?",
	)
	
	arri_arg := make([]interface{}, 0, 4)
	arrs_sets__removed := make([]string, 0, 4)
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
	if _arg__u1_num == nil {
		arrs_sets__removed = append(arrs_sets__removed, "u1_num")
	} else {
		arri_arg = append(arri_arg, _arg__u1_num)
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

func (t *C_Delete) Init(
	_pc_db_job *C_DB_job,
) {
	t.pc_db_job = _pc_db_job
}

type C_Delete struct {
	pc_db_job *C_DB_job
}

