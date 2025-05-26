package driver

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"ariga.io/atlas/sql/schema"
	"github.com/thinkgos/carp"
)

const (
	Mysql         = "mysql"
	FileMysql     = "file+mysql"
	FileMysqlTidb = "file+tidb"
)

var drivers sync.Map

type Driver interface {
	InspectSchema(context.Context, *InspectOption) (*carp.Schema, error)
}

func RegisterDriver(name string, d Driver) {
	if _, ok := drivers.Load(name); ok {
		panic("driver: Register called twice for " + name)
	}
	drivers.Store(name, d)
}

func LoadDriver(name string) (Driver, error) {
	d, ok := drivers.Load(name)
	if !ok {
		return nil, fmt.Errorf("unsupported schema, only support [%v]", strings.Join(DriverNames(), ", "))
	}
	return d.(Driver), nil
}

func DriverNames() []string {
	var names []string

	drivers.Range(func(key, value any) bool {
		names = append(names, key.(string))
		return true
	})
	return names
}

type InspectOption struct {
	// URL See: https://atlasgo.io/url
	URL string
	// sql data, for file
	Data string
	// InspectOptions describes options for Inspector.
	schema.InspectOptions
}
