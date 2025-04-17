package ens

import "github.com/thinkgos/ens/utils"

type Option struct {
	EnableInt          bool              `yaml:"enableInt" json:"enableInt"`                   // 使能int8,uint8,int16,uint16,int32,uint32输出为int,uint
	EnableBoolInt      bool              `yaml:"enableBoolInt" json:"enableBoolInt"`           // 使能bool输出int
	DisableNullToPoint bool              `yaml:"disableNullToPoint" json:"disableNullToPoint"` // 禁用字段为null时输出指针类型,将输出为sql.Nullxx
	EnableForeignKey   bool              `yaml:"enableForeignKey" json:"enableForeignKey"`     // 输出外键
	IgnoreOmitempty    bool              `yaml:"ignoreOmitempty" json:"ignoreOmitempty"`       // 忽略tags标签的 omitempty 标签
	Tags               map[string]string `yaml:"tags" json:"tags"`                             // tags标签列表, 如 json: snakeCase, support smallCamelCase, pascalCase, snakeCase, kebab
	EscapeName         []string          `yaml:"escapeName" json:"escapeName"`                 // 需要转义的字段
}

func defaultOption() *Option {
	return &Option{
		EnableInt:          false,
		EnableBoolInt:      false,
		DisableNullToPoint: false,
		EnableForeignKey:   false,
		IgnoreOmitempty:    false,
		Tags:               map[string]string{"json": utils.StyleSmallCamelCase},
		EscapeName:         []string{},
	}
}
