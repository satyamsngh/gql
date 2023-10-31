package graph

import (
	"gorm.io/gorm"
	"graphql/graph/auth"
	"graphql/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DB             *gorm.DB
	UsersService   *model.UserService
	CompanyService *model.CompanyService
	auth           auth.Auth
}
