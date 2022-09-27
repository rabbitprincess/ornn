package bp

import (
	"fmt"
	"strings"
)

func Clear__in_quot(_s_sql string) (s_ret string) {
	const (
		DEF_s_quot_single string = "'"
		DEF_s_quot_double string = "\""
		DEF_s_apostrophe  string = "`"
	)

	s_ret = ""
	is_in_quot_single := false
	is_in_quot_double := false
	is_in_quot_apostrophe := false
	for i := 0; i < len(_s_sql); i++ {
		s_one := _s_sql[i : i+1]

		is_find_any := false
		if s_one == DEF_s_quot_single {
			is_find_any = true
			is_in_quot_single = !is_in_quot_single
		}
		if s_one == DEF_s_quot_double {
			is_find_any = true
			is_in_quot_double = !is_in_quot_double
		}
		if s_one == DEF_s_apostrophe {
			is_find_any = true
			is_in_quot_apostrophe = !is_in_quot_apostrophe
		}

		if is_in_quot_single == true ||
			is_in_quot_double == true ||
			is_in_quot_apostrophe == true ||
			is_find_any == true {
			s_ret += " "
			continue
		}

		s_ret += s_one
	}
	return s_ret
}

func Util__split_by_delimiter(_s_sql, _s_delimiter string) (s_ret_before, s_ret_after string) {
	s_sql__lowercase := strings.ToLower(_s_sql)
	s_sql__lowercase = Clear__in_quot(s_sql__lowercase)

	n_pos := strings.Index(s_sql__lowercase, _s_delimiter)
	if n_pos == -1 {
		return _s_sql, ""
	}
	s_ret_before = _s_sql[:n_pos]
	s_ret_after = _s_sql[n_pos:]
	return s_ret_before, s_ret_after
}

func Util__export_str_between_delimiter(
	_s_input string,
	_s_delimiter string,
) (
	arrs_arg_name []string,
	err error,
) {
	s_input__without_in_quot := Clear__in_quot(_s_input)
	arrs_arg_name = make([]string, 0, 10)
	is_mode_in := false
	s_buf := ""
	for i := 0; i < len(_s_input); i++ {
		s_at := _s_input[i : i+1]
		s_at_quat := s_input__without_in_quot[i : i+1]

		if s_at_quat == _s_delimiter {
			if is_mode_in == true {
				for _, s_arg_name__after_where := range arrs_arg_name {
					if s_arg_name__after_where == s_buf {
						return nil, fmt.Errorf("duplicate field name | field name : %s", s_arg_name__after_where)
					}
				}
				arrs_arg_name = append(arrs_arg_name, s_buf)
				s_buf = ""
			}
			is_mode_in = !is_mode_in
			continue
		}
		if is_mode_in == true {
			s_buf += s_at
		}
	}
	return arrs_arg_name, nil
}

func Util__replace_str__between_delimiter(
	_s_input string,
	_s_delimiter string,
	_s_delimiter_after string,
) (
	s_output string,
) {
	s_input__without_in_quot := Clear__in_quot(_s_input)
	is_mode_in := false
	for i := 0; i < len(_s_input); i++ {
		s_at := _s_input[i : i+1]
		s_at_quat := s_input__without_in_quot[i : i+1]

		if s_at_quat == _s_delimiter {
			if is_mode_in == true {
				s_output += _s_delimiter_after
			}
			is_mode_in = !is_mode_in
			continue
		}
		if is_mode_in == false {
			s_output += s_at
		}
	}
	return s_output
}

func Util__Clear_delimiter(
	_s_input string,
	_s_delimiter string,
) (
	s_output string,
) {
	s_input__without_in_quot := Clear__in_quot(_s_input)
	for i := 0; i < len(_s_input); i++ {
		s_at := _s_input[i : i+1]
		if s_input__without_in_quot[i:i+1] == _s_delimiter {
			continue
		}
		s_output += s_at
	}
	return s_output
}

func Util__Replace_str__in_delimiter_value(
	_s_input string,
	_s_delimiter string,
	_s_spliter string,
) (
	s_output string,
) {
	//  xxxx#AAAA#xxxx 		-> xxxxAAAAxxxx
	//  xxxx#AAAA/BBBB#xxxx -> xxxxBBBBxxxx
	s_input__without_in_quot := Clear__in_quot(_s_input)

	var s_buf string
	var is_mode_in bool

	for i := 0; i < len(_s_input); i++ {
		s_at := _s_input[i : i+1]
		s_at_quat := s_input__without_in_quot[i : i+1]

		if s_at_quat == _s_delimiter {
			is_mode_in = !is_mode_in

			if is_mode_in == true { // 구분자 진입 시점
				s_buf = ""
			} else { // 구분자 탈출 시점
				s_output += s_buf
			}
			continue
		}

		if is_mode_in == false {
			s_output += s_at
		} else {
			if s_at_quat == _s_spliter {
				s_buf = ""
				continue
			}
			s_buf += s_at
		}
	}
	return s_output
}

// insert into `aaa` values (a, b, c, d) -> a, b, c, d 추출
func Util__export__insert_query_values(_s_sql_insert string) string {
	s_sql_insert__lowercase := strings.ToLower(_s_sql_insert)
	n_idx__after_values := strings.Index(s_sql_insert__lowercase, "values")
	if n_idx__after_values == -1 {
		return ""
	}
	var n_pos_start, n_pos_end int
	for i := n_idx__after_values; i < len(_s_sql_insert); i++ {
		if _s_sql_insert[i:i+1] == "(" {
			n_pos_start = i
		}
		if _s_sql_insert[i:i+1] == ")" {
			n_pos_end = i
		}
	}
	return _s_sql_insert[n_pos_start+1 : n_pos_end]
}

func Util__conv_first_upper_case(_s string) string {
	if len(_s) == 0 {
		return ""
	}

	var s_ret string
	s_ret = strings.ToUpper(_s[0:1]) + _s[1:]
	return s_ret
}

func Util__is_parser_val__arg(_bt_val []byte) bool {
	s_val := string(_bt_val)
	if len(s_val) < 3 { // :v0 | :v1 | :vN  임으로 3보다 작으면 추출 할 대상이 아님
		return false
	}
	if s_val[0:2] != ":v" { // ? 종류인지 구분 = 아니면 무시
		return false
	}
	return true
}
