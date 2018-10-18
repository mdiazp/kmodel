package kmodel

import (
	"database/sql"
	"fmt"
	"strconv"
)

///////////////////////////////////////////////////////////////////////////////////

// CRUD ...
type CRUD interface {
	Create(o ObjectModel) error
	Retrieve(o ObjectModel) error
	Update2(o ObjectModel) error
	Delete2(o ObjectModel) error
	RetrieveCollection(hfilter *string, limit *int, offset *int,
		orderby *string, ordDesc *bool, c ObjectCollection) error
}

///////////////////////////////////////////////////////////////////////////////////

// Create ...
func (m *model) Create(o ObjectModel) error {
	cnames := columnNames(o, !o.AutoPKey())
	qnums := ""
	ln := len(cnames)

	q := "INSERT INTO " + o.TableName() + " ("
	paramID := 1

	for i := 0; i < ln; i++ {
		q += cnames[i]
		qnums += "$" + strconv.Itoa(paramID)
		paramID++

		if i+1 < ln {
			q += ", "
			qnums += ", "
		}
	}

	q += ") VALUES (" + qnums + ")"
	q += " RETURNING " + fComa(columnNames(o, true)...)

	stmt, e := m.DB.Prepare(q)
	if e == nil {
		defer stmt.Close()
		e = stmt.QueryRow(columnValues(o, !o.AutoPKey())...).
			Scan(columnPointers(o, true)...)
	}
	return e
}

///////////////////////////////////////////////////////////////////////////////////

// Retrieve ...
func (m *model) Retrieve(o ObjectModel) error {
	q := "SELECT " + fComa(columnNames(o, true)...)
	q += " FROM " + o.TableName() + " WHERE " + o.PkeyName() + "=$1"

	println("---------> KMODEL --> ", q)

	stmt, e := m.DB.Prepare(q)
	if e == nil {
		defer stmt.Close()
		e = stmt.QueryRow(o.PkeyValue()).
			Scan(columnPointers(o, true)...)
	}
	return e
}

///////////////////////////////////////////////////////////////////////////////////

// Update ...
func (m *model) Update2(o ObjectModel) error {
	q := "UPDATE " + o.TableName() + " SET "

	cnames := columnNames(o)
	ln := len(cnames)
	for i := 0; i < ln; i++ {
		q += cnames[i] + "=$" + strconv.Itoa(i+1)
		if i+1 < ln {
			q += ", "
		}
	}
	q += " WHERE " + o.PkeyName() + "=$" + strconv.Itoa(ln+1)
	q += " RETURNING " + fComa(columnNames(o, true)...)

	stmt, e := m.DB.Prepare(q)
	if e == nil {
		defer stmt.Close()
		e = stmt.QueryRow(columnValues(o, true)...).
			Scan(columnPointers(o, true)...)
	}
	return e
}

///////////////////////////////////////////////////////////////////////////////////

// Delete ...
func (m *model) Delete2(o ObjectModel) error {
	q := "DELETE FROM " + o.TableName() +
		" WHERE " + o.PkeyName() + "=$1"

	stmt, e := m.DB.Prepare(q)
	if e == nil {
		defer stmt.Close()
		_, e = stmt.Exec(o.PkeyValue())
	}
	return e
}

///////////////////////////////////////////////////////////////////////////////////

// RetrieveList ...

func (m *model) RetrieveCollection(
	hfilter *string, limit *int, offset *int,
	orderby *string, ordDesc *bool, c ObjectCollection) error {

	o := c.NewObjectModel()

	q := "SELECT " + fComa(columnNames(o, true)...)

	q += " FROM " + o.TableName()

	if hfilter != nil {
		q += " WHERE " + *hfilter
	}

	if limit != nil {
		q += " LIMIT " + strconv.Itoa(*limit)
	}

	if offset != nil {
		q += " OFFSET " + strconv.Itoa(*offset)
	}

	if orderby != nil {
		q += " ORDER BY " + *orderby
		if ordDesc != nil && *ordDesc {
			q += " DESC"
		} else {
			q += " ASC"
		}
	}

	fmt.Println(q)

	stmt, e := m.DB.Prepare(q)
	if e == nil {
		defer stmt.Close()
		var rows *sql.Rows
		rows, e = stmt.Query()

		for rows.Next() {
			no := c.Add()
			e = rows.Scan(columnPointers(no, true)...)
		}
	}
	return e
}
