package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.27

import (
	"context"
	"fmt"
	"go-web/ent"
	"go-web/ent/user"
)

// UpdatePasswordByAccount is the resolver for the updatePasswordByAccount field.
func (r *mutationResolver) UpdatePasswordByAccount(ctx context.Context, account string, password string) (*ent.User, error) {
	panic(fmt.Errorf("not implemented: UpdatePasswordByAccount - updatePasswordByAccount"))
}

// UserByAccount is the resolver for the userByAccount field.
func (r *queryResolver) UserByAccount(ctx context.Context, account string) (*ent.User, error) {
	return r.client.User.Query().Where(user.Account(account)).First(ctx)
}
