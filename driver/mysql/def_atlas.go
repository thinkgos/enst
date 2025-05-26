package mysql

import (
	"fmt"
	"strings"

	"ariga.io/atlas/sql/mysql"
	"ariga.io/atlas/sql/schema"

	"github.com/thinkgos/carp"
	"github.com/thinkgos/carp/internal/insql"
	"github.com/thinkgos/carp/utils"
)

func intoColumnSql(col *schema.Column) string {
	nullable := col.Type.Null
	autoIncrement := autoIncrement(col.Attrs)

	b := &strings.Builder{}
	b.Grow(64)
	b.WriteString(col.Type.Raw)
	if !nullable {
		b.WriteString(" NOT NULL")
	}
	if autoIncrement {
		b.WriteString(" AUTO_INCREMENT")
	} else {
		switch x := schema.UnderlyingExpr(col.Default).(type) {
		case *schema.Literal:
			fmt.Fprintf(b, " DEFAULT '%s'", strings.Trim(x.V, `"`))
		case *schema.RawExpr:
			fmt.Fprintf(b, " DEFAULT %s", x.X)
		case nil:
			if col.Type.Null {
				b.WriteString(" DEFAULT NULL")
			}
		default:
			// do nothing
		}
	}
	if v, ok := onUpdate(col.Attrs); ok {
		fmt.Fprintf(b, " ON UPDATE %s", v)
	}
	return b.String()
}

func intoIndexSql(index *schema.Index) string {
	fields := insql.IndexPartColumnNames(index.Parts)
	indexType := findIndexType(index.Attrs)
	fieldList := "`" + strings.Join(fields, "`,`") + "`"
	if insql.IndexEqual(index.Table.PrimaryKey, index) {
		return fmt.Sprintf("PRIMARY KEY (%s) USING %s", fieldList, indexType)
	} else if index.Unique {
		return fmt.Sprintf("UNIQUE KEY `%s` (%s) USING %s", index.Name, fieldList, indexType)
	} else {
		return fmt.Sprintf("KEY `%s` (%s) USING %s", index.Name, fieldList, indexType)
	}
}

func intoForeignKeySql(fk *schema.ForeignKey) string {
	columnNameList := "`" + strings.Join(insql.ColumnNames(fk.Columns), "`,`") + "`"
	refColumnNameList := "`" + strings.Join(insql.ColumnNames(fk.RefColumns), "`,`") + "`"
	return fmt.Sprintf(
		"CONSTRAINT `%s` FOREIGN KEY (%s) REFERENCES `%s` (%s) ON DELETE %s ON UPDATE %s",
		fk.Symbol, columnNameList, fk.RefTable.Name, refColumnNameList, fk.OnDelete, fk.OnUpdate,
	)
}

func intoTableSql(tb *schema.Table) string {
	b := &strings.Builder{}
	b.Grow(256)
	fmt.Fprintf(b, "CREATE TABLE `%s` (\n", tb.Name)

	remain := len(tb.Columns) + len(tb.Indexes) + len(tb.ForeignKeys)
	if tb.PrimaryKey != nil {
		remain++
	}
	suffixOrEmpty := func(r int) string {
		if r == 0 {
			return ""
		}
		return ","
	}
	//* columns
	for _, col := range tb.Columns {
		remain--
		suffix := suffixOrEmpty(remain)
		comment, ok := insql.Comment(col.Attrs)
		if ok {
			comment = fmt.Sprintf(" COMMENT '%s'", comment)
		}
		fmt.Fprintf(b, "  `%s` %s%s%s\n", col.Name, intoColumnSql(col), comment, suffix)
	}
	//* pk + indexes
	if tb.PrimaryKey != nil {
		remain--
		suffix := suffixOrEmpty(remain)
		fmt.Fprintf(b, "  %s%s\n", intoIndexSql(tb.PrimaryKey), suffix)
	}
	for _, index := range tb.Indexes {
		remain--
		if insql.IndexEqual(tb.PrimaryKey, index) { // ignore primary key, maybe include
			continue
		}
		suffix := suffixOrEmpty(remain)
		fmt.Fprintf(b, "  %s%s\n", intoIndexSql(index), suffix)
	}
	//* foreignKeys
	for _, fk := range tb.ForeignKeys {
		remain--
		suffix := suffixOrEmpty(remain)
		fmt.Fprintf(b, "  %s%s\n", intoForeignKeySql(fk), suffix)
	}

	engine := mysql.EngineInnoDB
	charset := "utf8mb4"
	collate := ""
	comment := ""
	for _, attr := range tb.Attrs {
		switch val := attr.(type) {
		case *mysql.Engine:
			engine = val.V
		case *schema.Charset:
			charset = val.V
		case *schema.Collation:
			collate = val.V
		case *schema.Comment:
			comment = val.Text
			// case *mysql.AutoIncrement: // ignore this
		}
	}
	fmt.Fprintf(b, ") ENGINE=%s DEFAULT CHARSET=%s", engine, charset)
	if collate != "" {
		fmt.Fprintf(b, " COLLATE=%s", collate)
	}
	if comment != "" {
		fmt.Fprintf(b, " COMMENT='%s'", comment)
	}
	return b.String()
}

