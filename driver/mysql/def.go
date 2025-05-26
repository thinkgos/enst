package mysql

import (
	"fmt"
	"regexp"

	"ariga.io/atlas/sql/schema"
	"github.com/thinkgos/carp"
)

// \b([(]\d+[)])? 匹配0个或1个(\d+)
var typeDictMatchList = []struct {
	Key     string
	NewType func() carp.GoType
}{
	{`^(bool)`, carp.BoolType},                                 // bool
	{`^(tinyint)\b[(]1[)] unsigned`, carp.BoolType},            // bool
	{`^(tinyint)\b[(]1[)]`, carp.BoolType},                     // bool
	{`^(tinyint)\b([(]\d+[)])? unsigned`, carp.Uint8Type},      // uint8
	{`^(tinyint)\b([(]\d+[)])?`, carp.Int8Type},                // int8
	{`^(smallint)\b([(]\d+[)])? unsigned`, carp.Uint16Type},    // uint16
	{`^(smallint)\b([(]\d+[)])?`, carp.Int16Type},              // int16
	{`^(mediumint)\b([(]\d+[)])? unsigned`, carp.Uint32Type},   // uint32
	{`^(mediumint)\b([(]\d+[)])?`, carp.Int32Type},             // int32
	{`^(int)\b([(]\d+[)])? unsigned`, carp.Uint32Type},         // uint32
	{`^(int)\b([(]\d+[)])?`, carp.Int32Type},                   // int32
	{`^(integer)\b([(]\d+[)])? unsigned`, carp.Uint32Type},     // uint32
	{`^(integer)\b([(]\d+[)])?`, carp.Int32Type},               // int32
	{`^(bigint)\b([(]\d+[)])? unsigned`, carp.Uint64Type},      // uint64
	{`^(bigint)\b([(]\d+[)])?`, carp.Int64Type},                // int64
	{`^(float)\b([(]\d+,\d+[)])? unsigned`, carp.Float32Type},  // float32
	{`^(float)\b([(]\d+,\d+[)])?`, carp.Float32Type},           // float32
	{`^(double)\b([(]\d+,\d+[)])? unsigned`, carp.Float64Type}, // float64
	{`^(double)\b([(]\d+,\d+[)])?`, carp.Float64Type},          // float64
	{`^(char)\b[(]\d+[)]`, carp.StringType},                    // string
	{`^(varchar)\b[(]\d+[)]`, carp.StringType},                 // string
	{`^(datetime)\b([(]\d+[)])?`, carp.TimeType},               // time.Time
	{`^(date)\b([(]\d+[)])?`, carp.TimeType},                   // datatypes.Date
	{`^(timestamp)\b([(]\d+[)])?`, carp.TimeType},              // time.Time
	{`^(time)\b([(]\d+[)])?`, carp.TimeType},                   // time.Time
	{`^(year)\b([(]\d+[)])?`, carp.TimeType},                   // time.Time
	{`^(text)\b([(]\d+[)])?`, carp.StringType},                 // string
	{`^(tinytext)\b([(]\d+[)])?`, carp.StringType},             // string
	{`^(mediumtext)\b([(]\d+[)])?`, carp.StringType},           // string
	{`^(longtext)\b([(]\d+[)])?`, carp.StringType},             // string
	{`^(blob)\b([(]\d+[)])?`, carp.BytesType},                  // []byte
	{`^(tinyblob)\b([(]\d+[)])?`, carp.BytesType},              // []byte
	{`^(mediumblob)\b([(]\d+[)])?`, carp.BytesType},            // []byte
	{`^(longblob)\b([(]\d+[)])?`, carp.BytesType},              // []byte
	{`^(bit)\b[(]\d+[)]`, carp.BytesType},                      // []uint8
	{`^(json)\b`, carp.JSONRawMessageType},                     // datatypes.JSON
	{`^(enum)\b[(](.)+[)]`, carp.StringType},                   // string
	{`^(set)\b[(](.)+[)]`, carp.StringType},                    // string
	{`^(decimal)\b[(]\d+,\d+[)]`, carp.DecimalType},            // string
	{`^(binary)\b[(]\d+[)]`, carp.BytesType},                   // []byte
	{`^(varbinary)\b[(]\d+[)]`, carp.BytesType},                // []byte
	{`^(geometry)`, carp.StringType},                           // string
}

func intoGoType(columnType string) carp.GoType {
	for _, v := range typeDictMatchList {
		ok, _ := regexp.MatchString(v.Key, columnType)
		if ok {
			return v.NewType()
		}
	}
	panic(fmt.Sprintf("type (%v) not match in any way, need to add on (https://github.com/thinkgos/ormat/blob/main/driver/mysql/def.go)", columnType))
}

type TableDef struct {
	tb *schema.Table
}

func NewTableDef(tb *schema.Table) carp.TableDef {
	return &TableDef{tb: tb}
}

func (d *TableDef) Table() *schema.Table { return d.tb }

func (d *TableDef) PrimaryKey() carp.IndexDef {
	if d.tb.PrimaryKey != nil {
		return NewIndexDef(d.tb.PrimaryKey)
	}
	return nil
}

func (d *TableDef) Definition() string {
	return intoTableSql(d.tb)
}

type ColumnDef struct {
	col *schema.Column
}

func NewColumnDef(col *schema.Column) carp.ColumnDef {
	return &ColumnDef{col: col}
}

func (d *ColumnDef) Column() *schema.Column { return d.col }

func (d *ColumnDef) Definition() string { return intoColumnSql(d.col) }

func (d *ColumnDef) GormTag(tb *schema.Table) string { return intoGormTag(tb, d.col) }

type IndexDef struct {
	index *schema.Index
}

func NewIndexDef(index *schema.Index) carp.IndexDef { return &IndexDef{index: index} }

func (d *IndexDef) Index() *schema.Index { return d.index }

func (d *IndexDef) Definition() string { return intoIndexSql(d.index) }

type ForeignKeyDef struct {
	fk *schema.ForeignKey
}

func NewForeignKey(fk *schema.ForeignKey) carp.ForeignKeyDef {
	return &ForeignKeyDef{fk: fk}
}

func (d *ForeignKeyDef) ForeignKey() *schema.ForeignKey { return d.fk }

func (d *ForeignKeyDef) Definition() string { return intoForeignKeySql(d.fk) }
