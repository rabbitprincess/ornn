package template

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func parseTemplate(fileName string, args map[string]interface{}) string {
	path, _ := os.Getwd()
	parse, err := template.New(fileName).ParseFiles(path + "/gen/template/" + fileName)
	if err != nil {
		return ""
	}
	builder := &strings.Builder{}
	parse.Execute(builder, args)
	return builder.String()
}

func genQuery_body_setArgs(_arrs_arg []string) (genArg string) {
	var items string
	items = genQuery_body_arg(_arrs_arg)
	if items != "" {
		items += "\n"
	}

	return fmt.Sprintf(`args := []interface{}{%s}
`,
		items,
	)
}

func genQuery_body_arg(args []string) (ret string) {
	for _, arg := range args {
		ret += fmt.Sprintf("\n\t%s,", arg)
	}
	return ret
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
