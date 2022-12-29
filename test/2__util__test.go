package go_orm_gen_test

import (
	"fmt"
	"module/solution/bp"
	"reflect"
	"strings"
	"testing"
)

//----------------------------------------------------------------------------------------//
// CRUD

/*
json 구조
	환경설정
		BASE
			BP-Version
			dbms 종류 + DSN 정보(id/pw)
	작업물
		DB
			쿼리정의
				DB
					함수 고유번호
							- group:
							- func_name:
							- CRUD 타입: select or insert or update or delete
								1.타입별-select
									select
										* or field_name...n, 원하는데로 (다중 테이블이면 aa.* , bb.* )
									from
										적고싶은데로 - 테이블 이름 인식하여 select 및 BP_FIELD_TYPE 인식
									where
										적고싶은데로 - ? 인식하여 반대편 형 감지하여 input arg 서술 ( 못찾으면 컴파일러 터지게 적음 ex> XXX )
									리턴 구조체 이름 (구현 보류-당장필요없음-나중엔 꼭 필요 양방향 수정을 위함)
										[ // return_struct_name: T_Ret_query_name ] 을 -주석으로-적어주고 디폴트 생성이름을 쓰거나 주석 풀고 직접 정의
											같은 select 를 반환하는 경우(타입이)
												-> 1. 같은 리턴 명을 사용 할 수 있고 .go source 상에 1번만 정의 된다.
												-> 2. 다른 리턴 명을 사용 할 수 있고 여러번 정의 될 수도 있다. (편의성 측면에서 필요하면 사용 가능)
										## 함수정의 내용 자체를 md5로 만들어 .go source 상에 함수명 같은줄에 // 주석으로 기록 해 놓고

								2.타입별-insert
									예시 - insert into table_name(field_name...n) VALUES(values...n)
									into table
										테이블이름 1개 (단일)
									target field name
										해당 테이블에 속해있는 필드 , 단 * 이면 전부 나열 ( 나중에 함수 입력 인자가 됨 - BP_FIELD_TYPE 추적  )
									XX - target field values
										json 에 받을게 없음 ( 안해도됨 )

								3.타입별-update
									insert 와 동일 + where 추가됨
									where 절은 select 참조

								4.타입별-delete
									where 절은 select 참조


정의
	DB
		설정 - 자동 구성
			BP_FIELD_TYPE 초기값 -> 실제 사용값으로 변경 필요
			DB SCHEMA md5 -> 변경 시 감지하기 위함
		설정 - 검증 ( 서비스 실행 시 )
			DB SCHEMA md5
		작업 - 쿼리 정의
			패키지명
			함수명
				Select 문
					select = * or 원하는 기본 필드(필드타입 자동 정의) + 커스텀 필드 ( 필드 타입 선택 필요 )
						Next_row 실행 판정 함수 = 필요시 인자에 포함

				Insert 문
					from = 대상 table
					set = * or 원하는 기본필드(필드타입 자동 정의)
						처리 콜백 이름 = 패키지명 + 함수명 + 입력 대상 필드 이름  | 리턴형은 입력 대상 필드의 필드 타입과 같도록 자동 생성되어야 함.
				Update 문
					Insert 와 동일

				Delete


	rest
		사용자 정의 타입을 먼저 정의
			cgi 로 들어오는 입력에 대한 사용자 정의 타입
			모든 입력은 일단 문자([]byte 혹은 string 으로 가정)
				사용자 정의 형에 대한 변환기를 함수 정의 하고 -> 소스에서 직접 구현 한뒤 -> 해당 rest 호출시 전처리를 자동으로 한다.
					ex> func Argment_conv__BT__TO__SNUM(_bt_input []byte) (SNUM,error) { .... return &SNUM{}, nil }
						-> TD_ASSET 같은 경우 ASSET 에 해당하지 않으면 에러 리턴 -> 해당 rest 에러 출력 후 수행 종료 = 전처리가 획일화 됨
						-> DB 관련 입력값은 SQl 인젝션 공격에 대해 미리 거르는 코드를 넣을 수 있음.
						-> base64 형태로 오는 데이터는 미리 디코딩 후 사용자에게 제공가능


		rest 함수를 정의
				rest 카테고리
				rest 이름
				argment
					인자1=이름,타입(사용자정의 타입)
					인자n=이름,타입(사용자정의 타입)
				schema - http or https or both
				method - get post
				return struct 정의
					이름,형(단 사용자 정의 타입이여야 함 - 프론트 앤드 개발자와 공유 하는 용도)





*/

