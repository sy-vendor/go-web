//go:build ignore
// +build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
)

func main() {
	ex, err := entgql.NewExtension(
		entgql.WithRelaySpec(true),
		entgql.WithWhereFilters(true),
		entgql.WithConfigPath("../gqlgen.yml"),
		entgql.WithSchemaPath("../graph/generate/ent.graphql"),
	)
	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}
	opts := []entc.Option{
		entc.Extensions(ex),
		entc.TemplateDir("./template"),
	}
	if err := entc.Generate("./schema", &gen.Config{
		Features: []gen.Feature{gen.FeatureVersionedMigration, gen.FeatureExecQuery, gen.FeatureLock, gen.FeatureModifier, gen.FeatureUpsert},
		IDType:   &field.TypeInfo{Type: field.TypeUint64},
	}, opts...); err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
