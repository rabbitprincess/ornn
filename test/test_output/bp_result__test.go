package go_orm_gen_db

import (
	"database/sql"
	"fmt"
	"module/db"
	"module/db/db_mysql/test_env"
	"testing"
)

const (
	DEF_s_test_db_name    = "dev_bp_test"
	DEF_s_test_table_name = "bp_gen_test_schema"
)

func Test__select__one(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, true)

	pc_dbms := &C_DB{}
	pc_dbms.Init(pc_job)

	arrpt, err := pc_dbms.C_select.One()
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
	for _, pt := range arrpt {
		fmt.Printf("%+v\n", pt)
	}
	pc_job.Commit()
}

func Test__select__all(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, true)

	pc_dbms := &C_DB{}
	pc_dbms.Init(pc_job)

	arrpt, err := pc_dbms.C_select.All()
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
	for _, pt := range arrpt {
		fmt.Printf("%+v\n", pt)
	}
	pc_job.Commit()
}

func Test__select__single__arg(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, true)

	pc_dbms := &C_DB{}
	pc_dbms.Init(pc_job)

	pt, err := pc_dbms.C_select.Where_arg(2, false)
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
	fmt.Println(pt)

	pc_job.Commit()
}

func Test__select__single__tpl(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, true)

	pc_dbms := &C_DB{}
	pc_dbms.Init(pc_job)

	pt, err := pc_dbms.C_select.Tpl__ret_single(DEF_s_test_table_name)
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
	fmt.Println(pt)
	pc_job.Commit()
}

func Test__select__tpl__n__arg__where(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, true)
	pc_dbms := &C_DB{}
	pc_dbms.Init(pc_job)

	arrpt, err := pc_dbms.C_select.Tpl__n__arg(DEF_s_test_table_name, 10, 1)
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
	for _, pt := range arrpt {
		fmt.Printf("%+v\n", pt)
	}
	pc_job.Commit()
}

func Test__select__tpl__n__arg__order_by_snum(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, true)
	pc_dbms := &C_DB{}
	pc_dbms.Init(pc_job)

	arrpt, err := pc_dbms.C_select.Where_arg__snum(DEF_s_test_table_name)
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
	for _, pt := range arrpt {
		fmt.Printf("%+v\n", pt)
	}
	pc_job.Commit()
}

// /*
// func Test__insert__tpl__all_field__multi_insert(_t *testing.T) {
// 	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, false)
// 	pc_dbms := &C_DB{}
// 	pc_dbms.Init(pc_job)

// 	arr_u8_seq := make([]*uint64, 0, 5)
// 	{
// 		var u8_seq_1 uint64 = 1
// 		arr_u8_seq = append(arr_u8_seq, &u8_seq_1)
// 		var u8_seq_2 uint64 = 2
// 		arr_u8_seq = append(arr_u8_seq, &u8_seq_2)
// 		var u8_seq_3 uint64 = 3
// 		arr_u8_seq = append(arr_u8_seq, &u8_seq_3)
// 		var u8_seq_4 uint64 = 4
// 		arr_u8_seq = append(arr_u8_seq, &u8_seq_4)
// 		var u8_seq_5 uint64 = 5
// 		arr_u8_seq = append(arr_u8_seq, &u8_seq_5)
// 	}
// 	arr_n1_bool := []int8{0, 1, 0, 1, 0}
// 	arr_s_str := []string{"a", "b", "c", "d", "e"}
// 	arr_n1_num := []int8{math.MaxInt8, -math.MaxInt8, 3, 4, 5}
// 	arr_n2_num := []int16{math.MaxInt16, -math.MaxInt16, 3, 4, 5}
// 	arr_n4_num := []int32{math.MaxInt32, -math.MaxInt32, 3, 4, 5}
// 	arr_n8_num := []int64{math.MaxInt64, -math.MaxInt64, 3, 4, 5}
// 	arr_u1_num := []uint8{math.MaxInt8, 2, 3, 4, 5}
// 	arr_u2_num := []uint16{math.MaxInt16, 2, 3, 4, 5}
// 	arr_u4_num := []uint32{math.MaxInt32, 2, 3, 4, 5}
// 	arr_u8_num := []uint64{math.MaxInt64, 2, 3, 4, 5}
// 	arr_bt_bin := [][]byte{{'a'}, {'b'}, {'c'}, {'d'}, {'e'}}
// 	arr_f_num := []float32{0.0, 1.0, 0.00001, -0.00001, 94194939.12314345}
// 	arr_d_num := []float64{0.0, 1.0, 0.00001, -0.00001, 94194939.12314345}
// 	arr_str_js := []string{"a", "b", "c", "d", "e"}
// 	arr_u4_dt := []uint32{0, 1, 2, 3, 4}

