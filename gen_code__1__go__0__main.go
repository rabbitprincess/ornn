package bp

type T_Gen__go struct {
	T_global T_Go__global
}

func (t *T_Gen__go) Code() (s_code string, err error) { // for I_Gen interface
	pt_w := &Gen_writer{}
	pt_w.Init()

	t.T_global.Code(pt_w)

	return pt_w.String(), nil
}