// column, type, not null, authIncrement, default, [primaryKey|index], comment
func intoGormTag(tb *schema.Table, col *schema.Column) string {
	pkPriority, isPk := 0, false
	if pk := tb.PrimaryKey; pk != nil {
		pkPriority, isPk = insql.FindIndexPartSeq(pk.Parts, col)
	}
	autoIncrement := autoIncrement(col.Attrs)

	b := &strings.Builder{}
	b.Grow(64)
	fmt.Fprintf(b, `gorm:"column:%s`, col.Name)
	if !isPk || !autoIncrement {
		fmt.Fprintf(b, ";type:%s", col.Type.Raw)
	}
	if !col.Type.Null {
		fmt.Fprintf(b, ";not null")
	}

	if isPk {
		if autoIncrement {
			fmt.Fprintf(b, ";autoIncrement:true")
		}
	} else {
		dv := ""
		switch x := schema.UnderlyingExpr(col.Default).(type) {
		case *schema.Literal:
			dv = x.V
			if dv == `""` || dv == "" {
				dv = "''"
			} else {
				dv = strings.Trim(dv, `"`) // format: `"xxx"` or `'xxx'`
			}
		case *schema.RawExpr:
			dv = x.X
		case nil:
			if col.Type.Null {
				dv = "null"
			}
		default:
			// do nothing
		}
		// FIXME: bug, gorm目前 on update放在default上. 所以必须default存在.
		if dv != "" {
			if v, ok := onUpdate(col.Attrs); ok {
				dv = fmt.Sprintf("%s ON UPDATE %s", dv, v)
			}
			fmt.Fprintf(b, ";default:%s", dv)
		}
	}

	//* pk + indexes
	if isPk && tb.PrimaryKey != nil {
		fmt.Fprintf(b, ";primaryKey")
		if len(tb.PrimaryKey.Parts) > 1 {
			fmt.Fprintf(b, ",priority:%d", pkPriority)
		}
	}
	for _, val := range col.Indexes {
		if insql.IndexEqual(tb.PrimaryKey, val) { // ignore primary key, may be include
			continue
		}
		if val.Unique {
			fmt.Fprintf(b, ";uniqueIndex:%s", val.Name)
		} else {
			fmt.Fprintf(b, ";index:%s", val.Name)
			// 	mysql.IndexTypeFullText
			// if v.IndexType == "FULLTEXT" {
			// 	b.WriteString(",class:FULLTEXT")
			// }
		}
		if len(val.Parts) > 1 {
			priority, ok := insql.FindIndexPartSeq(val.Parts, col)
			if ok {
				fmt.Fprintf(b, ",priority:%d", priority)
			}
		}
	}
	if comment, ok := insql.Comment(col.Attrs); ok && comment != "" {
		fmt.Fprintf(b, ";comment:%s", utils.TrimFieldComment(comment))
	}
	b.WriteString(`"`)
	return b.String()
}

func intoSchema(tb *schema.Table) *carp.EntityDescriptor {
	// * columns
	fielders := make([]*carp.FieldDescriptor, 0, len(tb.Columns))
	for _, col := range tb.Columns {
		fielders = append(fielders, &carp.FieldDescriptor{
			ColumnName: col.Name,
			Comment:    insql.MustComment(col.Attrs),
			Nullable:   col.Type.Null,
			Column:     NewColumnDef(col),
			Type:       intoGoType(col.Type.Raw),
			GoName:     utils.PascalCase(col.Name),
			GoPointer:  col.Type.Null,
			Tags:       []string{intoGormTag(tb, col)},
		})
	}
	// * indexes
	indexers := make([]*carp.IndexDescriptor, 0, len(tb.Indexes))
	for _, index := range tb.Indexes {
		indexers = append(indexers, &carp.IndexDescriptor{
			Name:   index.Name,
			Fields: insql.IndexPartColumnNames(index.Parts),
			Index:  NewIndexDef(index),
		})
	}
	//* foreignKeys
	fks := make([]*carp.ForeignKeyDescriptor, 0, len(tb.ForeignKeys))
	for _, fk := range tb.ForeignKeys {
		fks = append(fks, &carp.ForeignKeyDescriptor{
			Symbol:     fk.Symbol,
			Table:      fk.Table.Name,
			Columns:    insql.ColumnNames(fk.Columns),
			RefTable:   fk.RefTable.Name,
			RefColumns: insql.ColumnNames(fk.RefColumns),
			OnUpdate:   fk.OnUpdate,
			OnDelete:   fk.OnDelete,
			ForeignKey: NewForeignKey(fk),
		})
	}

	// * table
	return &carp.EntityDescriptor{
		Name:        tb.Name,
		Comment:     insql.MustComment(tb.Attrs),
		Table:       NewTableDef(tb),
		Fields:      fielders,
		Indexes:     indexers,
		ForeignKeys: fks,
	}
}
