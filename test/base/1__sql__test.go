package bp_test

import (
	"fmt"
	"module/db/db_mysql/test_env"
	"module/solution/bp"
	"testing"

	"github.com/blastrain/vitess-sqlparser/sqlparser"
)

//------------------------------------------------------------------------------------------//
// insert

func Test__parser__sql_create_table(_t *testing.T) {
	s_table_name := "bp_gen_test_schema"
	s_sql := "show create table `" + s_table_name + "`"

	var s_query string
	pc_db := test_env.DB__get_test_mysql__local("dev_bp_test")
	_, err := pc_db.Query(s_sql).Row_next(&s_table_name, &s_query)
	if err != nil {
		_t.Fatal(err)
	}

	stmt, err := sqlparser.Parse(s_query)
	if err != nil {
		_t.Fatal(err)
	}
	pt_parser, is_ok := stmt.(*sqlparser.CreateTable)
	if is_ok == false {
		_t.Fatal("fail")
	}

	// print
	{
		// table name
		s_table_name := pt_parser.NewName.Name.String()
		fmt.Printf("table name - %v\n", s_table_name)

		// field
		for _, column := range pt_parser.Columns {
			fmt.Printf("\tfield name, type : %v | %v\n", column.Name, column.Type)
		}

		// index
		for _, constr := range pt_parser.Constraints {
			fmt.Printf("\tindex name, type : %v | %v\n", constr.Name, constr.Type)
			for _, key := range constr.Keys {
				fmt.Printf("\tkeys : %v\n", key.String())
			}
		}
	}
}

func Test__parser__sql_insert(_t *testing.T) {
	type T_Testcase struct {
		s_sql     string
		s_explain string
	}

	arrpt_testcase := make([]*T_Testcase, 0, 10)
	{

		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_sql:     "insert into test_table (u8_seq,is_bool,s_str,n1_num,n2_num,n4_num,n8_num,bt_bin,sn_snum,f_num,d_num,js_str,dt_num,dtn_num) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			s_explain: "all field name",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_sql:     "insert into test_table (u8_seq,is_bool) values (?, ?)",
			s_explain: "some field name",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_sql:     "insert into test_table values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
			s_explain: "no field name",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_sql:     "insert into test_table set u8_seq = ?, is_bool = ?",
			s_explain: "insert with set",
		})

		// on duplicate update
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_sql:     "insert into test_table set u8_seq = ?, is_bool = ? on duplicate key update u8_seq = ?, is_bool = ?",
			s_explain: "insert with set",
		})

		// 임시 - 파서 미지원 ( 추후 bp 고도화 예정 )
		/*
			arrpt_testcase = append(arrpt_testcase, &T_Testcase{
				s_sql:     "insert into test_table values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
				s_explain: "multi insert",
			})
			// multi insert
			arrpt_testcase = append(arrpt_testcase, &T_Testcase{
				s_sql: "insert into test_table values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?), (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
				s_explain: "multi insert",
			})
			arrpt_testcase = append(arrpt_testcase, &T_Testcase{
				s_sql: "insert into test_table SELECT * from test_table",
				s_explain: "insert select",
			})
			arrpt_testcase = append(arrpt_testcase, &T_Testcase{
				s_sql: "insert into test_table (u8_seq, is_bool, s_str) select u8_seq, is_bool, s_str from test_table where u8_seq = ?",
				s_explain: "insert select with where",
			})
		*/
	}

	for _, pt_testcase := range arrpt_testcase {
		pt_sql := bp.T_SQL{}
		i_sql, err := pt_sql.Get_parser(pt_testcase.s_sql)
		if err != nil {
			_t.Fatal(err)
		}
		fmt.Printf("\n===========================================================\n\n")
		fmt.Printf("sql\n\t%s\n", pt_testcase.s_sql)
		fmt.Printf("explain\n\t%s\n", pt_testcase.s_explain)
		fmt.Printf("table name\n\t%v\n", i_sql.(*bp.T_SQL__insert).S_table_name)
		// fmt.Printf("fields name\n\t%v\n\n", i_sql.(*bp.T_SQL__insert).Arrpt_field)
	}
}
