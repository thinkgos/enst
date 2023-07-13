// Code generated by internal/float.tmpl, DO NOT EDIT.

package ens

import (
	"reflect"
)

{{ $t := $.Kind }}

var _ Fielder = (*{{ $t }}Builder)(nil)
var {{ $t }}Type = reflect.TypeOf({{ $t }}(0))

{{ $title := title $t.String }}

func {{ $title }}Type() *GoType {
	return NewGoType(Type{{ $title }}, {{ $t }}(0))
}

// {{ $title }} returns a new Field with type {{ $t }}.
func {{ $title }}(name string) *{{ $t }}Builder { 
    return &{{ $t }}Builder{
        &FieldDescriptor{
			Name: name,
			SchemaType: "",
		    Type: {{ $title }}Type(),
	    },
    }
}

{{ $builder := printf "%sBuilder" $t }}

// {{ $builder }} is the builder for {{ $t }} fields.
type {{ $builder }} struct {
	inner *FieldDescriptor
}

// SchemaType overrides the default database type with a custom
// schema type (per dialect) for {{ $t.String }}.
//
//	field.{{ title $t.String }}("amount").
//		SchemaType("decimal(5, 2)")
func (b *{{ $builder }}) SchemaType(ct string) *{{ $builder }} {
//	b.inner.SchemaType = ct
	return b
}

// Comment sets the comment of the field.
func (b *{{ $builder }}) Comment(c string) *{{ $builder }} {
	b.inner.Comment = c
	return b
}

// Nullable indicates that this field is a nullable.
func (b *{{ $builder }}) Nullable() *{{ $builder }} {
	b.inner.Nullable = true
	return b
}

// Definition set the sql definition of the field.
func (b *{{ $builder }}) Definition(s string) *{{ $builder }} {
	b.inner.Definition = s
	return b
}

{{ $tt := title $t.String }}
// GoType overrides the default Go type with a custom one.
//
//	field.{{ $tt }}("{{ $t }}").
//		GoType(pkg.{{ $tt }}(0))
//
func (b *{{ $builder }}) GoType(typ any) *{{ $builder }} {
	b.inner.goType(typ)
	return b
}

// Optional indicates that this field is optional.
// Unlike "Nullable" only fields,
// "Optional" fields are pointers in the generated struct.
func (b *{{ $builder }}) Optional() *{{ $builder }} {
	b.inner.Optional = true
	return b
}

// Tags adds a list of tags to the field tag.
//
//	field.{{ $tt }}("{{ $t }}").
//		Tags("yaml:"xxx"")
func (b *{{ $builder }}) Tags(tags ...string) *{{ $builder }} {
	b.inner.Tags = append(b.inner.Tags, tags...)
	return b
}

// Build implements the Fielder interface by returning its descriptor.
func (b *{{ $builder }}) Build(opt *Option) *FieldDescriptor {
//	b.inner.checkGoType({{ $t }}Type)
	b.inner.build(opt)
	return b.inner
}
