package enst

import (
	"github.com/thinkgos/enst/proto"
	"github.com/thinkgos/enst/rapier"
	"github.com/thinkgos/enst/sqlx"
)

// Schema
type Schema struct {
	Name     string              // schema name
	Entities []*EntityDescriptor // schema entity.
}

func (s *Schema) IntoProto() *proto.Schema {
	entities := make([]*proto.Message, 0, len(s.Entities))
	for _, entity := range s.Entities {
		entities = append(entities, entity.IntoProto())
	}
	return &proto.Schema{
		Name:     s.Name,
		Entities: entities,
	}
}

func (s *Schema) IntoRapier() *rapier.Schema {
	entities := make([]*rapier.Struct, 0, len(s.Entities))
	for _, entity := range s.Entities {
		entities = append(entities, entity.IntoRapier())
	}
	return &rapier.Schema{
		Name:     s.Name,
		Entities: entities,
	}
}

func (s *Schema) IntoSQL() *sqlx.Schema {
	entities := make([]*sqlx.Table, 0, len(s.Entities))
	for _, entity := range s.Entities {
		entities = append(entities, entity.IntoSQL())
	}
	return &sqlx.Schema{
		Name:     s.Name,
		Entities: entities,
	}
}
