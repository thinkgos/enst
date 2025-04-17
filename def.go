package enst

import "ariga.io/atlas/sql/schema"

type TableDef interface {
	Table() *schema.Table
	Definition() string
	PrimaryKey() IndexDef
}

type ColumnDef interface {
	Column() *schema.Column
	Definition() string
	GormTag(*schema.Table) string
}

type IndexDef interface {
	Index() *schema.Index
	Definition() string
}

type ForeignKeyDef interface {
	ForeignKey() *schema.ForeignKey
	Definition() string
}