// 	arr_sn_snum := make([]*snum.T_Snum, 0, 5)
// 	{
// 		sn_snum := &snum.T_Snum{}
// 		sn_snum.String__set("1")
// 		arr_sn_snum = append(arr_sn_snum, sn_snum)

// 		sn_snum = &snum.T_Snum{}
// 		sn_snum.String__set("-1")
// 		arr_sn_snum = append(arr_sn_snum, sn_snum)

// 		sn_snum = &snum.T_Snum{}
// 		sn_snum.String__set("9999999")
// 		arr_sn_snum = append(arr_sn_snum, sn_snum)

// 		sn_snum = &snum.T_Snum{}
// 		sn_snum.String__set("-9999999")
// 		arr_sn_snum = append(arr_sn_snum, sn_snum)

// 		sn_snum = &snum.T_Snum{}
// 		sn_snum.String__set("0")
// 		arr_sn_snum = append(arr_sn_snum, sn_snum)
// 	}
// 	arr_u8_dtn := make([]*db.T_Col__dtn, 0, 5)
// 	{
// 		n_dtn := &db.T_Col__dtn{}
// 		n_dtn.N8_dtn = 1534435
// 		arr_u8_dtn = append(arr_u8_dtn, n_dtn)
// 		n_dtn = &db.T_Col__dtn{}
// 		n_dtn.N8_dtn = 15
// 		arr_u8_dtn = append(arr_u8_dtn, n_dtn)
// 		n_dtn = &db.T_Col__dtn{}
// 		n_dtn.N8_dtn = 153
// 		arr_u8_dtn = append(arr_u8_dtn, n_dtn)
// 		n_dtn = &db.T_Col__dtn{}
// 		n_dtn.N8_dtn = 1534
// 		arr_u8_dtn = append(arr_u8_dtn, n_dtn)
// 		n_dtn = &db.T_Col__dtn{}
// 		n_dtn.N8_dtn = 15344
// 		arr_u8_dtn = append(arr_u8_dtn, n_dtn)
// 	}

// 	nn_last_insert_id, err := pc_dbms.C_insert.All_field__multi_insert(
// 		arr_u8_seq,
// 		arr_n1_bool,
// 		arr_s_str,
// 		arr_n1_num,
// 		arr_n2_num,
// 		arr_n4_num,
// 		arr_n8_num,
// 		arr_u1_num,
// 		arr_u2_num,
// 		arr_u4_num,
// 		arr_u8_num,
// 		arr_bt_bin,
// 		arr_sn_snum,
// 		arr_f_num,
// 		arr_d_num,
// 		arr_str_js,
// 		arr_u4_dt,
// 		arr_u8_dtn,
// 	)

// 	if err != nil {
// 		_t.Fatalf("err - %v\n", err)
// 	}
// 	fmt.Printf("LastInsertId - %d\n", nn_last_insert_id)
// 	pc_job.Commit()
// }

// // seq 는 auto increment 에 따라 작동
// func Test__insert__tpl__all_field__seq_null__multi_insert(_t *testing.T) {
// 	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, false)
// 	pc_dbms := &C_DB{}
// 	pc_dbms.Init(pc_job)

