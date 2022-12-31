package template

import "fmt"

func Insert(arg, query string, tpls []string, multiInsert, structName, instanceName string) string {
	return fmt.Sprintf(`
%s
sql := fmt.Sprintf(
	"%s",%s%s
)

exec, err := %s.%s.Exec(
	sql,
	args...,
)
if err != nil {
	return 0, err
}

return exec.LastInsertId()
`,
		arg,
		query,
		genQuery_body_arg(tpls),
		multiInsert,
		structName,
		instanceName,
	)
}

func genQuery_body_arg(args []string) (ret string) {
	for _, arg := range args {
		ret += fmt.Sprintf("\n\t%s,", arg)
	}
	return ret
}
