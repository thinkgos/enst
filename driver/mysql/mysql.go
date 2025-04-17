package mysql

import (
	"context"

	"ariga.io/atlas/sql/schema"
	"ariga.io/atlas/sql/sqlclient"

	"github.com/thinkgos/ens"
	"github.com/thinkgos/ens/driver"

	_ "ariga.io/atlas/sql/mysql"
	_ "github.com/go-sql-driver/mysql"
)

var _ driver.Driver = (*MySQL)(nil)

type MySQL struct{}

func (ms *MySQL) InspectSchema(ctx context.Context, arg *driver.InspectOption) (*ens.Schema, error) {
	schemaes, err := ms.inspectSchema(ctx, arg)
	if err != nil {
		return nil, err
	}
	entities := make([]*ens.EntityDescriptor, 0, len(schemaes.Tables))
	for _, tb := range schemaes.Tables {
		entities = append(entities, intoSchema(tb))
	}
	return &ens.Schema{
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
