package authctx

import (
	"context"

	"github.com/androsyz/nexus-user-svc/constant"
	"github.com/androsyz/nexus-user-svc/handler/middleware"

	"github.com/99designs/gqlgen/graphql"
)

func GetAuthUserID(ctx context.Context) (string, error) {
	authUser, ok := ctx.Value(middleware.UserCtxKey).(*middleware.AuthUser)
	if !ok || authUser == nil || authUser.UserID == "" {
		return "", graphql.ErrorOnPath(ctx, constant.ErrForbiddenAccess)
	}
	return authUser.UserID, nil
}
