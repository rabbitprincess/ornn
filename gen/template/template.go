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
