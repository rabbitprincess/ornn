package go_orm_gen_db

import (
	"database/sql"
	"fmt"
	"math"
	"module/db"
	"module/db/db_mysql/test_env"
	"testing"
)

const (
	DEF_s_test_db_name    = "dev_bp_sample"
	DEF_s_test_table_name = "test_table"
)

//-------------------------------------------------------------------------------------//
// common

func Test__common__arg(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, true)
	defer pc_job.Commit()
	pc_db := &C_DB{}
	pc_db.Init(pc_job)

	arrpt, err := pc_db.C_common.Arg(1)
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
	for _, pt := range arrpt {
		fmt.Printf("%+v\n", pt)
	}
}

func Test__common__arg__limit_one(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, true)
	defer pc_job.Commit()
	pc_db := &C_DB{}
	pc_db.Init(pc_job)

	pt, err := pc_db.C_common.Arg__limit_one(1)
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
	fmt.Printf("%+v\n", pt)
}

func Test__common__tpl(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, false)
	defer pc_job.Commit()
	pc_db := &C_DB{}
	pc_db.Init(pc_job)

	s_str := "hello world"
	n1_num := int8(1)
	pt, err := pc_db.C_common.Tpl(DEF_s_test_table_name, &s_str, &n1_num)
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
	fmt.Printf("%+v\n", pt)
}

//-------------------------------------------------------------------------------------//
// select

func Test__select__limit__offset__group_by(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, true)
	defer pc_job.Commit()
	pc_db := &C_DB{}
	pc_db.Init(pc_job)

	arrpt, err := pc_db.C_select.Group_by__limit__offset(0, "u1_num", 10, 10)
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
	for _, pt := range arrpt {
		fmt.Printf("%+v\n", pt)
	}
}

func Test__select__tpl(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, true)
	defer pc_job.Commit()
	pc_db := &C_DB{}
	pc_db.Init(pc_job)

	arrpt, err := pc_db.C_select.Tpl(DEF_s_test_table_name)
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
	for _, pt := range arrpt {
		fmt.Printf("%+v\n", pt)
	}
}

func Test__select__custom_type(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, true)
	defer pc_job.Commit()
	pc_db := &C_DB{}
	pc_db.Init(pc_job)

	arrpt, err := pc_db.C_select.Custom_field()
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
	for i, pt := range arrpt {
		// bool, snum
		fmt.Printf("idx : %v\n\tbool : %v\n\tsnum : %v\n", i, pt.Is_bool, pt.Bt_snum.String())
	}
	pc_job.Commit()
}

//-------------------------------------------------------------------------------------//
// insert

func Test__insert__auto_increment(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, false)
	defer pc_job.Commit()
	pc_db := &C_DB{}
	pc_db.Init(pc_job)

	var is_bool bool = true
	var s_str string = "hello"
	var bt_bin []byte = []byte("hello")
	var n1_num int8 = math.MaxInt8
	var n2_num int16 = math.MaxInt16
	var n4_num int32 = math.MaxInt32
	var n8_num int64 = math.MaxInt64
	var u1_num uint8 = math.MaxInt8
	var u2_num uint16 = math.MaxInt16
	var u4_num uint32 = math.MaxInt32
	var u8_num uint64 = math.MaxInt64
	var f4_num float32 = math.MaxFloat32
	var f8_num float64 = math.MaxFloat64
	pt_snum := &db.T_Col__snum{}
	pt_snum.Set__str("123")
	_, err := pc_db.C_insert.Auto_increment(
		&is_bool,
		&s_str,
		&bt_bin,
		&n1_num,
		&n2_num,
		&n4_num,
		&n8_num,
		&u1_num,
		&u2_num,
		&u4_num,
		&u8_num,
		&f4_num,
		&f8_num,
		pt_snum,
	)
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
}

func Test__insert__multi(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, false)
	defer pc_job.Commit()
	pc_db := &C_DB{}
	pc_db.Init(pc_job)

	arrpis_bool := make([]*bool, 0, 100)
	arrps_str := make([]*string, 0, 100)
	arrpn1_num := make([]*int8, 0, 100)
	for i := 0; i < 100; i++ {
		is_bool := ([]bool{true, false})[i%2]
		s_str := fmt.Sprintf("hello %v", i)
		n1_num := int8(i)

		arrpis_bool = append(arrpis_bool, &is_bool)
		arrps_str = append(arrps_str, &s_str)
		arrpn1_num = append(arrpn1_num, &n1_num)
	}

	pc_db.C_insert.Multi_insert(arrpis_bool, arrps_str, arrpn1_num)

}

//-------------------------------------------------------------------------------------//
// update

func Test__update__add(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, false)
	defer pc_job.Commit()
	pc_db := &C_DB{}
	pc_db.Init(pc_job)

	var u8_num__add uint64 = 100
	_, err := pc_db.C_update.Add(&u8_num__add, 1)
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
}

func Test__update__null_ignore(_t *testing.T) {
	pc_job, _ := test_env.DB__get_test_mysql__local(DEF_s_test_db_name).Begin(sql.LevelDefault, false)
	defer pc_job.Commit()
	pc_db := &C_DB{}
	pc_db.Init(pc_job)

	var s_str = "hello_after"
	var u1_num uint8 = 123
	_, err := pc_db.C_update.Null_ignore(nil, &s_str, &u1_num, 1)
	if err != nil {
		_t.Fatalf("err - %v\n", err)
	}
}
