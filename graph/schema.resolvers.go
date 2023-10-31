package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.40

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"graphql/graph/auth"
	"graphql/graph/model"
)

type auths struct {
	ath *auth.Auth
}

// CreateCompany is the resolver for the createCompany field.
func (r *mutationResolver) CreateCompany(ctx context.Context, input model.NewCompany) (*model.Company, error) {
	company, err := r.CompanyService.CreateCompany(input)
	if err != nil {
		return nil, gqlerror.Errorf("Error creating company")
	}
	return company, nil
}
func (r *mutationResolver) SignIn(ctx context.Context, input model.UserSignIn) (*model.User, error) {
	claims, err := r.UsersService.Authenticate(ctx, input.Email, input.Password)

	if err == ErrInvalidCredentials {
		return nil, gqlerror.Errorf(err.Error())
	}
	if err != nil {
		return nil, gqlerror.Errorf("Error signing in. Please try again")
	}
	err = r.signIn(ctx, user)
	if err != nil {
		return nil, gqlerror.Errorf("Error signing in")
	}
	return user, nil
}

// CreateUser is the resolver for the createUser field.
func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	user, err := r.UsersService.CreateUser(input)
	if err != nil {
		return nil, gqlerror.Errorf("Error creating user")
	}
	err = r.signIn(ctx, user)
	if err != nil {
		return nil, gqlerror.Errorf("Error signing in")
	}
	return user, nil
}
func (r *mutationResolver) signIn(ctx context.Context, user *model.User) error {
	claims, err := r.UsersService.Authenticate(ctx, user.Email, user.PasswordHash)
	if err != nil {
		log.Print("not able to generate clamis")
		return err

	}
	if user.Remember == "" {
		token, err := r.auth.GenerateToken(claims)
		if err != nil {
			return err
		}
		user.Remember = token
		err = r.UsersService.Update(user)
		if err != nil {
			return err
		}
	}
	return nil
}

var ErrInvalidCredentials = errors.New("Invalid username or password")

// SignIn is the resolver for the signIn field.

// Companies is the resolver for the companies field.
func (r *queryResolver) Companies(ctx context.Context) ([]*model.Company, error) {
	panic(fmt.Errorf("not implemented: Companies - companies"))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//   - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//     it when you're done.
//   - You have helper methods in this file. Move them out to keep these resolver files clean.
