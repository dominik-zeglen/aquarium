package api

import (
	"log"

	graphql "github.com/graph-gophers/graphql-go"
	"github.com/graph-gophers/graphql-go/relay"

	"github.com/dominik-zeglen/aquarium/api/schema"
	"github.com/dominik-zeglen/aquarium/sim"
)

func GetSchemaStr() (*string, error) {
	schemaData, err := schema.Asset("api/schema/schema.graphql")
	if err != nil {
		return nil, err
	}

	schemaStr := string(schemaData)

	return &schemaStr, nil
}

func GetSchema(sim *sim.Sim) (*graphql.Schema, error) {
	schemaStr, err := GetSchemaStr()
	if err != nil {
		return nil, err
	}

	opts := []graphql.SchemaOpt{graphql.UseFieldResolvers()}
	schema := graphql.MustParseSchema(*schemaStr, &Query{sim}, opts...)

	return schema, nil
}

func InitAPI(sim *sim.Sim) *relay.Handler {
	schema, err := GetSchema(sim)
	if err != nil {
		log.Fatal(err)
	}

	return &relay.Handler{Schema: schema}
}
