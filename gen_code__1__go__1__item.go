package bp

import (
	"module/debug/logf"
	"strings"
)

type I_Gen__global__item interface {
	I_Gen_go__global__item()
	Code(*Gen_writer)
}

//--------------------------------------------------------------------------------------------------------------//
// global

type T_Go__global struct {
	S_comment__do_not_edit string
	S_package_name         string
	t_import               T_Go__import
	arri                   []I_Gen__global__item
}

func (t *T_Go__global) Set_package_name(_s_package_name string) {
	t.S_package_name = _s_package_name
}

func (t *T_Go__global) Add(_i I_Gen__global__item) {
	if t.arri == nil {
		t.arri = make([]I_Gen__global__item, 0, 10)
	}

	t.arri = append(t.arri, _i)
}

func (t *T_Go__global) Code(_pt_w *Gen_writer) {
	_pt_w.W("%s\n", t.S_comment__do_not_edit)   // comment ( do not edit )
	_pt_w.W("package %s\n\n", t.S_package_name) // package
	t.t_import.Code(_pt_w)                      // import

	for _, pt := range t.arri {
		pt.Code(_pt_w)
	}
}

//--------------------------------------------------------------------------------------------------------------//
// struct

type T_Go__struct struct {
	S_struct_name string
	t_var_field   T_Go__var
	arrpt_func    []*T_Go__func
}

func (t *T_Go__struct) I_Gen_go__global__item() {}

func (t *T_Go__struct) Add_field(_pt_item *T_Go__var__item) {
	t.t_var_field.Add(_pt_item)
}

func (t *T_Go__struct) Add_func(_pt_item *T_Go__func) {
	if t.arrpt_func == nil {
		t.arrpt_func = make([]*T_Go__func, 0, 10)
	}

	t.arrpt_func = append(t.arrpt_func, _pt_item)
}

func (t *T_Go__struct) Code(_pt_w *Gen_writer) {
	// set struct
	_pt_w.W("type %s struct {\n", t.S_struct_name)
	_pt_w.Indent__in()
	t.t_var_field.Code(_pt_w)
	_pt_w.Indent__out()
	_pt_w.W("}\n\n")

	// set func
	for _, pt := range t.arrpt_func {
		pt.Code(_pt_w)
	}
}

//--------------------------------------------------------------------------------------------------------------//
// func

type T_Go__func struct {
	S_struct_name  string
	S_struct_type  string
	S_func_name    string
	S_in_body_code string

	t_var_arg    T_Go__var
	t_var_return T_Go__var
	t_const      T_Go__const
}

func (t *T_Go__func) I_Gen_go__global__item() {}

func (t *T_Go__func) Add_arg(_pt_item *T_Go__var__item) {
	t.t_var_arg.Add(_pt_item)
}

func (t *T_Go__func) Add_return(_pt_item *T_Go__var__item) {
	t.t_var_return.Add(_pt_item)
}

func (t *T_Go__func) Add_const(_pt_item *T_Go__const__item) {
	t.t_const.Add(_pt_item)
}

func (t *T_Go__func) Code(_pt_w *Gen_writer) {
	_pt_w.W("func ")

	if t.S_struct_type != "" {
		_pt_w.N("(%s %s) ", t.S_struct_name, t.S_struct_type)
	}

	_pt_w.N("%s(", t.S_func_name)
	if len(t.t_var_arg.arrpt) != 0 {
		_pt_w.N("\n")
		_pt_w.Indent__in()
		t.t_var_arg.Code(_pt_w)
		_pt_w.Indent__out()
	}
	_pt_w.W(")")

	if len(t.t_var_return.arrpt) != 0 {
		_pt_w.N(" ") // 인자와 리턴() 간에 빈칸 1개 추가
		t.t_var_return.Code(_pt_w)
	}

	// 함수 내용 출력
	_pt_w.N(" {\n")
	_pt_w.Indent__in()
	t.t_const.Code(_pt_w)

	// 코드 첫 빈줄 제거
	if len(t.S_in_body_code) != 0 && t.S_in_body_code[0:1] == "\n" {
		t.S_in_body_code = t.S_in_body_code[1:]
	}

	// 코드 끝 빈줄 제거
	if len(t.S_in_body_code) != 0 && t.S_in_body_code[len(t.S_in_body_code)-1:] == "\n" {
		t.S_in_body_code = t.S_in_body_code[:len(t.S_in_body_code)-1]
	}

	_pt_w.W("%s\n", strings.ReplaceAll(t.S_in_body_code, "\n", "\n"+_pt_w.Get_indent()))
	_pt_w.Indent__out()
	_pt_w.W("}\n\n")
}

//--------------------------------------------------------------------------------------------------------------//
// var

type TD_N1_code_format__var int8

const (
	TD_N1_code_format__var__in_global_var TD_N1_code_format__var = iota + 1
	TD_N1_code_format__var__in_struct_field
	TD_N1_code_format__var__in_func_arg
	TD_N1_code_format__var__in_func_return
)

type T_Go__var struct {
	td_n1_code_format TD_N1_code_format__var
	arrpt             []*T_Go__var__item
}

func (t *T_Go__var) I_Gen_go__global__item() {}

func (t *T_Go__var) Add(_pt_item *T_Go__var__item) {
	if t.arrpt == nil {
		t.arrpt = make([]*T_Go__var__item, 0, 10)
	}
	t.arrpt = append(t.arrpt, _pt_item)
}

