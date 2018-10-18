package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type entity struct {
	StructName string   `json:"StructName"`
	TableName  string   `json:"TableName"`
	Columns    []column `json:"Columns"`
	AutoPKey   bool     `json:"AutoPKey"`

	Manys []relationMany `json:"Manys"`
	Ones  []relationOne  `json:"Ones"`
}

func (et *entity) getPkey() column {
	for _, col := range et.Columns {
		if col.PrimaryKey {
			return col
		}
	}
	panic("Error: Don't exist PrimaryKey")
}

type column struct {
	PropertyName string `json:"PropertyName"`
	PropertyType string `json:"PropertyType"`
	ColumnName   string `json:"ColumnName"`
	PrimaryKey   bool   `json:"PrimaryKey"`
}

type relationMany struct {
	RelationClass  string `json:"RelationClass"`
	RelationColumn string `json:"RelationColumn"`
}

type relationOne struct {
	RelationClass    string `json:"RelationClass"`
	RelationProperty string `json:"RelationProperty"`
}

func main() {
	var (
		entitys, dirToGenerate string
	)

	flag.StringVar(&entitys, "entitys", "area,user", "")
	flag.StringVar(&dirToGenerate, "dirToGenerate", "", "")
	flag.Parse()

	enames := strings.Split(entitys, ",")

	println("------------> dirToGenerate =", dirToGenerate)

	for _, ename := range enames {
		println("-------------> ", ename)
		file, e := os.Open(ename + "Model.json")
		if e != nil {
			panic("Cannot open config file")
		}

		et := entity{}

		//Parsing json file
		decoder := json.NewDecoder(file)
		e = decoder.Decode(&et)
		if e != nil {
			panic("Cannot get configuration from file")
		}
		file.Close()

		file, e = os.Create(dirToGenerate + et.StructName + "Model.go")
		if e != nil {
			panic(e.Error())
		}
		_, e = file.WriteString(generateModel(et))
		if e != nil {
			panic(e.Error())
		}
		file.Close()
	}
}

func toLowerCaseFirstLetter(s string) string {
	s = (string)('a'+(s[0]-'A')) + s[1:]
	return s
}

func getUpperAsLower(s string) string {
	r := ""
	for _, c := range s {
		if 'A' <= c && c <= 'Z' {
			r += (string)('a' + (c - 'A'))
		}
	}
	return r
}

