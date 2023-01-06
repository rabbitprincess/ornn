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

type User_select struct {
	Seq  uint64
	Id   string
	Ord  int64
	Name string
	Pw   []byte
}

func (t *User) Select(
	arg_seq uint64,
) (
	selects []*User_select,
	err error,
) {
	args := []interface{}{
		arg_seq,
	}
	
	sql := fmt.Sprintf(
		"SELECT * FROM user WHERE seq = ?",
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

func (t *User) Update(
	arg_seq uint64,
	arg_id string,
	arg_ord int64,
	arg_name string,
	arg_pw []byte,
	arg_where_seq uint64,
) (
	rowAffected int64,
	err error,
) {
	sql := fmt.Sprintf(
		"UPDATE user SET seq = ?, id = ?, ord = ?, name = ?, pw = ? WHERE seq = ?",
	)
	args := []interface{}{
		arg_seq,
		arg_id,
		arg_ord,
		arg_name,
		arg_pw,
		arg_where_seq,
	}
	
	exec, err := t.job.Exec(
		sql,
		args...,
	)
	if err != nil {
		return 0, err
	}
	
	return exec.RowsAffected()
}

func (t *User) Delete(
	arg_seq uint64,
) (
	rowAffected int64,
	err error,
) {
	args := []interface{}{
		arg_seq,
	}
	
	sql := fmt.Sprintf(
		"DELETE FROM user WHERE seq = ?",
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

func (t *User) Insert(
	arg_seq uint64,
	arg_id string,
	arg_ord int64,
	arg_name string,
	arg_pw []byte,
) (
	lastInsertId int64,
	err error,
) {
	args := []interface{}{
		arg_seq,
		arg_id,
		arg_ord,
		arg_name,
		arg_pw,
	}
	
	sql := fmt.Sprintf(
		"INSERT INTO user VALUES (?, ?, ?, ?, ?)",
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

