package codegen

import (
	"strconv"
	"strings"

	"github.com/things-go/ens"
	"github.com/things-go/ens/utils"
)

func (g *CodeGen) GenMapper() *CodeGen {
	if !g.disableDocComment {
		g.P("// Code generated by ", g.byName, ". DO NOT EDIT.")
		g.P("// version: ", g.version)
		g.P()
	}
	g.P(`syntax = "proto3";`)
	g.P()
	g.P("package ", g.packageName, ";")
	g.P()
	if len(g.options) > 0 {
		for k, v := range g.options {
			g.P(`option `, k, ` = "`, v, `";`)
		}
		g.P()
	}

	g.P(`import "protoc-gen-openapiv2/options/annotations.proto";`)
	g.P(`import "protosaber/seaql/seaql.proto";`)
	g.P(`// import "protosaber/enumerate/enumerate.proto";`)
	g.P()

	for _, et := range g.entities {
		structName := utils.CamelCase(et.Name)
		commaOrEmpty := func(r int) string {
			if r == 0 {
				return ""
			}
			return ","
		}

		g.P("// ", structName, " ", trimStructComment(et.Comment, "\n", "\n// "))
		g.P("message ", structName, " {")
		if (et.Table != nil && et.Table.PrimaryKey() != nil) ||
			len(et.Indexes) > 0 {
			g.P("option (things_go.seaql.options) = {")
			g.P("index: [")
			remain := len(et.Indexes)
			if et.Table != nil && et.Table.PrimaryKey() != nil {
				ending := commaOrEmpty(remain)
				g.P("'", et.Table.PrimaryKey().Definition(), "'", ending)
			}
			for _, index := range et.Indexes {
				ending := commaOrEmpty(remain)
				g.P("'", index.Index.Definition(), "'", ending)
			}
			g.P("];")
			g.P("};")
		}
		g.P()
		for i, m := range et.ProtoMessage {
			if m.Comment != "" {
				g.P("// ", m.Comment)
			}
			g.P(genMapperMessageField(i+1, m))
		}
		g.P("}")
	}
	return g
}

func genMapperMessageField(seq int, m *ens.ProtoMessage) string {
	b := strings.Builder{}
	b.Grow(256)
	b.WriteString(m.DataType)
	b.WriteString(" ")
	b.WriteString(m.Name)
	b.WriteString(" = ")
	b.WriteString(strconv.Itoa(seq))
	if len(m.Annotations) > 0 {
		b.WriteString(" [")
		b.WriteString(strings.Join(m.Annotations, ", "))
		b.WriteString("]")
	}
	b.WriteString(";")
	return b.String()
}
