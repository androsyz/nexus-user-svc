package middleware

import (
	"context"
	"strings"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserCtxKey contextKey = "x-user-id"
)

type AuthUser struct {
	UserID string
}

func ApplyAuthMiddleware(srv *handler.Server, secret string) {
	srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		rc := graphql.GetOperationContext(ctx)
		if rc == nil {
			return next(ctx)
		}

		authHeader := rc.Headers.Get("Authorization")
		if authHeader == "" {
			return next(ctx)
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return next(ctx)
		}

		userID, _ := claims["sub"].(string)
		authUser := &AuthUser{UserID: userID}
		ctx = context.WithValue(ctx, UserCtxKey, authUser)

		return next(ctx)
	})
}
