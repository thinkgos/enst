package carp

import (
	"github.com/thinkgos/carp/proto"
	"github.com/thinkgos/carp/rapier"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type FieldDescriptor struct {
	ColumnName string // 列名, snake case
	Comment    string // 注释
	Nullable   bool   // Nullable reports whether the column may be null.
	Column     ColumnDef
	// for go
	Type      GoType   // go type information.
	GoName    string   // Go name
	GoPointer bool     // go field is pointer.
	Tags      []string // Tags struct tag
}

func (field *FieldDescriptor) GoType(typ any) {
	field.Type = NewGoType(field.Type.Type, typ)
}

func (field *FieldDescriptor) IntoProto() *proto.MessageField {
	goType := field.Type
	k, n := goType.Type.IntoProtoKind()
	cardinality := protoreflect.Required
	if field.Nullable {
		cardinality = protoreflect.Optional
	}
	return &proto.MessageField{
		Cardinality: cardinality,
		Type:        k,
		TypeName:    n,
		Name:        field.ColumnName,
		ColumnName:  field.ColumnName,
		Comment:     field.Comment,
	}
}

func (field *FieldDescriptor) IntoRapier() *rapier.StructField {
	return &rapier.StructField{
		Type:       field.Type.Type.IntoRapierType(),
		GoName:     field.GoName,
		Nullable:   field.Nullable,
		ColumnName: field.ColumnName,
		Comment:    field.Comment,
	}
}
