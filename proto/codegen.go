package proto

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/things-go/ens/utils"
	"golang.org/x/tools/imports"
	"google.golang.org/protobuf/reflect/protoreflect"
)

const googleProtobufTimestamp = "google.protobuf.Timestamp"

type CodeGen struct {
	buf                       bytes.Buffer
	Messages                  []*Message        // required, proto Message
	ByName                    string            // required, 生成名称
	Version                   string            // required, 生成版本
	PackageName               string            // required, proto 包名
	Options                   map[string]string // required, proto option
	Style                     string            // 字段代码风格, snakeCase, smallCamelCase, camelCase
	DisableDocComment         bool              // 禁用doc注释
	DisableBool               bool              // 禁用bool,使用int32
	DisableTimestamp          bool              // 禁用google.protobuf.Timestamp,使用int64
	EnableOpenapiv2Annotation bool              // 启用int64的openapiv2注解
}

// Bytes returns the CodeBuf's buffer.
func (g *CodeGen) Bytes() []byte {
	return g.buf.Bytes()
}

// FormatSource return formats and adjusts imports contents of the CodeGen's buffer.
func (g *CodeGen) FormatSource() ([]byte, error) {
	data := g.buf.Bytes()
	if len(data) == 0 {
		return data, nil
	}
	// return format.Source(data)
	return imports.Process("", data, nil)
}

// Write appends the contents of p to the buffer,
func (g *CodeGen) Write(b []byte) (n int, err error) {
	return g.buf.Write(b)
}

// Print formats using the default formats for its operands and writes to the generated output.
// Spaces are added between operands when neither is a string.
// It returns the number of bytes written and any write error encountered.
func (g *CodeGen) Print(a ...any) (n int, err error) {
	return fmt.Fprint(&g.buf, a...)
}

// Printf formats according to a format specifier for its operands and writes to the generated output.
// It returns the number of bytes written and any write error encountered.
func (g *CodeGen) Printf(format string, a ...any) (n int, err error) {
	return fmt.Fprintf(&g.buf, format, a...)
}

// Fprintln formats using the default formats to the generated output.
// Spaces are always added between operands and a newline is appended.
// It returns the number of bytes written and any write error encountered.
func (g *CodeGen) Println(a ...any) (n int, err error) {
	return fmt.Fprintln(&g.buf, a...)
}

func (g *CodeGen) Gen() *CodeGen {
	if !g.DisableDocComment {
		g.Printf("// Code generated by %s. DO NOT EDIT.\n", g.ByName)
		g.Printf("// version: %s\n", g.Version)
		g.Println()
	}
	g.Println(`syntax = "proto3";`)
	g.Println()
	g.Printf("package %s;\n", g.PackageName)
	g.Println()

	if len(g.Options) > 0 {
		for k, v := range g.Options {
			g.Printf("option %s = \"%s\";\n", k, v)
		}
		g.Println()
	}

	if g.needGoogleProtobufTimestamp(g.Messages) {
		g.Println(`import "google/protobuf/timestamp.proto";`)
	}
	if g.needOpenapiv2Annotation(g.Messages) {
		g.Println(`import "protoc-gen-openapiv2/options/annotations.proto";`)
	}
	g.Println()

	for _, et := range g.Messages {
		structName := utils.CamelCase(et.Name)

		g.Printf("// %s %s\n", structName, strings.ReplaceAll(strings.TrimSpace(et.Comment), "\n", "\n// "))
		g.Printf("message %s {\n", structName)
		for i, m := range et.Fields {
			if m.Comment != "" {
				g.Printf("  // %s\n", m.Comment)
			}
			typeName, annotations := g.intoTypeNameAndAnnotation(m)
			fieldName := utils.StyleName(g.Style, m.Name)
			annotation := ""
			if len(annotations) > 0 {
				annotation = fmt.Sprintf(" [%s]", strings.Join(annotations, ", "))
			}
			seq := i + 1
			if m.Cardinality == protoreflect.Required {
				g.Printf("  %s %s = %d%s;\n", typeName, fieldName, seq, annotation)
			} else {
				g.Printf("  %s %s %s = %d%s;\n", m.Cardinality.String(), typeName, fieldName, seq, annotation)
			}
		}
		g.Println("}")
	}
	return g
}

func (g *CodeGen) intoTypeNameAndAnnotation(field *MessageField) (string, []string) {
	annotations := make([]string, 0, 8)
	switch {
	case g.DisableBool && field.Type == protoreflect.BoolKind:
		return protoreflect.Int32Kind.String(), annotations
	case field.Type == protoreflect.MessageKind && field.TypeName == googleProtobufTimestamp:
		if g.DisableTimestamp {
			if g.EnableOpenapiv2Annotation {
				annotations = append(annotations, `(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { type: [ INTEGER ] }`)
			}
			return protoreflect.Int64Kind.String(), annotations
		} else {
			return field.TypeName, annotations
		}
	case (field.Type == protoreflect.Int64Kind || field.Type == protoreflect.Uint64Kind) && g.EnableOpenapiv2Annotation:
		annotations = append(annotations, `(grpc.gateway.protoc_gen_openapiv2.options.openapiv2_field) = { type: [ INTEGER ] }`)
		fallthrough
	default:
		return field.Type.String(), annotations
	}
}

func (g *CodeGen) needOpenapiv2Annotation(messages []*Message) bool {
	if !g.EnableOpenapiv2Annotation {
		return false
	}
	for _, msg := range messages {
		for _, v := range msg.Fields {
			if v.Type == protoreflect.Int64Kind ||
				g.DisableTimestamp && v.Type == protoreflect.MessageKind && v.TypeName == googleProtobufTimestamp {
				return true
			}
		}
	}
	return false
}

func (g *CodeGen) needGoogleProtobufTimestamp(messages []*Message) bool {
	for _, msg := range messages {
		for _, v := range msg.Fields {
			if !g.DisableTimestamp && v.Type == protoreflect.MessageKind && v.TypeName == googleProtobufTimestamp {
				return true
			}
		}
	}
	return false
}
