// Code generated by internal/integer.tpl, DO NOT EDIT.

package ens

import (
	"reflect"
)

var _ Fielder = (*uint16Builder)(nil)
var uint16Type = reflect.TypeOf(uint16(0))

func Uint16Type() *GoType {
	return NewGoType(TypeUint16, uint16(0))
}

// Uint16 returns a new Field with type uint16.
func Uint16(name string) *uint16Builder {
	return &uint16Builder{
		&FieldDescriptor{
			Name: name,
			Type: Uint16Type(),
		},
	}
}

// uint16Builder is the builder for uint16 field.
type uint16Builder struct {
	inner *FieldDescriptor
}

// SchemaType sets the column type of the field.
func (b *uint16Builder) SchemaType(ct string) *uint16Builder {
	b.inner.SchemaType = ct
	return b
}

// Comment sets the comment of the field.
func (b *uint16Builder) Comment(c string) *uint16Builder {
	b.inner.Comment = c
	return b
}

// Nullable indicates that this field is a nullable.
func (b *uint16Builder) Nullable() *uint16Builder {
	b.inner.Nullable = true
	return b
}

// Definition set the sql definition of the field.
func (b *uint16Builder) Definition(s string) *uint16Builder {
	b.inner.Definition = s
	return b
}

// GoType overrides the default Go type with a custom one.
//
//	field.Uint16("uint16").
//		GoType(pkg.Uint16(0))
func (b *uint16Builder) GoType(typ any) *uint16Builder {
	b.inner.goType(typ)
	return b
}

// Optional indicates that this field is optional.
// Unlike "Nullable" only fields,
// "Optional" fields are pointers in the generated struct.
func (b *uint16Builder) Optional() *uint16Builder {
	b.inner.Optional = true
	return b
}

// Tags adds a list of tags to the field tag.
//
//	field.Uint16("uint16").
//		Tags("yaml:"xxx"")
func (b *uint16Builder) Tags(tags ...string) *uint16Builder {
	b.inner.Tags = append(b.inner.Tags, tags...)
	return b
}

// Build implements the Fielder interface by returning its descriptor.
func (b *uint16Builder) Build(opt *Option) *FieldDescriptor {
	//	b.inner.checkGoType(uint16Type)
	b.inner.build(opt)
	return b.inner
}
