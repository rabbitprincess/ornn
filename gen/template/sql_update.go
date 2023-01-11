package template

func Update(args []string, tpls []string, query string, structName, instanceName string) string {
	return parseTemplate("sql_update.template", map[string]interface{}{
		"query":    query,
		"tpl":      genQuery_body_arg(tpls),
		"arg":      genQuery_body_setArgs(args),
		"struct":   structName,
		"instance": instanceName,
	})

}
