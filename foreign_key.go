package carp

import (
	"ariga.io/atlas/sql/schema"
)

type ForeignKeyDescriptor struct {
	Symbol     string
	Table      string
	Columns    []string
	RefTable   string
	RefColumns []string
	OnUpdate   schema.ReferenceOption
	OnDelete   schema.ReferenceOption
	ForeignKey ForeignKeyDef
}