// 	arr_n1_bool := []int8{0, 1, 0, 1, 0}
// 	arr_s_str := []string{"a", "b", "c", "d", "e"}
// 	arr_n1_num := []int8{math.MaxInt8, -math.MaxInt8, 3, 4, 5}
// 	arr_n2_num := []int16{math.MaxInt16, -math.MaxInt16, 3, 4, 5}
// 	arr_n4_num := []int32{math.MaxInt32, -math.MaxInt32, 3, 4, 5}
// 	arr_n8_num := []int64{math.MaxInt64, -math.MaxInt64, 3, 4, 5}
// 	arr_u1_num := []uint8{math.MaxInt8, 2, 3, 4, 5}
// 	arr_u2_num := []uint16{math.MaxInt16, 2, 3, 4, 5}
// 	arr_u4_num := []uint32{math.MaxInt32, 2, 3, 4, 5}
// 	arr_u8_num := []uint64{math.MaxInt64, 2, 3, 4, 5}

// 	arr_bt_bin := [][]byte{{'a'}, {'b'}, {'c'}, {'d'}, {'e'}}
// 	arr_f_num := []float32{0.0, 1.0, 0.00001, -0.00001, 94194939.12314345}
// 	arr_d_num := []float64{0.0, 1.0, 0.00001, -0.00001, 94194939.12314345}
// 	arr_str_js := []string{"a", "b", "c", "d", "e"}
// 	arr_u4_dt := []uint32{0, 1, 2, 3, 4}

// 	arr_sn_snum := make([]*snum.T_Snum, 0, 5)
// 	{
// 		sn_snum := &snum.T_Snum{}
// 		sn_snum.String__set("1")
// 		arr_sn_snum = append(arr_sn_snum, sn_snum)

// 		sn_snum = &snum.T_Snum{}
// 		sn_snum.String__set("-1")
// 		arr_sn_snum = append(arr_sn_snum, sn_snum)

// 		sn_snum = &snum.T_Snum{}
// 		sn_snum.String__set("9999999")
// 		arr_sn_snum = append(arr_sn_snum, sn_snum)

// 		sn_snum = &snum.T_Snum{}
// 		sn_snum.String__set("-9999999")
// 		arr_sn_snum = append(arr_sn_snum, sn_snum)

// 		sn_snum = &snum.T_Snum{}
// 		sn_snum.String__set("0")
// 		arr_sn_snum = append(arr_sn_snum, sn_snum)
// 	}
// 	arr_u8_dtn := make([]*db.T_Col__dtn, 0, 5)
// 	{
// 		n_dtn := &db.T_Col__dtn{}
// 		n_dtn.N8_dtn = 1534435
// 		arr_u8_dtn = append(arr_u8_dtn, n_dtn)
// 		n_dtn = &db.T_Col__dtn{}
// 		n_dtn.N8_dtn = 15
// 		arr_u8_dtn = append(arr_u8_dtn, n_dtn)
// 		n_dtn = &db.T_Col__dtn{}
// 		n_dtn.N8_dtn = 153
// 		arr_u8_dtn = append(arr_u8_dtn, n_dtn)
// 		n_dtn = &db.T_Col__dtn{}
// 		n_dtn.N8_dtn = 1534
// 		arr_u8_dtn = append(arr_u8_dtn, n_dtn)
// 		n_dtn = &db.T_Col__dtn{}
// 		n_dtn.N8_dtn = 15344
// 		arr_u8_dtn = append(arr_u8_dtn, n_dtn)
// 	}

// 	nn_last_insert_id, err := pc_dbms.C_insert.Tpl__all_field__seq_null__multi_insert(
// 		DEF_s_test_table_name,
// 		arr_n1_bool,
// 		arr_s_str,
// 		arr_n1_num,
// 		arr_n2_num,
// 		arr_n4_num,
// 		arr_n8_num,
// 		arr_u1_num,
// 		arr_u2_num,
// 		arr_u4_num,
// 		arr_u8_num,
// 		arr_bt_bin,
// 		arr_sn_snum,
// 		arr_f_num,
// 		arr_d_num,
// 		arr_str_js,
// 		arr_u4_dt,
// 		arr_u8_dtn,
// 	)

