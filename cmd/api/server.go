package main

import (
	"context"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
	"go.uber.org/zap"
	"log"
	"net/http"
	"os"
	graph2 "realtime-bike-go/cmd/api/graph"
	"realtime-bike-go/config"
	"realtime-bike-go/internal/store"
	"realtime-bike-go/pkg/rblogger"
)

const defaultPort = "8080"

var (
	configPrefix = "rb-graphql"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	ctx := context.Background()
	logger := rblogger.New()
	cfg := config.NewBase(ctx, configPrefix)
	st, err := store.New(ctx, logger, cfg.URI, cfg.DatabaseName, "graphql", cfg.ServerSelectionTimeout)
	if err != nil {
		logger.Error("could not init store", zap.Error(err))
		return
	}

	srv := handler.New(graph2.NewExecutableSchema(graph2.Config{Resolvers: &graph2.Resolver{
		Store: st,
	}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
