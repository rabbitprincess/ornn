package codegen

import (
	"log"
	"strings"
)

type Item interface {
	Code(w *Writer)
}

//--------------------------------------------------------------------------------------------------------------//
// global

type Global struct {
	DoNotEdit string
	Package   string
	Imports   Import
	Items     []Item
}

func (t *Global) Init() {
	t.Items = make([]Item, 0, 10)
}

func (t *Global) AddImport(item *ImportItem) {
	t.Imports.Add(item)
}

func (t *Global) AddItem(i Item) {
	t.Items = append(t.Items, i)
}

func (t *Global) Code(w *Writer) {
	w.W("%s\n", t.DoNotEdit)         // comment ( do not edit )
	w.W("package %s\n\n", t.Package) // package
	t.Imports.Code(w)                // import

	for _, pt := range t.Items {
		pt.Code(w)
	}
}

//--------------------------------------------------------------------------------------------------------------//
// struct

type Struct struct {
	Name    string
	Field   Var
	Methods []*Func
}

func (t *Struct) Init() {
	t.Field.Scope = VarScopeStructField
	t.Methods = make([]*Func, 0, 10)
}

func (t *Struct) AddField(item *VarItem) {
	t.Field.Add(item)
}

func (t *Struct) AddFunc(item *Func) {
	t.Methods = append(t.Methods, item)
}

func (t *Struct) Code(w *Writer) {
	// struct
	w.W("type %s struct {\n", t.Name)
	w.IndentIn()
	t.Field.Code(w)
	w.IndentOut()
	w.W("}\n\n")

	// func
	for _, method := range t.Methods {
		method.Code(w)
	}
}

//--------------------------------------------------------------------------------------------------------------//
// func

type Func struct {
	StructName string
	StructType string
	FuncName   string
	InlineCode string

	Arg   Var
	Ret   Var
	Const Const
}

func (t *Func) Init() {
	t.Arg.Scope = VarScopeFuncArg
	t.Ret.Scope = VarScopeFuncRet
}

func (t *Func) AddArg(item *VarItem) {
	t.Arg.Add(item)
}

func (t *Func) AddRet(item *VarItem) {
	t.Ret.Add(item)
}

func (t *Func) AddConst(item *ConstItem) {
	t.Const.Add(item)
}

func (t *Func) Code(w *Writer) {
	w.W("func ")

	if t.StructType != "" {
		w.N("(%s %s) ", t.StructName, t.StructType)
	}

	w.N("%s(", t.FuncName)
	if len(t.Arg.Items) != 0 {
		w.N("\n")
		w.IndentIn()
		t.Arg.Code(w)
		w.IndentOut()
	}
	w.W(")")

	if len(t.Ret.Items) != 0 {
		w.N(" ") // 인자와 리턴() 간에 빈칸 1개 추가
		t.Ret.Code(w)
	}

	// 함수 내용 출력
	w.N(" {\n")
	w.IndentIn()
	t.Const.Code(w)

	// 코드 첫 빈줄 제거
	if len(t.InlineCode) != 0 && t.InlineCode[0:1] == "\n" {
		t.InlineCode = t.InlineCode[1:]
	}

	// 코드 끝 빈줄 제거
	if len(t.InlineCode) != 0 && t.InlineCode[len(t.InlineCode)-1:] == "\n" {
		t.InlineCode = t.InlineCode[:len(t.InlineCode)-1]
	}

	w.W("%s\n", strings.ReplaceAll(t.InlineCode, "\n", "\n"+w.Indent()))
	w.IndentOut()
	w.W("}\n\n")
}

//--------------------------------------------------------------------------------------------------------------//
// var

type VarScope int8

const (
	VarScopeGlobal VarScope = iota + 1
	VarScopeStructField
	VarScopeFuncArg
	VarScopeFuncRet
)

type Var struct {
	Scope VarScope
	Items []*VarItem
}

func (t *Var) Add(item *VarItem) {
	if t.Items == nil {
		t.Items = make([]*VarItem, 0, 10)
	}
	t.Items = append(t.Items, item)
}

