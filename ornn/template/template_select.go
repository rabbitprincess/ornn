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
	return fmt.Sprintf(`
%s
sql := fmt.Sprintf(
	"%s",%s
)
ret, err := %s.%s.Query(
	sql,
	args...,
)
if err != nil {
	return nil, err
}
defer ret.Close()
%s
for ret.Next() {
	scan := &%s{}
	err := ret.Scan(scan)
	if err != nil {
		return nil, err
	}
	%s
}

return %s, nil
`,
		genQuery_body_setArgs(args),
		query,
		genQuery_body_arg(tpls),
		structName,
		instanceName,
		bodyRetDeclare,
		retName,
		bodyRetSet,
		retItemName,
	)
}
