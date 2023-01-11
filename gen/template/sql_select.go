package template

import (
	"fmt"
)

func Select(args []string, tpls []string, query string, selectSingle bool, structName string, instanceName string, retName, retItemName, retItemType string) string {
	var bodyRetDeclare, bodyRetSet string
	if selectSingle == true {
		bodyRetSet = fmt.Sprintf("%s = scan\n\tbreak", retItemName)
	} else {
		bodyRetDeclare = fmt.Sprintf("\n%s = make(%s, 0, 100)", retItemName, retItemType)
		bodyRetSet = fmt.Sprintf("%s = append(%s, scan)", retItemName, retItemName)
	}
	return parseTemplate("sql_select.template", map[string]interface{}{
		"arg":      genQuery_body_setArgs(args),
		"query":    query,
		"tpl":      genQuery_body_arg(tpls),
		"struct":   structName,
		"instance": instanceName,
		"body":     bodyRetDeclare,
		"scan":     retName,
		"retSet":   bodyRetSet,
		"ret":      retItemName,
	})
}
