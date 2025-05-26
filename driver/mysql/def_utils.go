package mysql

import (
	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"

	"github.com/thinkgos/carp/internal/insql"
)

func autoIncrement(attrs []schema.Attr) bool {
	return insql.Has(attrs, &mysql.AutoIncrement{})
}

func onUpdate(attrs []schema.Attr) (string, bool) {
	var val mysql.OnUpdate
	ok := insql.Has(attrs, &val)
	return val.A, ok
}

func findIndexType(attrs []schema.Attr) string {
	var t mysql.IndexType
	if insql.Has(attrs, &t) && t.T != "" {
		return t.T
	} else {
		return "BTREE"
	}
}