func (t *T_Go__var) Get_max_length_of_name() int {
	var n_len__max int
	for _, pt := range t.arrpt {
		n_len := len(pt.S_var_name)
		if n_len__max < n_len {
			n_len__max = n_len
		}
	}
	return n_len__max
}

func (t *T_Go__var) Code(_pt_w *Gen_writer) {
	// _pt_w.Wf
	switch t.td_n1_code_format {
	case TD_N1_code_format__var__in_global_var:
		{
			_pt_w.W("var (\n")
			_pt_w.Indent__in()
			for _, pt := range t.arrpt {
				_pt_w.W("")
				pt.Code(_pt_w, 0)
				_pt_w.N("\n")
			}
			_pt_w.Indent__out()
			_pt_w.W(")\n")
		}
	case TD_N1_code_format__var__in_struct_field:
		{
			n_len__max := t.Get_max_length_of_name()

			for _, pt := range t.arrpt {
				n_len__blank := n_len__max - len(pt.S_var_name)
				_pt_w.W("")
				pt.Code(_pt_w, n_len__blank)
				_pt_w.N("\n")
			}
		}
	case TD_N1_code_format__var__in_func_arg:
		{
			for i := 0; i < len(t.arrpt); i++ {
				_pt_w.W("")
				t.arrpt[i].Code(_pt_w, 0)
				_pt_w.N(",\n")
			}
		}
	case TD_N1_code_format__var__in_func_return:
		{
			n_len := len(t.arrpt)
			if n_len == 0 {
				// 출력 할 것이 없음
				break
			}

			if n_len == 1 && t.arrpt[0].S_var_name == "" {
				t.arrpt[0].Code(_pt_w, 0)
				// 더 출력 할 것이 없음 = 단일 리턴값의 type 만 출력 (name 이 없음)
				break
			}

			_pt_w.N("(\n")
			_pt_w.Indent__in()
			for i := 0; i < len(t.arrpt); i++ {
				_pt_w.W("")
				t.arrpt[i].Code(_pt_w, 0)
				_pt_w.N(",\n")
			}
			_pt_w.Indent__out()
			_pt_w.W(")")
		}
	default:
		logf.GC.Fatal("bp", "input undefined type of td_n1_code_format - %d", t.td_n1_code_format)
	}
}

//--------------------------------------------------------------------------------------------------------------//
// var - item

type T_Go__var__item struct {
	S_var_name string
	S_var_type string
}

func (t *T_Go__var__item) Code(_pt_w *Gen_writer, _n_len__blank int) {
	if t.S_var_name != "" {
		_pt_w.N("%s%s %s", t.S_var_name, strings.Repeat(" ", _n_len__blank), t.S_var_type)
	} else {
		_pt_w.N("%s", t.S_var_type)
	}
}

//--------------------------------------------------------------------------------------------------------------//
// const

type T_Go__const struct {
	arrpt []*T_Go__const__item
}

func (t *T_Go__const) I_Gen_go__global__item() {}

func (t *T_Go__const) Add(_pt_item *T_Go__const__item) {
	if t.arrpt == nil {
		t.arrpt = make([]*T_Go__const__item, 0, 10)
	}
	t.arrpt = append(t.arrpt, _pt_item)
}

func (t *T_Go__const) Code(_pt_w *Gen_writer) {
	if len(t.arrpt) > 0 {
		_pt_w.W("const (\n")
		_pt_w.Indent__in()
		for _, pt := range t.arrpt {
			pt.Code(_pt_w)
		}
		_pt_w.Indent__out()
		_pt_w.W(")\n")
	}
}

//--------------------------------------------------------------------------------------------------------------//
// const - item

type T_Go__const__item struct {
	S_const_name  string
	S_const_type  string // type 이 없으면 value 로 형이 결정 됨
	S_const_value string // value 가 없으면 상위 개체가 iota 일 수 있음
}

func (t *T_Go__const__item) Code(_pt_w *Gen_writer) {
	_pt_w.W("%s %s = %s\n", t.S_const_name, t.S_const_type, t.S_const_value)
}

//--------------------------------------------------------------------------------------------------------------//
// import

type T_Go__import struct {
	arrpt []*T_Go__import__item
}

func (t *T_Go__import) Add(_pt_item *T_Go__import__item) {
	if t.arrpt == nil {
		t.arrpt = make([]*T_Go__import__item, 0, 10)
	}
	t.arrpt = append(t.arrpt, _pt_item)
}

func (t *T_Go__import) Code(_pt_w *Gen_writer) {
	if len(t.arrpt) > 0 {
		_pt_w.W("import (\n")
		_pt_w.Indent__in()
		for _, pt := range t.arrpt {
			pt.Code(_pt_w)
		}
		_pt_w.Indent__out()
		_pt_w.W(")\n\n")
	}
}

//--------------------------------------------------------------------------------------------------------------//
// import - item

type T_Go__import__item struct {
	S_path     string
	S_use_name string // 기본 ""(없음)  or . or  임의이름
}

func (t *T_Go__import__item) Code(_pt_w *Gen_writer) {
	if t.S_use_name != "" {
		_pt_w.W("%s \"%s\"\n", t.S_use_name, t.S_path)
	} else {
		_pt_w.W("\"%s\"\n", t.S_path)
	}
}