func (t *Var) MaxLenName() int {
	var max int
	for _, pt := range t.Items {
		if max < len(pt.Name) {
			max = len(pt.Name)
		}
	}
	return max
}

func (t *Var) Code(w *Writer) {
	switch t.Scope {
	case VarScopeGlobal:
		{
			w.W("var (\n")
			w.IndentIn()
			for _, pt := range t.Items {
				w.W("")
				pt.Code(w, 0)
				w.N("\n")
			}
			w.IndentOut()
			w.W(")\n")
		}
	case VarScopeStructField:
		{
			max := t.MaxLenName()
			for _, pt := range t.Items {
				blank := max - len(pt.Name)
				w.W("")
				pt.Code(w, blank)
				w.N("\n")
			}
		}
	case VarScopeFuncArg:
		{
			for i := 0; i < len(t.Items); i++ {
				w.W("")
				t.Items[i].Code(w, 0)
				w.N(",\n")
			}
		}
	case VarScopeFuncRet:
		{
			lenItem := len(t.Items)
			if lenItem == 0 {
				// 출력 할 것이 없음
				break
			}

			if lenItem == 1 && t.Items[0].Name == "" {
				t.Items[0].Code(w, 0)
				// 더 출력 할 것이 없음 = 단일 리턴값의 type 만 출력 (name 이 없음)
				break
			}

			w.N("(\n")
			w.IndentIn()
			for i := 0; i < len(t.Items); i++ {
				w.W("")
				t.Items[i].Code(w, 0)
				w.N(",\n")
			}
			w.IndentOut()
			w.W(")")
		}
	default:
		log.Fatalf("input undefined type of scope - %v", t.Scope)
	}
}

//--------------------------------------------------------------------------------------------------------------//
// var - item

type VarItem struct {
	Name string
	Type string
}

func (t *VarItem) Code(w *Writer, blank int) {
	if t.Name != "" {
		w.N("%s%s %s", t.Name, strings.Repeat(" ", blank), t.Type)
	} else {
		w.N("%s", t.Type)
	}
}

//--------------------------------------------------------------------------------------------------------------//
// const

type Const struct {
	items []*ConstItem
}

func (t *Const) Add(item *ConstItem) {
	if t.items == nil {
		t.items = make([]*ConstItem, 0, 10)
	}
	t.items = append(t.items, item)
}

func (t *Const) Code(w *Writer) {
	if len(t.items) > 0 {
		w.W("const (\n")
		w.IndentIn()
		for _, pt := range t.items {
			pt.Code(w)
		}
		w.IndentOut()
		w.W(")\n")
	}
}

//--------------------------------------------------------------------------------------------------------------//
// const - item

type ConstItem struct {
	Name  string
	Type  string // type 이 없으면 value 로 형이 결정 됨
	Value string // value 가 없으면 상위 개체가 iota 일 수 있음
}

func (t *ConstItem) Code(w *Writer) {
	w.W("%s %s = %s\n", t.Name, t.Type, t.Value)
}

//--------------------------------------------------------------------------------------------------------------//
// import

type Import struct {
	items []*ImportItem
}

func (t *Import) Add(item *ImportItem) {
	if t.items == nil {
		t.items = make([]*ImportItem, 0, 10)
	}
	t.items = append(t.items, item)
}

func (t *Import) Code(w *Writer) {
	if len(t.items) > 0 {
		w.W("import (\n")
		w.IndentIn()
		for _, pt := range t.items {
			pt.Code(w)
		}
		w.IndentOut()
		w.W(")\n\n")
	}
}

//--------------------------------------------------------------------------------------------------------------//
// import - item

type ImportItem struct {
	Path  string
	Alias string // 기본 ""(없음)  or . or  임의이름
}

func (t *ImportItem) Code(w *Writer) {
	if t.Alias != "" {
		w.W("%s \"%s\"\n", t.Alias, t.Path)
	} else {
		w.W("\"%s\"\n", t.Path)
	}
}
