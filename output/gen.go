// Code generated - DO NOT EDIT.
// This file is a generated and any changes will be lost.

package gen

import (
	"fmt"
	. "github.com/gokch/ornn/db"
)

type Gen struct {
	Newtable Newtable
}

func (t *Gen) Init(
	job *Job,
) {
	t.Newtable.Init(job)
}

func (t *Newtable) Init(
	job *Job,
) {
	t.job = job
}

type Newtable struct {
	job *Job
}

func (t *Newtable) Insert(
	arg_a string,
	arg_b string,
) (
	lastInsertId int64,
	err error,
) {
	args := make([]interface{}, 0, 2)
	args = append(args, I_to_arri(
		arg_a,
		arg_b,
	)...)
	
	sql := fmt.Sprintf(
		"INSERT INTO newtable VALUES (?, ?)",
	)
	
	exec, err := t.job.Exec(
		sql,
		args...,
	)
	if err != nil {
		return 0, err
	}
	
	return exec.LastInsertId()
}

type Newtable_select struct {
	B string
	A string
}

func (t *Newtable) Select() (
	selects []*Newtable_select,
	err error,
) {
	args := make([]interface{}, 0, 0)
	args = append(args, I_to_arri()...)
	
	sql := fmt.Sprintf(
		"SELECT * FROM newtable",
	)
	ret, err := t.job.Query(
		sql,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer ret.Close()
	
	selects = make([]*Newtable_select, 0, 100)
	for ret.Next() {
		scan := &Newtable_select{}
		err := ret.Scan(scan)
		if err != nil {
			return nil, err
		}
		selects = append(selects, scan)
	}
	
	return selects, nil
}

func (t *Newtable) Delete() (
	rowAffected int64,
	err error,
) {
	args := make([]interface{}, 0, 0)
	args = append(args, I_to_arri()...)
	
	sql := fmt.Sprintf(
		"DELETE FROM newtable",
	)
			
	exec, err := t.job.Exec(
		sql,
		args...,
	)
	if err != nil {
		return 0, err
	}
	
	return exec.RowsAffected()
}