// 	if err != nil {
// 		_t.Fatalf("err - %v\n", err)
// 	}
// 	fmt.Printf("LastInsertId - %d\n", nn_last_insert_id)
// 	pc_job.Commit()
// }
// */

func Test__insert__tpl__snum(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, false)
	pc_dbms := &C_DB{}
	pc_dbms.Init(pc_job)
	sn_snum := &db.T_Col__snum{}
	err := sn_snum.Encode([]byte("-999967676"))
	if err != nil {
		_t.Fatal(err)
	}

	nn_last_insert_id, err := pc_dbms.C_insert.Tpl__snum(DEF_s_test_table_name, sn_snum)
	if err != nil {
		_t.Fatal(err)
	}
	err = pc_job.Commit()
	if err != nil {
		_t.Fatal(err)
	}
	fmt.Println(nn_last_insert_id)
}

func Test__insert__tpl__snum__multi_insert(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, false)
	pc_dbms := &C_DB{}
	pc_dbms.Init(pc_job)

	arr_sn_snum := make([]*db.T_Col__snum, 0, 5)
	{
		sn_snum := &db.T_Col__snum{}
		sn_snum.Set__str("-9999999.99999999999999999999")
		arr_sn_snum = append(arr_sn_snum, sn_snum)

		sn_snum = &db.T_Col__snum{}
		sn_snum.Set__str("-9999999.999999999999999999")
		arr_sn_snum = append(arr_sn_snum, sn_snum)

		sn_snum = &db.T_Col__snum{}
		sn_snum.Set__str("-9999999.9999999999")
		arr_sn_snum = append(arr_sn_snum, sn_snum)

		sn_snum = &db.T_Col__snum{}
		sn_snum.Set__str("-9999999")
		arr_sn_snum = append(arr_sn_snum, sn_snum)

		sn_snum = &db.T_Col__snum{}
		sn_snum.Set__str("-9999998.999999")
		arr_sn_snum = append(arr_sn_snum, sn_snum)

		sn_snum = &db.T_Col__snum{}
		sn_snum.Set__str("-1.0000000000000001")
		arr_sn_snum = append(arr_sn_snum, sn_snum)

		sn_snum = &db.T_Col__snum{}
		sn_snum.Set__str("-1")
		arr_sn_snum = append(arr_sn_snum, sn_snum)

		sn_snum = &db.T_Col__snum{}
		sn_snum.Set__str("-0.999999999999")
		arr_sn_snum = append(arr_sn_snum, sn_snum)

		sn_snum = &db.T_Col__snum{}
		sn_snum.Set__str("-0.0000000000000001")
		arr_sn_snum = append(arr_sn_snum, sn_snum)

		sn_snum = &db.T_Col__snum{}
		sn_snum.Set__str("-0.00000000000000009")
		arr_sn_snum = append(arr_sn_snum, sn_snum)

		sn_snum = &db.T_Col__snum{}
		sn_snum.Set__str("-0")
		arr_sn_snum = append(arr_sn_snum, sn_snum)

		sn_snum = &db.T_Col__snum{}
		sn_snum.Set__str("0")
		arr_sn_snum = append(arr_sn_snum, sn_snum)

		sn_snum = &db.T_Col__snum{}
		sn_snum.Set__str("0.00000000000000009")
		arr_sn_snum = append(arr_sn_snum, sn_snum)

		sn_snum = &db.T_Col__snum{}
		sn_snum.Set__str("0.0000000000000001")
		arr_sn_snum = append(arr_sn_snum, sn_snum)

		sn_snum = &db.T_Col__snum{}
		sn_snum.Set__str("0.9999999999999999")
		arr_sn_snum = append(arr_sn_snum, sn_snum)

		sn_snum = &db.T_Col__snum{}
		sn_snum.Set__str("1")
		arr_sn_snum = append(arr_sn_snum, sn_snum)

		sn_snum = &db.T_Col__snum{}
		sn_snum.Set__str("1.00000000000000001")
		arr_sn_snum = append(arr_sn_snum, sn_snum)

		sn_snum = &db.T_Col__snum{}
		sn_snum.Set__str("9999998.999999")
		arr_sn_snum = append(arr_sn_snum, sn_snum)

		sn_snum = &db.T_Col__snum{}
		sn_snum.Set__str("9999999")
		arr_sn_snum = append(arr_sn_snum, sn_snum)

		sn_snum = &db.T_Col__snum{}
		sn_snum.Set__str("9999999.9999999999999")
		arr_sn_snum = append(arr_sn_snum, sn_snum)
	}

	nn_last_insert_id, err := pc_dbms.C_insert.Tpl__snum__multi_insert(DEF_s_test_table_name, arr_sn_snum)

	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
	fmt.Printf("LastInsertId - %d\n", nn_last_insert_id)
	pc_job.Commit()
}

