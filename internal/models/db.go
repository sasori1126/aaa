package models

type FindByField struct {
	Field string
	Value interface{}
}

type QueryByField struct {
	Query string
	Value interface{}
}
