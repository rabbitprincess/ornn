package template

import "fmt"

func Update(query string, tpls []string, body, structName, instanceName string) string {
	return fmt.Sprintf(`
sql := fmt.Sprintf(
	"%s",%s
)
%s
exec, err := %s.%s.Exec(
	sql,
	args...,
)
if err != nil {
	return 0, err
}

return exec.RowsAffected()
`,
		query,
		genQuery_body_arg(tpls),
		body,
		structName,
		instanceName,
	)
}
