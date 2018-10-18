package kmodel

///////////////////////////////////////////////////////////////////////////////////

// ObjectCollection ...
type ObjectCollection interface {
	Add() ObjectModel
	NewObjectModel() ObjectModel
}

///////////////////////////////////////////////////////////////////////////////////

// ObjectCollection ...
type objectCollection struct {
	model     Model
	List      *[]ObjectModel
	newObject func() ObjectModel
}

func (c *objectCollection) NewObjectModel() ObjectModel {
	return c.newObject()
}

func (c *objectCollection) Add() ObjectModel {
	o := c.newObject()
	*c.List = append(*c.List, o)
	return o
}

func (c *objectCollection) Collection() []ObjectModel {
	return *(c.List)
}

///////////////////////////////////////////////////////////////////////////////////

func (m *model) NewObjectCollection(newObject func() ObjectModel) ObjectCollection {
	list := new([]ObjectModel)
	*list = make([]ObjectModel, 0)
	return &objectCollection{
		model:     m,
		List:      list,
		newObject: newObject,
	}
}
