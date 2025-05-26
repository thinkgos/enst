package command

import (
	"cmp"
	"fmt"
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/thinkgos/carp/proto"
	"github.com/thinkgos/carp/utils"
)

type protoOpt struct {
	source

	// output directory
	OutputDir string

	// codegen
	PackageName               string            // required, proto 包名
	Options                   map[string]string // required, proto option
	Style                     string            // 字段代码风格, snakeCase, smallCamelCase, pascalCase
	DisableDocComment         bool              // 禁用doc注释
	DisableBool               bool              // 禁用bool,使用int32
	DisableTimestamp          bool              // 禁用google.protobuf.Timestamp,使用int64
	EnableOpenapiv2Annotation bool              // 启用int64的openapiv2注解
}

type protoCmd struct {
	cmd *cobra.Command
	protoOpt
}

func newProtoCmd() *protoCmd {
	root := &protoCmd{}

	cmd := &cobra.Command{
		Use:     "proto",
		Short:   "Generate proto from database",
		Example: "ormat proto",
		RunE: func(*cobra.Command, []string) error {
			sc, err := getSchema(&root.source)
			if err != nil {
				return err
			}
			protoSchemaes := sc.IntoProto()
			packageName := cmp.Or(root.PackageName, utils.GetPkgName(root.OutputDir))
			for _, msg := range protoSchemaes.Entities {
				codegen := &proto.CodeGen{
					Messages:                  []*proto.Message{msg},
					ByName:                    "ormat",
					Version:                   version,
					PackageName:               packageName,
					Options:                   root.Options,
					Style:                     root.Style,
					DisableDocComment:         root.DisableDocComment,
					DisableBool:               root.DisableBool,
					DisableTimestamp:          root.DisableTimestamp,
					EnableOpenapiv2Annotation: root.EnableOpenapiv2Annotation,
				}
				data := codegen.Gen().Bytes()
				filename := joinFilename(root.OutputDir, msg.TableName, ".proto")
				err := WriteFile(filename, data)
				if err != nil {
					return fmt.Errorf("%v: %w", msg.TableName, err)
				}
				slog.Info("👉 " + filename)
			}
			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&root.InputFile, "input", "i", nil, "input file")
	cmd.Flags().StringVarP(&root.Schema, "schema", "s", "file+mysql", "parser file driver, [file+mysql,file+tidb](仅input时有效)")

	// database url
	cmd.Flags().StringVarP(&root.URL, "url", "u", "", "mysql://root:123456@127.0.0.1:3306/test")
	cmd.Flags().StringSliceVarP(&root.Tables, "table", "t", nil, "only out custom table(仅url时有效)")
	cmd.Flags().StringSliceVarP(&root.Exclude, "exclude", "e", nil, "exclude table pattern(仅url时有效)")

	cmd.Flags().StringVarP(&root.OutputDir, "out", "o", "./mapper", "out directory")

	cmd.Flags().StringVar(&root.PackageName, "package", "", "proto package name")
	cmd.Flags().StringToStringVar(&root.Options, "options", nil, "proto options key/value")
	cmd.Flags().StringVar(&root.Style, "style", "", "字段代码风格, [snakeCase,smallCamelCase,pascalCase]")
	cmd.Flags().BoolVar(&root.DisableDocComment, "disableDocComment", false, "禁用文档注释")
	cmd.Flags().BoolVar(&root.DisableBool, "disableBool", false, "禁用bool,使用int32")
	cmd.Flags().BoolVar(&root.DisableTimestamp, "disableTimestamp", false, "禁用google.protobuf.Timestamp,使用int64")
	cmd.Flags().BoolVar(&root.EnableOpenapiv2Annotation, "enableOpenapiv2Annotation", false, "启用用int64的openapiv2注解")

	cmd.MarkFlagsOneRequired("url", "input")

	root.cmd = cmd
	return root
}
