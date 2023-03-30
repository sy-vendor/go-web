package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("name").MaxLen(50).Comment("姓名"),
		field.Bool("sex").Comment("性别"),
		field.Int("age").Comment("年龄"),
		field.String("account").MaxLen(20).Comment("账号"),
		field.String("password").MaxLen(20).Comment("密码"),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{}
}
