package sql

import (
	"fmt"
	"strings"
)

const (
	DEF_s_sql__prepare_statement__delimiter = "%"
	DEF_s_sql__prepare_statement__after     = "?"

	DEF_s_sql__tpl__delimiter = "#"
	DEF_s_sql__tpl__after     = "%s"
	DEF_s_sql__tpl__split     = "/"
)

func Util_ClearInQuot(input string) (output string) {
	var inQuotSingle, inQuotDouble, inQuotApostrophe bool

	for i := 0; i < len(input); i++ {
		one := input[i : i+1]

		find := false
		if one == "'" {
			find = true
			inQuotSingle = !inQuotSingle
		}
		if one == "\"" {
			find = true
			inQuotDouble = !inQuotDouble
		}
		if one == "`" {
			find = true
			inQuotApostrophe = !inQuotApostrophe
		}

		if inQuotSingle == true ||
			inQuotDouble == true ||
			inQuotApostrophe == true ||
			find == true {
			output += " "
			continue
		}

		output += one
	}
	return output
}

func Util_SplitByDelimiter(sql, delimiter string) (before, after string) {
	sqlLower := strings.ToLower(sql)
	sqlLower = Util_ClearInQuot(sqlLower)

	pos := strings.Index(sqlLower, delimiter)
	if pos == -1 {
		return sql, ""
	}
	return sql[:pos], sql[pos:]
}

func Util_ExportBetweenDelimiter(input string, delimiter string) (outputs []string, err error) {
	input_withoutQuot := Util_ClearInQuot(input)

	var in bool
	var buf string
	outputs = make([]string, 0, 10)

	for i := 0; i < len(input); i++ {
		at := input[i : i+1]
		quat := input_withoutQuot[i : i+1]

		if quat == delimiter {
			if in == true {
				for _, name := range outputs {
					if name == buf {
						return nil, fmt.Errorf("duplicate field name | field name : %s", name)
					}
				}
				outputs = append(outputs, buf)
				buf = ""
			}
			in = !in
			continue
		}
		if in == true {
			buf += at
		}
	}
	return outputs, nil
}

func Util_ReplaceBetweenDelimiter(input string, delimiter string, delimiterAfter string) (output string) {
	input_withoutQuot := Util_ClearInQuot(input)

	var in bool
	for i := 0; i < len(input); i++ {
		at := input[i : i+1]
		quat := input_withoutQuot[i : i+1]

		if quat == delimiter {
			if in == true {
				output += delimiterAfter
			}
			in = !in
			continue
		}
		if in == false {
			output += at
		}
	}
	return output
}

func Util_ClearDelimiter(input string, delimiter string) (output string) {
	input_withoutQuot := Util_ClearInQuot(input)
	for i := 0; i < len(input); i++ {
		at := input[i : i+1]
		if input_withoutQuot[i:i+1] == delimiter {
			continue
		}
		output += at
	}
	return output
}

func Util_ReplaceInDelimiter(input string, delimiter string, spliter string) (output string) {
	//  xxxx#AAAA#xxxx 		-> xxxxAAAAxxxx
	//  xxxx#AAAA/BBBB#xxxx -> xxxxBBBBxxxx
	input_withoutQuot := Util_ClearInQuot(input)

	var buf string
	var in bool
	for i := 0; i < len(input); i++ {
		at := input[i : i+1]
		quat := input_withoutQuot[i : i+1]

		if quat == delimiter {
			in = !in
			if in == true { // 구분자 진입 시점
				buf = ""
			} else { // 구분자 탈출 시점
				output += buf
			}
			continue
		}

		if in == false {
			output += at
		} else {
			if quat == spliter {
				buf = ""
				continue
			}
			buf += at
		}
	}
	return output
}

// insert into `aaa` values (a, b, c, d) -> a, b, c, d 추출
func Util_ExportInsertQueryValues(sqlInsert string) string {
	sqlInsertLower := strings.ToLower(sqlInsert)
	idxValues := strings.Index(sqlInsertLower, "values")
	if idxValues == -1 {
		return ""
	}

	var start, end int
	for i := idxValues; i < len(sqlInsert); i++ {
		if sqlInsert[i:i+1] == "(" {
			start = i
		}
		if sqlInsert[i:i+1] == ")" {
			end = i
		}
	}
	return sqlInsert[start+1 : end]
}

func Util_ConvFirstToUpper(input string) (output string) {
	if len(input) == 0 {
		return ""
	}

	output = strings.ToUpper(input[0:1]) + input[1:]
	return output
}

func Util_IsParserValArg(val []byte) bool {
	s_val := string(val)
	if len(s_val) < 3 { // :v0 | :v1 | :vN  임으로 3보다 작으면 추출 할 대상이 아님
		return false
	}
	if s_val[0:2] != ":v" { // ? 종류인지 구분 = 아니면 무시
		return false
	}
	return true
}
