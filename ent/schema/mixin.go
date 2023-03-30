package schema

import (
	"go-web/pkg/uid"
	"go-web/pkg/util"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
)

// Mixin holds the schema definition for the Mixin entity.
type Mixin struct {
	mixin.Schema
}

// Fields of the Mixin.
func (Mixin) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").
			DefaultFunc(uid.NewUniqueIDGenerator().NewID).
			Validate(util.ValidateID),
	}
}

// Edges of the Mixin.
func (Mixin) Edges() []ent.Edge {
	return nil
}
