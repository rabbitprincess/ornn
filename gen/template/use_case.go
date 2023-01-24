package template

func UseCase(packageName, className string) string {
	return parseTemplate("use_case.template", map[string]interface{}{
		"package": packageName,
		"class":   className,
	})
}
