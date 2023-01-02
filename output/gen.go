// Code generated - DO NOT EDIT.
// This file is a generated and any changes will be lost.

package gen

import (
	"fmt"
	. "github.com/gokch/ornn/db"
)

type Schema struct {
	test   Test
	user   User
	custom Custom
}

func (t *Schema) Init(
	job *Job,
) {
	t.test.Init(job)
	t.user.Init(job)
	t.custom.Init(job)
}

func (t *Test) Init(
	job *Job,
) {
	t.Job = job
}

type Test struct {
	Job *Job
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
	
	exec, err := t.Job.Exec(
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
	ret, err := t.Job.Query(
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
			
	exec, err := t.Job.Exec(
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
	t.Job = job
}

type User struct {
	Job *Job
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

