// Code generated - DO NOT EDIT.
// This file is a generated and any changes will be lost.

package gen

import (
	"fmt"
	. "github.com/gokch/ornn/db"
)

type Gen struct {
	User User
}

func (t *Gen) Init(
	job *Job,
) {
	t.User.Init(job)
}

func (t *User) Init(
	job *Job,
) {
	t.job = job
}

type User struct {
	job *Job
}

func (t *User) Insert(
	arg_seq int64,
	arg_id string,
	arg_name string,
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
	
	exec, err := t.job.Exec(
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
	ret, err := t.job.Query(
		sql,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer ret.Close()
	
	selects = make([]*User_select, 0, 100)
	for ret.Next() {
		scan := &User_select{}
		err := ret.Scan(scan)
		if err != nil {
			return nil, err
		}
		selects = append(selects, scan)
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
			
	exec, err := t.job.Exec(
		sql,
		args...,
	)
	if err != nil {
		return 0, err
	}
	
	return exec.RowsAffected()
}

