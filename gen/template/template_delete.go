package template

import "fmt"

func Delete(args []string, query string, tpls []string, structName, instanceName string) string {
	return fmt.Sprintf(`
%s
sql := fmt.Sprintf(
	"%s",%s
)
		
exec, err := %s.%s.Exec(
	sql,
	args...,
)
if err != nil {
	return 0, err
}

return exec.RowsAffected()
`,
		genQuery_body_setArgs(args),
		query,
		genQuery_body_arg(tpls),
		structName,
		instanceName,
	)
}