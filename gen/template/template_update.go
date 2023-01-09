package template

import (
	"fmt"
	"strings"
)

func Update(args []string, tpls []string, query string, nullIgnore bool, structName, instanceName string) string {
	var body string
	if nullIgnore == true {
		body = genQuery_body_removeSets(args)
	} else {
		body = genQuery_body_setArgs(args)
	}

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

func genQuery_body_removeSets(args []string) (removeSets string) {
	var isNil string
	for _, arg := range args {
		fieldName := strings.TrimPrefix(arg, "arg_")

		isNil += fmt.Sprintf(`if %s == nil {
	setsRemoved = append(setsRemoved, "%s")
} else {
	args = append(args, %s)
}
`,
			arg,
			fieldName,
			arg,
		)
	}

	removeSets = fmt.Sprintf(`
args := make([]interface{}, 0, %d)
setsRemoved := make([]string, 0, %d)
%s
if len(setsRemoved) != 0 {
	sql, _ = RemoveNull(sql, setsRemoved)
}
`,
		len(args),
		len(args),
		isNil,
	)
	return removeSets
}
