// Code generated - DO NOT EDIT.
// This file is a generated and any changes will be lost.

package gen

import (
	"fmt"
	. "github.com/gokch/ornn/db"
)

type Schema struct {
	authors Authors
}

func (t *Schema) Init(
	job *Job,
) {
	t.authors.Init(job)
}

func (t *Authors) Init(
	job *Job,
) {
	t.job = job
}

type Authors struct {
	job *Job
}

type Authors__select struct {
	Id   int
	Name string
	Bio   string
}

func (t *Authors) Select() (
	arrpt_select []*Authors__select,
	err error,
) {
	arri_arg := make([]interface{}, 0, 0)
	arri_arg = append(arri_arg)
	
	s_sql := fmt.Sprintf(
		"SELECT * FROM authors",
	)
	pc_ret, err := t.job.Query(
		s_sql,
		arri_arg...,
	)
	if err != nil {
		return nil, err
	}
	defer pc_ret.Close()
	
	arrpt_select = make([]*Authors__select, 0, 100)
	for pc_ret.Next(){
		pt_struct := &Authors__select{}
		err := pc_ret.Scan(pt_struct)
		if err != nil {
			return nil, err
		}
		arrpt_select = append(arrpt_select, pt_struct)
	}
	
	return arrpt_select, nil
}