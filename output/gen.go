// Code generated - DO NOT EDIT.
// This file is a generated and any changes will be lost.

package gen

import (
	"fmt"
	. "github.com/gokch/ornn/db"
)

type Gen struct {
	Tbltest Tbltest
	User    User
}

func (t *Gen) Init(
	job *Job,
) {
	t.Tbltest.Init(job)
	t.User.Init(job)
}

func (t *Tbltest) Init(
	job *Job,
) {
	t.job = job
}

type Tbltest struct {
	job *Job
}

func (t *Tbltest) Insert(
	arg_seq uint32,
	arg_id2 string,
	arg_address []byte,
	arg_registered bool,
) (
	lastInsertId int64,
	err error,
) {
	args := make([]interface{}, 0, 4)
	args = append(args, I_to_arri(
		arg_seq,
		arg_id2,
		arg_address,
		arg_registered,
	)...)
	
	sql := fmt.Sprintf(
		"",
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

func (t *User) Init(
	job *Job,
) {
	t.job = job
}

type User struct {
	job *Job
}

func (t *User) Insert(
	arg_seq uint32,
	arg_id2 string,
	arg_address []byte,
	arg_registered bool,
) (
	lastInsertId int64,
	err error,
) {
	args := make([]interface{}, 0, 4)
	args = append(args, I_to_arri(
		arg_seq,
		arg_id2,
		arg_address,
		arg_registered,
	)...)
	
	sql := fmt.Sprintf(
		"",
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

