// Code generated - DO NOT EDIT.
// This file is a generated and any changes will be lost.

package gen

import (
	"fmt"
	"github.com/gokch/ornn/db"
)

type Schema struct {
	user User
}

func (t *Schema) Init(
	job *Job,
) {
	t.user.Init(job)
}

func (t *User) Init(
	job *Job,
) {
	t.Job = job
}

type User struct {
	Job *Job
}

func (t *User) Insert(
	arg_seq *int64,
	arg_id *string,
	arg_name *string,
) (
	lastInsertId int64,
	err error,
) {
	args := make([]interface{}, 0, 3)
	args = append(args, I_to_arri(
		arg_seq,
		arg_id,
		arg_name,
	)...)
	
	sql := fmt.Sprintf(
		"INSERT INTO user VALUES (?, ?, ?)",
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

type User_select struct {
	Seq  int64
	Id   string
	Name string
}

func (t *User) Select() (
	selects []*User_select,
	err error,
) {
	args := make([]interface{}, 0, 0)
	args = append(args, I_to_arri()...)
	
	sql := fmt.Sprintf(
		"SELECT * FROM user",
	)
	ret, err := t.Job.Query(
		sql,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer ret.Close()
	
	selects = make([]*User_select, 0, 100)
	for ret.Next() {
		pt_struct := &User_select{}
		err := ret.Scan(pt_struct)
		if err != nil {
			return nil, err
		}
		selects = append(selects, pt_struct)
	}
	
	return selects, nil
}

func (t *User) Delete() (
	rowAffected int64,
	err error,
) {
	args := make([]interface{}, 0, 0)
	args = append(args, I_to_arri()...)
	
	sql := fmt.Sprintf(
		"DELETE FROM user",
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