func Test__clear_str__in_quota(_t *testing.T) {
	if bp.Clear__in_quot("'will delete'") !=
		fmt.Sprintf("%s", strings.Repeat(" ", len("'will delete'"))) {
		_t.Error("not same")
	}

	if bp.Clear__in_quot(`"will delete"`) !=
		fmt.Sprintf("%s", strings.Repeat(" ", len(`"will delete"`))) {
		_t.Error("not same")
	}

	fmt.Printf("==%s==\n", bp.Clear__in_quot("`will delete`"))
	if bp.Clear__in_quot("`will delete`") !=
		fmt.Sprintf("%s", strings.Repeat(" ", len("`will delete`"))) {
		_t.Error("not same")
	}
}

func Test__split__by_delimiter(_t *testing.T) {
	type T_Testcase struct {
		s_testcase   string
		s_ret_before string
		s_ret_after  string
	}

	arrpt_testcase := make([]*T_Testcase, 0, 10)
	{
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase:   "'?' \"?\" `?` ?",
			s_ret_before: "'?' \"?\" `?` ?",
			s_ret_after:  "",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase:   "??????????",
			s_ret_before: "??????????",
			s_ret_after:  "",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase:   "?'?'?'?'?'?'?'?'?'?'",
			s_ret_before: "?'?'?'?'?'?'?'?'?'?'",
			s_ret_after:  "",
		})

		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase:   "SELECT * FROM table_name WHERE field_name_1 = ? and field_name_2 = ? and field_name_3 = ? and field_name_4 = ?;",
			s_ret_before: "SELECT * FROM table_name ",
			s_ret_after:  "WHERE field_name_1 = ? and field_name_2 = ? and field_name_3 = ? and field_name_4 = ?;",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase:   "INSERT INTO table_name (field_name_1,field_name_2,field_name_3,field_name_4) VALUES (?,?,?,?);",
			s_ret_before: "INSERT INTO table_name (field_name_1,field_name_2,field_name_3,field_name_4) VALUES (?,?,?,?);",
			s_ret_after:  "",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase:   "UPDATE table_name SET field_name_1=?, field_name_2=?, field_name_3=?, field_name_4=? WHERE field_name_1 = ? and field_name_2 = ? and field_name_3 = ? and field_name_4 = ?;",
			s_ret_before: "UPDATE table_name SET field_name_1=?, field_name_2=?, field_name_3=?, field_name_4=? ",
			s_ret_after:  "WHERE field_name_1 = ? and field_name_2 = ? and field_name_3 = ? and field_name_4 = ?;",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase:   "DELETE FROM table_name WHERE field_name_1 = ? and field_name_2 = ? and field_name_3 = ? and field_name_4 = ?;",
			s_ret_before: "DELETE FROM table_name ",
			s_ret_after:  "WHERE field_name_1 = ? and field_name_2 = ? and field_name_3 = ? and field_name_4 = ?;",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase:   "SELECT * FROM table_name WHERE field_name_1 = ? and field_name_2 = '?' and field_name_3 = \"?\" and field_name_4 = `?`;",
			s_ret_before: "SELECT * FROM table_name ",
			s_ret_after:  "WHERE field_name_1 = ? and field_name_2 = '?' and field_name_3 = \"?\" and field_name_4 = `?`;",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase:   "INSERT INTO table_name (field_name_1,field_name_2,field_name_3,field_name_4) VALUES (?,'?', \"?\",`?`);",
			s_ret_before: "INSERT INTO table_name (field_name_1,field_name_2,field_name_3,field_name_4) VALUES (?,'?', \"?\",`?`);",
			s_ret_after:  "",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase:   "UPDATE table_name SET field_name_1=?, field_name_2='?', field_name_3=\"?\", field_name_4=`?` WHERE field_name_1 = ? and field_name_2 = '?' and field_name_3 = \"?\" and field_name_4 = `?`;",
			s_ret_before: "UPDATE table_name SET field_name_1=?, field_name_2='?', field_name_3=\"?\", field_name_4=`?` ",
			s_ret_after:  "WHERE field_name_1 = ? and field_name_2 = '?' and field_name_3 = \"?\" and field_name_4 = `?`;",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase:   "DELETE FROM table_name WHERE field_name_1 = ? and field_name_2 = '?' and field_name_3 = \"?\" and field_name_4 = `?`;",
			s_ret_before: "DELETE FROM table_name ",
			s_ret_after:  "WHERE field_name_1 = ? and field_name_2 = '?' and field_name_3 = \"?\" and field_name_4 = `?`;",
		})
	}

	for _, pt_testcase := range arrpt_testcase {
		s_ret_before, s_ret_after := bp.Util__split_by_delimiter(pt_testcase.s_testcase, "where")
		if s_ret_before != pt_testcase.s_ret_before {
			_t.Errorf("Clear_after_where is not same\n\tinput : %s\n\tresult : %s", s_ret_before, pt_testcase.s_ret_before)
		}

		if s_ret_after != pt_testcase.s_ret_after {
			_t.Errorf("Clear_after_where is not same\n\tinput : %s\n\tresult : %s", s_ret_before, pt_testcase.s_ret_before)
		}

	}
}

