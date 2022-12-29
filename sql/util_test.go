package sql

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test__clear_str__in_quota(_t *testing.T) {
	if Clear__in_quot("'will delete'") !=
		fmt.Sprintf("%s", strings.Repeat(" ", len("'will delete'"))) {
		_t.Error("not same")
	}

	if Clear__in_quot(`"will delete"`) !=
		fmt.Sprintf("%s", strings.Repeat(" ", len(`"will delete"`))) {
		_t.Error("not same")
	}

	fmt.Printf("==%s==\n", Clear__in_quot("`will delete`"))
	if Clear__in_quot("`will delete`") !=
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
		s_ret_before, s_ret_after := Util__split_by_delimiter(pt_testcase.s_testcase, "where")
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
		arrs_ret, err := Util__export_str_between_delimiter(pt_testcase.s_testcase, "%")
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
		s_ret := Util__replace_str__between_delimiter(pt_testcase.s_testcase, "%", " ")
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
		s_ret := Util__Clear_delimiter(pt_testcase.s_testcase, "%")
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
		s_ret := Util__Replace_str__in_delimiter_value(pt_testcase.s_input, pt_testcase.s_delimiter, pt_testcase.s_spliter)
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
		s_ret := Util__export__insert_query_values(pt_testcase.s_input)
		if s_ret != pt_testcase.s_ret {
			_t.Errorf("export__insert_query_values is not same\n\tinput : %s\n\treturn : %s\n\tresult : %s", pt_testcase.s_input, s_ret, pt_testcase.s_ret)
		}
	}
}
