// Code generated - DO NOT EDIT.
// This file is a generated and any changes will be lost.

package gen

import (
	"fmt"
	. "github.com/gokch/ornn/db"
)

type Gen struct {
	Test   Test
	User   User
	Custom Custom
}

func (t *Gen) Init(
	job *Job,
) {
	t.Test.Init(job)
	t.User.Init(job)
	t.Custom.Init(job)
}

func (t *Test) Init(
	job *Job,
) {
	t.job = job
}

type Test struct {
	job *Job
}

func (t *Test) Insert(
	arg_seq *uint32,
	arg_id2 *string,
	arg_address *[]byte,
	arg_registered *int8,
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
		"INSERT INTO test VALUES (?, ?, ?, ?)",
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

type Test_select struct {
	Seq        uint32
	Id2        string
	Address    []byte
	Registered int8
}

func (t *Test) Select() (
	selects []*Test_select,
	err error,
) {
	args := make([]interface{}, 0, 0)
	args = append(args, I_to_arri()...)
	
	sql := fmt.Sprintf(
		"SELECT * FROM test",
	)
	ret, err := t.job.Query(
		sql,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer ret.Close()
	
	selects = make([]*Test_select, 0, 100)
	for ret.Next() {
		scan := &Test_select{}
		err := ret.Scan(scan)
		if err != nil {
			return nil, err
		}
		selects = append(selects, scan)
	}
	
	return selects, nil
}

func (t *Test) Delete() (
	rowAffected int64,
	err error,
) {
	args := make([]interface{}, 0, 0)
	args = append(args, I_to_arri()...)
	
	sql := fmt.Sprintf(
		"DELETE FROM test",
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

func (t *User) Init(
	job *Job,
) {
	t.job = job
}

type User struct {
	job *Job
}

func (t *User) Insert(
	arg_seq *uint32,
	arg_id2 *string,
	arg_address *[]byte,
	arg_registered *int8,
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
		"INSERT INTO user VALUES (?, ?, ?, ?)",
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
	Seq        uint32
	Id2        string
	Address    []byte
	Registered int8
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

func (t *Custom) Init(
	job *Job,
) {
	t.job = job
}

type Custom struct {
	job *Job
}