func Test__export_between_delimiter(_t *testing.T) {
	type T_Testcase struct {
		s_testcase string
		arrs_ret   []string
	}

	arrpt_testcase := make([]*T_Testcase, 0, 10)
	{
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "%a%",
			arrs_ret:   []string{"a"},
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "%a%b",
			arrs_ret:   []string{"a"},
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "a%b%",
			arrs_ret:   []string{"b"},
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "a%b%c",
			arrs_ret:   []string{"b"},
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "%a%b%c%",
			arrs_ret:   []string{"a", "c"},
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "a%b%c%d%e",
			arrs_ret:   []string{"b", "d"},
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "%a%b%c%d%e%",
			arrs_ret:   []string{"a", "c", "e"},
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "%a%`%b%`%c%d%e%",
			arrs_ret:   []string{"a", "c", "e"},
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "%a`%b%`c%d%e%",
			arrs_ret:   []string{"a`%b%`c", "e"},
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "%a%'%b%c%d%'%e%",
			arrs_ret:   []string{"a", "e"},
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "%a'%b%c%d%'e%",
			arrs_ret:   []string{"a'%b%c%d%'e"},
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "'%a'%b%c%d%'e%'",
			arrs_ret:   []string{"b", "d"},
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "'%a%b%c%d%e%'",
			arrs_ret:   []string{},
		})
	}

	for _, pt_testcase := range arrpt_testcase {
		arrs_ret, err := bp.Util__export_str_between_delimiter(pt_testcase.s_testcase, "%")
		if err != nil {
			_t.Errorf("err - %v", err)
		} else if reflect.DeepEqual(arrs_ret, pt_testcase.arrs_ret) == false {
			_t.Errorf("Clear_after_where is not same\n\tinput : %s\n\tresult : %s", arrs_ret, pt_testcase.arrs_ret)
		}
	}
}

func Test__change_delimiter_to_input_mark(_t *testing.T) {
	type T_Testcase struct {
		s_testcase string
		s_ret      string
	}

	arrpt_testcase := make([]*T_Testcase, 0, 10)
	{
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "%a%",
			s_ret:      " ",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "%a%b",
			s_ret:      " b",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "a%b%",
			s_ret:      "a ",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "a%b%c",
			s_ret:      "a c",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "%a%b%c%",
			s_ret:      " b ",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "a%b%c%d%e",
			s_ret:      "a c e",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "%a%b%c%d%e%",
			s_ret:      " b d ",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "%a%`%b%`%c%d%e%",
			s_ret:      " `%b%` d ",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "%a`%b%`c%d%e%",
			s_ret:      " d ",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "%a%'%b%c%d%'%e%",
			s_ret:      " '%b%c%d%' ",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "%a'%b%c%d%'e%",
			s_ret:      " ",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "'%a'%b%c%d%'e%'",
			s_ret:      "'%a' c 'e%'",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "'%a%b%c%d%e%'",
			s_ret:      "'%a%b%c%d%e%'",
		})
	}

	for _, pt_testcase := range arrpt_testcase {
		s_ret := bp.Util__replace_str__between_delimiter(pt_testcase.s_testcase, "%", " ")
		if s_ret != pt_testcase.s_ret {
			_t.Errorf("change_delimiter is not same\n\tinput : %s\n\treturn : %s\n\tresult : %s", pt_testcase.s_testcase, s_ret, pt_testcase.s_ret)
		}
	}
}

