package resolver

import (
	"context"

	"github.com/androsyz/nexus-user-svc/graph/model"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type ucUserInterface interface {
	UpdateUser(ctx context.Context, input model.UpdateUserInput) (*model.User, error)
	DeleteUser(ctx context.Context, id string) (bool, error)
	User(ctx context.Context, id string) (*model.User, error)
}

type ucRoleInterface interface {
	CreateRole(ctx context.Context, input model.CreateRoleInput) (*model.Role, error)
	UpdateRole(ctx context.Context, input model.UpdateRoleInput) (*model.Role, error)
	DeleteRole(ctx context.Context, roleID string) (bool, error)
	AssignRole(ctx context.Context, input model.AssignRoleInput) (bool, error)
	RevokeRole(ctx context.Context, input model.AssignRoleInput) (bool, error)
	Roles(ctx context.Context) ([]*model.Role, error)
}

type ucAuthInterface interface {
	Register(ctx context.Context, request model.RegisterRequest) (*model.AuthResponse, error)
	Login(ctx context.Context, request model.LoginRequest) (*model.AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*model.AuthResponse, error)
	VerifyToken(ctx context.Context, token string) (bool, error)
	RevokeToken(ctx context.Context) (bool, error)
}

func NewResolver(
	ucUser ucUserInterface,
	ucAuth ucAuthInterface,
	ucRole ucRoleInterface,
) (*Resolver, error) {
	return &Resolver{
		ucUser: ucUser,
		ucAuth: ucAuth,
		ucRole: ucRole,
	}, nil
}

type Resolver struct {
	ucUser ucUserInterface
	ucAuth ucAuthInterface
	ucRole ucRoleInterface
}
