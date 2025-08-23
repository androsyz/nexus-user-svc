package main

import (
	"context"
	"log"
	"net/http"
	"slices"
	"time"

	"github.com/androsyz/nexus-user-svc/cmd/initialize"
	"github.com/androsyz/nexus-user-svc/config"
	"github.com/androsyz/nexus-user-svc/graph/generated"
	"github.com/androsyz/nexus-user-svc/handler/middleware"
	"github.com/androsyz/nexus-user-svc/handler/resolver"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gorilla/websocket"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	var ctx = context.Background()
	zlog := config.NewLogger()

	cfg, err := config.NewViper()
	if err != nil {
		zlog.Err(err)
		return
	}

	app, err := initialize.Bootstrap(ctx, cfg, zlog)
	if err != nil {
		zlog.Err(err)
		return
	}

	rsvl, err := resolver.NewResolver(app.UcUser)
	if err != nil {
		zlog.Err(err)
		return
	}

	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: rsvl}))
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")
				if origin == "" || origin == r.Header.Get("Host") {
					return true
				}

				return slices.Contains([]string{"http://localhost:8000"}, origin)
			},
		},
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	middleware.ApplyAuthMiddleware(srv, cfg.Settings.JWTSecret)

	address := cfg.Server.Address
	if address == "" {
		address = defaultPort
	}

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	zlog.Info().Msgf("connect to http://localhost:%s for GraphQL playground", address)
	log.Fatal(http.ListenAndServe(":"+address, nil))
}