func Test__clear_delimiter(_t *testing.T) {
	type T_Testcase struct {
		s_testcase string
		s_ret      string
	}

	arrpt_testcase := make([]*T_Testcase, 0, 10)
	{
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "%a%",
			s_ret:      "a",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "%a%b%",
			s_ret:      "ab",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "%%%",
			s_ret:      "",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_testcase: "a%a'%a'a%a",
			s_ret:      "aa'%a'aa",
		})
	}

	for _, pt_testcase := range arrpt_testcase {
		s_ret := bp.Util__Clear_delimiter(pt_testcase.s_testcase, "%")
		if s_ret != pt_testcase.s_ret {
			_t.Errorf("change_delimiter is not same\n\tinput : %s\n\treturn : %s\n\tresult : %s", pt_testcase.s_testcase, s_ret, pt_testcase.s_ret)
		}
	}
}

func Test__replace_str__in_delimiter_value(_t *testing.T) {
	type T_Testcase struct {
		s_input     string
		s_ret       string
		s_delimiter string
		s_spliter   string
	}

	arrpt_testcase := make([]*T_Testcase, 0, 10)
	{
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_input:     "%a%",
			s_ret:       "a",
			s_delimiter: "%",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_input:     "%a%",
			s_ret:       "a",
			s_delimiter: "%",
			s_spliter:   "/",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_input:     "%a%b%",
			s_ret:       "ab",
			s_delimiter: "%",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_input:     "%%%",
			s_ret:       "",
			s_delimiter: "%",
			s_spliter:   "/",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_input:     "a%a'%a'a%a",
			s_ret:       "aa'%a'aa",
			s_delimiter: "%",
			s_spliter:   "/",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_input:     "%a/b%",
			s_ret:       "a/b",
			s_delimiter: "%",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_input:     "%a/b%",
			s_ret:       "b",
			s_delimiter: "%",
			s_spliter:   "/",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_input:     "%ab%",
			s_ret:       "ab",
			s_delimiter: "%",
			s_spliter:   "/",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_input:     "%a/b%c/d%e/f%",
			s_ret:       "bc/df",
			s_delimiter: "%",
			s_spliter:   "/",
		})
	}

	for _, pt_testcase := range arrpt_testcase {
		s_ret := bp.Util__Replace_str__in_delimiter_value(pt_testcase.s_input, pt_testcase.s_delimiter, pt_testcase.s_spliter)
		if s_ret != pt_testcase.s_ret {
			_t.Errorf("change_delimiter is not same\n\tinput : %s\n\treturn : %s\n\tresult : %s", pt_testcase.s_input, s_ret, pt_testcase.s_ret)
		}
	}
}

func Test__export__insert_query_values(_t *testing.T) {
	type T_Testcase struct {
		s_input string
		s_ret   string
	}
	arrpt_testcase := make([]*T_Testcase, 0, 10)
	{
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_input: "INSERT INTO `test table` VALUES (?, ?, ?, ?)",
			s_ret:   "?, ?, ?, ?",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_input: "INSERT INTO `test table` VALUES (?, ?, ?, ?);",
			s_ret:   "?, ?, ?, ?",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_input: "INSERT INTO `test table` (`seq`, `str`, `num`, `dtn`) VALUES (?, ?, ?, ?);",
			s_ret:   "?, ?, ?, ?",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_input: "INSERT INTO `test table` (`seq`, `str`, `num`, `dtn`) VALUES (`1`, `안녕하세요`, `12345`, ?);",
			s_ret:   "`1`, `안녕하세요`, `12345`, ?",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_input: "INSERT INTO Customers (CustomerName, ContactName, Address, City, PostalCode, Country) VALUES (`SupplierName`, `ContactName`, ?, ?, ?, `Country`)",
			s_ret:   "`SupplierName`, `ContactName`, ?, ?, ?, `Country`",
		})

		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_input: "INSERT INTO Customers (CustomerName, ContactName, Address, City, PostalCode, Country) SELECT SupplierName, ContactName, Address, City, PostalCode, Country FROM Suppliers;",
			s_ret:   "",
		})
		arrpt_testcase = append(arrpt_testcase, &T_Testcase{
			s_input: "INSERT INTO test_table values (?, ?, ?, ?);",
			s_ret:   "?, ?, ?, ?",
		})
	}

	for _, pt_testcase := range arrpt_testcase {
		s_ret := bp.Util__export__insert_query_values(pt_testcase.s_input)
		if s_ret != pt_testcase.s_ret {
			_t.Errorf("export__insert_query_values is not same\n\tinput : %s\n\treturn : %s\n\tresult : %s", pt_testcase.s_input, s_ret, pt_testcase.s_ret)
		}
	}
}
