package template

import (
	"fmt"
	"strings"

	"github.com/gokch/ornn/gen/sql"
)

func Insert(args []string, tpls []string, query string, insertMulti bool, structName, instanceName string) string {
	var multiInsert, genArgs string
	if insertMulti == true { // multi insert
		queryVal := sql.Util_ExportInsertQueryValues(query)
		query = strings.TrimSuffix(query, ";")
		query += "%s"
		genArgs = genQuery_body_multiInsertProc(args)
		multiInsert = genQuery_body_multiInsert(queryVal)
	} else { // insert
		genArgs = genQuery_body_setArgs(args)
	}
	return parseTemplate("sql_insert.template", map[string]interface{}{
		"arg":      genArgs,
		"query":    query,
		"tpl":      genQuery_body_arg(tpls),
		"multi":    multiInsert,
		"struct":   structName,
		"instance": instanceName,
	})
}

func genQuery_body_multiInsertProc(args []string) (multiInsertProc string) {
	var checkLen string
	for i, arg := range args {
		checkLen += fmt.Sprintf("argLen != len(%s)", arg)
		if i != len(args)-1 {
			checkLen += fmt.Sprintf(" || ")
		}
	}

	var append string
	for i, arg := range args {
		append += fmt.Sprintf("%s[i]", arg)
		if i != len(args)-1 {
			append += fmt.Sprintf(",\n\t\t")
		}
	}

	multiInsertProc = fmt.Sprintf(`argLen := len(%s)
if argLen == 0 {
	return 0, fmt.Errorf("arg len is zero")
}
if %s {
	return 0, fmt.Errorf("arg len is not same")
}

args := make([]interface{}, 0, argLen*%d)
for i := 0; i < argLen; i++ {
	args = append(args, I_to_arri(
		%s,
	)...)
}
`,
		args[0],
		checkLen,
		len(args),
		append)
	return multiInsertProc
}

func genQuery_body_multiInsert(query string) (multiInsert string) {
	return fmt.Sprintf("\n\tstrings.Repeat(\", (%s)\", argLen-1),", query)
}
