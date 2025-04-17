package mysql

import (
	"fmt"
	"regexp"

	"ariga.io/atlas/sql/schema"
	"github.com/thinkgos/ens"
)

// \b([(]\d+[)])? 匹配0个或1个(\d+)
var typeDictMatchList = []struct {
	Key     string
	NewType func() ens.GoType
}{
	{`^(bool)`, ens.BoolType},                                 // bool
	{`^(tinyint)\b[(]1[)] unsigned`, ens.BoolType},            // bool
	{`^(tinyint)\b[(]1[)]`, ens.BoolType},                     // bool
	{`^(tinyint)\b([(]\d+[)])? unsigned`, ens.Uint8Type},      // uint8
	{`^(tinyint)\b([(]\d+[)])?`, ens.Int8Type},                // int8
	{`^(smallint)\b([(]\d+[)])? unsigned`, ens.Uint16Type},    // uint16
	{`^(smallint)\b([(]\d+[)])?`, ens.Int16Type},              // int16
	{`^(mediumint)\b([(]\d+[)])? unsigned`, ens.Uint32Type},   // uint32
	{`^(mediumint)\b([(]\d+[)])?`, ens.Int32Type},             // int32
	{`^(int)\b([(]\d+[)])? unsigned`, ens.Uint32Type},         // uint32
	{`^(int)\b([(]\d+[)])?`, ens.Int32Type},                   // int32
	{`^(integer)\b([(]\d+[)])? unsigned`, ens.Uint32Type},     // uint32
	{`^(integer)\b([(]\d+[)])?`, ens.Int32Type},               // int32
	{`^(bigint)\b([(]\d+[)])? unsigned`, ens.Uint64Type},      // uint64
	{`^(bigint)\b([(]\d+[)])?`, ens.Int64Type},                // int64
	{`^(float)\b([(]\d+,\d+[)])? unsigned`, ens.Float32Type},  // float32
	{`^(float)\b([(]\d+,\d+[)])?`, ens.Float32Type},           // float32
	{`^(double)\b([(]\d+,\d+[)])? unsigned`, ens.Float64Type}, // float64
	{`^(double)\b([(]\d+,\d+[)])?`, ens.Float64Type},          // float64
	{`^(char)\b[(]\d+[)]`, ens.StringType},                    // string
	{`^(varchar)\b[(]\d+[)]`, ens.StringType},                 // string
	{`^(datetime)\b([(]\d+[)])?`, ens.TimeType},               // time.Time
	{`^(date)\b([(]\d+[)])?`, ens.TimeType},                   // datatypes.Date
	{`^(timestamp)\b([(]\d+[)])?`, ens.TimeType},              // time.Time
	{`^(time)\b([(]\d+[)])?`, ens.TimeType},                   // time.Time
	{`^(year)\b([(]\d+[)])?`, ens.TimeType},                   // time.Time
	{`^(text)\b([(]\d+[)])?`, ens.StringType},                 // string
	{`^(tinytext)\b([(]\d+[)])?`, ens.StringType},             // string
	{`^(mediumtext)\b([(]\d+[)])?`, ens.StringType},           // string
	{`^(longtext)\b([(]\d+[)])?`, ens.StringType},             // string
	{`^(blob)\b([(]\d+[)])?`, ens.BytesType},                  // []byte
	{`^(tinyblob)\b([(]\d+[)])?`, ens.BytesType},              // []byte
	{`^(mediumblob)\b([(]\d+[)])?`, ens.BytesType},            // []byte
	{`^(longblob)\b([(]\d+[)])?`, ens.BytesType},              // []byte
	{`^(bit)\b[(]\d+[)]`, ens.BytesType},                      // []uint8
	{`^(json)\b`, ens.JSONRawMessageType},                     // datatypes.JSON
	{`^(enum)\b[(](.)+[)]`, ens.StringType},                   // string
	{`^(set)\b[(](.)+[)]`, ens.StringType},                    // string
	{`^(decimal)\b[(]\d+,\d+[)]`, ens.DecimalType},            // string
	{`^(binary)\b[(]\d+[)]`, ens.BytesType},                   // []byte
	{`^(varbinary)\b[(]\d+[)]`, ens.BytesType},                // []byte
	{`^(geometry)`, ens.StringType},                           // string
}

func intoGoType(columnType string) ens.GoType {
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

func NewTableDef(tb *schema.Table) ens.TableDef {
	return &TableDef{tb: tb}
}

func (d *TableDef) Table() *schema.Table { return d.tb }

func (d *TableDef) PrimaryKey() ens.IndexDef {
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

func NewColumnDef(col *schema.Column) ens.ColumnDef {
	return &ColumnDef{col: col}
}

func (d *ColumnDef) Column() *schema.Column { return d.col }

func (d *ColumnDef) Definition() string { return intoColumnSql(d.col) }

func (d *ColumnDef) GormTag(tb *schema.Table) string { return intoGormTag(tb, d.col) }

type IndexDef struct {
	index *schema.Index
}

func NewIndexDef(index *schema.Index) ens.IndexDef { return &IndexDef{index: index} }

func (d *IndexDef) Index() *schema.Index { return d.index }

func (d *IndexDef) Definition() string { return intoIndexSql(d.index) }

type ForeignKeyDef struct {
	fk *schema.ForeignKey
}

func NewForeignKey(fk *schema.ForeignKey) ens.ForeignKeyDef {
	return &ForeignKeyDef{fk: fk}
}

func (d *ForeignKeyDef) ForeignKey() *schema.ForeignKey { return d.fk }

func (d *ForeignKeyDef) Definition() string { return intoForeignKeySql(d.fk) }