func Test__update__tpl(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, false)
	pc_dbms := &C_DB{}
	pc_dbms.Init(pc_job)

	var n1_bool int8 = 1
	var s_string string = "newstring"
	nn_row_affected, err := pc_dbms.C_update.Tpl__arg(DEF_s_test_table_name, &n1_bool, &s_string, 4)
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
	fmt.Printf("RowAffected - %d\n", nn_row_affected)
	pc_job.Commit()
}

func Test__update__tpl__arg__add__left(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, false)
	pc_dbms := &C_DB{}
	pc_dbms.Init(pc_job)

	var u8_num uint64 = 1
	nn_row_affected, err := pc_dbms.C_update.Tpl__arg__add__left(DEF_s_test_table_name, &u8_num, 101)
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
	fmt.Printf("RowAffected - %d\n", nn_row_affected)
	pc_job.Commit()
}

func Test__update__tpl__arg__add__right(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local("dev_bp_test").Begin(sql.LevelDefault, false)
	pc_dbms := &C_DB{}
	pc_dbms.Init(pc_job)

	var u8_num uint64 = 2
	nn_row_affected, err := pc_dbms.C_update.Tpl__arg__add__right(DEF_s_test_table_name, &u8_num, 101)
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
	fmt.Printf("RowAffected - %d\n", nn_row_affected)
	pc_job.Commit()
}

func Test__update__tpl__arg__add__both(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local("dev_bp_test").Begin(sql.LevelDefault, false)
	pc_dbms := &C_DB{}
	pc_dbms.Init(pc_job)

	var u8_num uint64 = 2
	var u4_num uint32 = 5
	nn_row_affected, err := pc_dbms.C_update.Tpl__arg__add__both(DEF_s_test_table_name, &u8_num, &u4_num, 101)
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
	fmt.Printf("RowAffected - %d\n", nn_row_affected)
	pc_job.Commit()
}

func Test__update__arg__update_if_not_null(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, false)
	pc_dbms := &C_DB{}
	pc_dbms.Init(pc_job)
	var u1_num uint8 = 4
	_ = u1_num
	var is_bool int8 = 0
	_ = is_bool
	var s_str string = "hello"
	_ = s_str

	nn_row_affected, err := pc_dbms.C_update.Tpl__arg__without_null(DEF_s_test_table_name, nil, nil, &s_str, 1)
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
	fmt.Printf("RowAffected - %d\n", nn_row_affected)
	pc_job.Commit()
}

func Test__delete__all(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, false)
	pc_dbms := &C_DB{}
	pc_dbms.Init(pc_job)
	nn_row_affected, err := pc_dbms.C_delete.Tpl__all(DEF_s_test_table_name)
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
	fmt.Printf("RowAffected - %d\n", nn_row_affected)
	pc_job.Commit()
}
