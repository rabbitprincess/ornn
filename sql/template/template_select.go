package template

import "fmt"

func Select(args []string, query string, tpls []string, structName string, instanceName string, bodyCodeDeclare, bodyCodeRet, retName, retItemName string) string {
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
		gen_query__add__func__body__set_args(args),
		query,
		gen_query__add__func__body__arg(tpls),
		structName,
		instanceName,
		bodyCodeDeclare,
		retName,
		bodyCodeRet,
		retItemName,
	)
}

func gen_query__add__func__body__set_args(_arrs_arg []string) (s_gen_arg string) {
	var s_gen_arg__item string
	s_gen_arg__item = gen_query__add__func__body__arg(_arrs_arg)
	if s_gen_arg__item != "" {
		s_gen_arg__item += "\n"
	}

	s_gen_arg += fmt.Sprintf(`arri_arg := make([]interface{}, 0, %d)
arri_arg = append(arri_arg, I_to_arri(%s)...)
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
