package template

func Delete(args []string, query string, tpls []string, structName, instanceName string) string {
	return parseTemplate("sql_delete.template", map[string]interface{}{
		"arg":      genQuery_body_setArgs(args),
		"query":    query,
		"tpl":      genQuery_body_arg(tpls),
		"struct":   structName,
		"instance": instanceName,
	})
}
