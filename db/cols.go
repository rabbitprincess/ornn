package db

func I_to_arri(_arri ...interface{}) (arri_ret []interface{}) {
	arri_ret = make([]interface{}, 0, 10)

	for _, i_item := range _arri {
		arri_ret = append(arri_ret, i_item)
	}

	return arri_ret
}
