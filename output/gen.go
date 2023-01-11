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
	val_seq uint32,
	val_id2 *string,
	val_address []byte,
	val_registered bool,
) (
	lastInsertId int64,
	err error,
) {
	args := []interface{}{
		val_seq,
		val_id2,
		val_address,
		val_registered,
	}
	
	sql := fmt.Sprintf(
		"INSERT INTO tbltest VALUES (?, ?, ?, ?)",
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

type Tbltest_select struct {
	Seq        uint32
	Id2        *string
	Address    []byte
	Registered bool
}

func (t *Tbltest) Select(
	where_seq uint32,
) (
	selects []*Tbltest_select,
	err error,
) {
	args := []interface{}{
		where_seq,
	}
	
	sql := fmt.Sprintf(
		"SELECT * FROM tbltest WHERE seq = ?",
	)
	ret, err := t.job.Query(
		sql,
		args...,
	)
	if err != nil {
		return nil, err
	}
	defer ret.Close()
	
	selects = make([]*Tbltest_select, 0, 100)
	for ret.Next() {
		scan := &Tbltest_select{}
		err := ret.Scan(scan)
		if err != nil {
			return nil, err
		}
		selects = append(selects, scan)
	}
	
	return selects, nil
}

func (t *Tbltest) Update(
	set_seq uint32,
	set_id2 *string,
	set_address []byte,
	set_registered bool,
	where_seq uint32,
) (
	rowAffected int64,
	err error,
) {
	sql := fmt.Sprintf(
		"UPDATE tbltest SET seq = ?, id2 = ?, address = ?, registered = ? WHERE seq = ?",
	)
	args := []interface{}{
		set_seq,
		set_id2,
		set_address,
		set_registered,
		where_seq,
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

func (t *Tbltest) Delete(
	where_seq uint32,
) (
	rowAffected int64,
	err error,
) {
	args := []interface{}{
		where_seq,
	}
	
	sql := fmt.Sprintf(
		"DELETE FROM tbltest WHERE seq = ?",
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
	val_seq uint32,
	val_id2 *string,
	val_address []byte,
	val_registered bool,
) (
	lastInsertId int64,
	err error,
) {
	args := []interface{}{
		val_seq,
		val_id2,
		val_address,
		val_registered,
	}
	
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
	Id2        *string
	Address    []byte
	Registered bool
}

func (t *User) Select(
	where_seq uint32,
) (
	selects []*User_select,
	err error,
) {
	args := []interface{}{
		where_seq,
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
	set_seq uint32,
	set_id2 *string,
	set_address []byte,
	set_registered bool,
	where_seq uint32,
) (
	rowAffected int64,
	err error,
) {
	sql := fmt.Sprintf(
		"UPDATE user SET seq = ?, id2 = ?, address = ?, registered = ? WHERE seq = ?",
	)
	args := []interface{}{
		set_seq,
		set_id2,
		set_address,
		set_registered,
		where_seq,
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
	where_seq uint32,
) (
	rowAffected int64,
	err error,
) {
	args := []interface{}{
		where_seq,
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