func generateModel(et entity) string {
	s := `
package models2

import (
	"github.com/mdiazp/kmodel"
)

///////////////////////////////////////////////////////////////////////////////////
`
	s += fmt.Sprintf(
		`
// %s ...
type %s struct {
`, et.StructName, et.StructName)

	for _, col := range et.Columns {
		s += fmt.Sprintf(
			`	%s %s
`, col.PropertyName, col.PropertyType)
	}

	s +=
		`	model Model
	
`

	for _, many := range et.Manys {
		s += fmt.Sprintf(
			`	%ss *[]*%s
`, toLowerCaseFirstLetter(many.RelationClass), many.RelationClass)
	}

	for _, one := range et.Ones {
		s += fmt.Sprintf(
			`	%s *%s
`, toLowerCaseFirstLetter(one.RelationClass), one.RelationClass)
	}

	pk := et.getPkey()
	lwc := getUpperAsLower(et.StructName)
	s += fmt.Sprintf(`}
/////////////////////////////////////////////////////

// TableName ...
func (%s *%s) TableName() string {
	return "%s"
}
`, lwc, et.StructName, et.TableName)

	s += fmt.Sprintf(`
// AutoPKey ...
func (%s *%s) AutoPKey() bool {
	return %t
}
`, lwc, et.StructName, et.AutoPKey)

	s += fmt.Sprintf(` 
// PkeyName ...
func (%s *%s) PkeyName() string {
	return "%s"
}

// PkeyValue ...
func (%s *%s) PkeyValue() interface{} {
	return %s.%s
}

// PkeyPointer ...
func (%s *%s) PkeyPointer() interface{} {
	return &%s.%s
}
`, lwc, et.StructName,
		pk.ColumnName,
		lwc, et.StructName,
		lwc, pk.PropertyName,
		lwc, et.StructName,
		lwc, pk.PropertyName)

	s += fmt.Sprintf(
		`
// ColumnNames ...
func (%s *%s) ColumnNames() []string {
	return []string{
`, lwc, et.StructName)

	for _, col := range et.Columns {
		if col.PrimaryKey {
			continue
		}
		s += fmt.Sprintf(`		"%s",	
`, col.ColumnName)
	}

	s += fmt.Sprintf(
		`	}
}

// ColumnValues ...
func (%s *%s) ColumnValues() []interface{} {
	return []interface{}{
`, lwc, et.StructName)

	for _, col := range et.Columns {
		if col.PrimaryKey {
			continue
		}
		s += fmt.Sprintf(`			%s.%s,	
`, lwc, col.PropertyName)
	}

	s += fmt.Sprintf(
		`	}
}

// ColumnPointers ...
func (%s *%s) ColumnPointers() []interface{} {
	return []interface{}{
`, lwc, et.StructName)

	for _, col := range et.Columns {
		if col.PrimaryKey {
			continue
		}
		s += fmt.Sprintf(`		&%s.%s,	
`, lwc, col.PropertyName)
	}

	s += fmt.Sprintf(
		`	}
}

/////////////////////////////////////////////////////

// Update ...
func (%s *%s) Update() error {
	return %s.model.Update2(%s)
}
`, lwc, et.StructName, lwc, lwc)

	s += fmt.Sprintf(
		`
// Load ...
func (%s *%s) Load() error {
	return %s.model.Retrieve(%s)
}
`, lwc, et.StructName,
		lwc, lwc)

	for _, one := range et.Ones {
		s += fmt.Sprintf(
			`
// %s ...
func (%s *%s) %s() (*%s, error) {
	var e error
	if %s.%s == nil {
		%s.%s = %s.model.New%s()
		%s.%s.ID = %s.%s
		e = %s.model.Retrieve(%s.%s)
	}
	return %s.%s, e
}
`, one.RelationClass,
			lwc, et.StructName, one.RelationClass, one.RelationClass,
			lwc, toLowerCaseFirstLetter(one.RelationClass),
			lwc, toLowerCaseFirstLetter(one.RelationClass), lwc, one.RelationClass,
			lwc, toLowerCaseFirstLetter(one.RelationClass), lwc, one.RelationProperty,
			lwc, lwc, toLowerCaseFirstLetter(one.RelationClass),
			lwc, toLowerCaseFirstLetter(one.RelationClass),
		)
	}

	for _, manys := range et.Manys {
		s += fmt.Sprintf(
			`	
// %ss ...
func (%s *%s) %ss() (*[]*%s, error) {
	var e error
	if %s.%ss == nil {
		tmp := %s.model.New%sCollection()
		hfilter := fmt.Sprintf("%s=%cd", %s.ID)
		e = %s.model.RetrieveCollection(&hfilter, nil, nil, nil, nil, tmp)
		if e == nil {
			%s.%ss = tmp.%ss
		}
	}
	return %s.%ss, e
}
`, manys.RelationClass,
			lwc, et.StructName, manys.RelationClass, manys.RelationClass,
			lwc, toLowerCaseFirstLetter(manys.RelationClass),
			lwc, manys.RelationClass,
			manys.RelationColumn, '%', lwc,
			lwc,
			lwc, toLowerCaseFirstLetter(manys.RelationClass), manys.RelationClass,
			lwc, toLowerCaseFirstLetter(manys.RelationClass),
		)
	}

	s += fmt.Sprintf(
		`

///////////////////////////////////////////////////////////////////////////////////

// %sCollection ...
type %sCollection struct {
	model Model
	%ss *[]*%s
}

// NewObjectModel ...
func (c *%sCollection) NewObjectModel() kmodel.ObjectModel {
	return c.model.New%s()
}

// Add ...
func (c *%sCollection) Add() kmodel.ObjectModel {
	%s := c.model.New%s()
	*(c.%ss) = append(*(c.%ss), %s)
	return %s
}
`, et.StructName,
		et.StructName,
		et.StructName, et.StructName,
		et.StructName,
		et.StructName,
		et.StructName,
		lwc, et.StructName,
		et.StructName, et.StructName, lwc,
		lwc,
	)

	s += fmt.Sprintf(
		`

///////////////////////////////////////////////////////////////////////////////////

// %sModel ...
type %sModel interface {
	New%s() *%s
	New%sCollection() *%sCollection
	%ss(limit, offset *int, orderby *string,
		orderDesc *bool) (*%sCollection, error)
}
`, et.StructName,
		et.StructName,
		et.StructName, et.StructName,
		et.StructName, et.StructName,
		et.StructName,
		et.StructName,
	)

	s += fmt.Sprintf(
		`
/////////////////////////////////////////////////////

// New%s ...
func (m *model) New%s() *%s {
	%s := &%s{
		model: m,
	}
	return %s
}

// New%sCollection ...
func (m *model) New%sCollection() *%sCollection {
	kk := make([]*%s, 0)
	return &%sCollection{
		model: m,
		%ss: &kk,
	}
}
`, et.StructName,
		et.StructName, et.StructName,
		lwc, et.StructName,
		lwc,
		et.StructName,
		et.StructName, et.StructName,
		et.StructName,
		et.StructName,
		et.StructName,
	)

	s += fmt.Sprintf(
		`

func (m *model) %ss(limit, offset *int, orderby *string,
	orderDesc *bool) (*%sCollection, error) {

	collection := m.New%sCollection()
	e := m.RetrieveCollection(nil, limit, offset, orderby, orderDesc, collection)
	return collection, e
}
`, et.StructName,
		et.StructName,
		et.StructName,
	)

	return s
}
