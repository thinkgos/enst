package mysql

import (
	"context"

	"ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/sqlclient"

	"github.com/thinkgos/carp"
	"github.com/thinkgos/carp/driver"

	_ "ariga.io/atlas/sql/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var _ driver.Driver = (*MySQL)(nil)

type MySQL struct{}

func (ms *MySQL) InspectSchema(ctx context.Context, arg *driver.InspectOption) (*carp.Schema, error) {
	schemaes, err := ms.inspectSchema(ctx, arg)
	if err != nil {
		return nil, err
	}
	entities := make([]*carp.EntityDescriptor, 0, len(schemaes.Tables))
	for _, tb := range schemaes.Tables {
		entities = append(entities, intoSchema(tb))
	}
	return &carp.Schema{
		Name:     schemaes.Name,
		Entities: entities,
	}, nil
}

func (ms *MySQL) inspectSchema(ctx context.Context, arg *driver.InspectOption) (*schema.Schema, error) {
	client, err := sqlclient.Open(ctx, arg.URL)
	if err != nil {
		return nil, err
	}
	return client.InspectSchema(ctx, "", &arg.InspectOptions)
}
