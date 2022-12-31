package db

func I_to_arri(items ...interface{}) (arrs []interface{}) {
	arrs = make([]interface{}, 0, 10)

	for _, item := range items {
		arrs = append(arrs, item)
	}

	return arrs
}
