// Code generated - DO NOT EDIT.
// This file is a generated and any changes will be lost.

package gen

import (
	"fmt"
	. "github.com/gokch/ornn/db"
)

type Schema struct {
	mytesttable MyTestTable
	custom      Custom
}

func (t *Schema) Init(
	job *Job,
) {
	t.mytesttable.Init(job)
	t.custom.Init(job)
}

func (t *MyTestTable) Init(
	job *Job,
) {
	t.Job = job
}

type MyTestTable struct {
	Job *Job
}

func (t *MyTestTable) Insert(
	arg_id *uint32,
) (
	lastInsertId int64,
	err error,
) {
	args := make([]interface{}, 0, 1)
	args = append(args, I_to_arri(
		arg_id,
	)...)
	
	sql := fmt.Sprintf(
		"INSERT INTO myTestTable VALUES (?)",
	)
	
	exec, err := t.Job.Exec(
		sql,
		args...,
	)
	if err != nil {
		return 0, err
	}
	
	return exec.LastInsertId()
}

type MyTestTable_select struct {
	Id interface{}
}

func (t *MyTestTable) Select() (
	selects []*MyTestTable_select,
	err error,
) {
	args := make([]interface{}, 0, 0)
	args = append(args, I_to_arri()...)
	
	sql := fmt.Sprintf(
		"SELECT * FROM myTestTable",
	)
	ret, err := t.Job.Query(
		sql,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer ret.Close()
	
	selects = make([]*MyTestTable_select, 0, 100)
	for ret.Next() {
		scan := &MyTestTable_select{}
		err := ret.Scan(scan)
		if err != nil {
			return nil, err
		}
		selects = append(selects, scan)
	}
	
	return selects, nil
}

func (t *MyTestTable) Delete() (
	rowAffected int64,
	err error,
) {
	args := make([]interface{}, 0, 0)
	args = append(args, I_to_arri()...)
	
	sql := fmt.Sprintf(
		"DELETE FROM myTestTable",
	)
			
	exec, err := t.Job.Exec(
		sql,
		args...,
	)
	if err != nil {
		return 0, err
	}
	
	return exec.RowsAffected()
}

func (t *Custom) Init(
	job *Job,
) {
	t.Job = job
}

type Custom struct {
	Job *Job
}

