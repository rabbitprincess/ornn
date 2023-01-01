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
		gen_query__add__func__body__arg(tpls),
		structName,
		instanceName,
		bodyRetDeclare,
		retName,
		bodyRetSet,
		retItemName,
	)
}

func genQuery_body_setArgs(_arrs_arg []string) (s_gen_arg string) {
	var s_gen_arg__item string
	s_gen_arg__item = gen_query__add__func__body__arg(_arrs_arg)
	if s_gen_arg__item != "" {
		s_gen_arg__item += "\n"
	}

	s_gen_arg += fmt.Sprintf(`args := make([]interface{}, 0, %d)
args = append(args, I_to_arri(%s)...)
`,
		len(_arrs_arg),
		s_gen_arg__item,
	)

	return s_gen_arg
}

func gen_query__add__func__body__arg(_arrs_arg []string) (s_arg string) {
	for _, s_arg__one := range _arrs_arg {
		s_arg += fmt.Sprintf("\n\t%s,", s_arg__one)
	}
	return s_arg
}
