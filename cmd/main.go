package main

import (
	"github.com/victor-nazario/kube-agent/internal/agent"
	"github.com/victor-nazario/kube-agent/internal/auth"
	"github.com/victor-nazario/kube-agent/internal/authz"
	"github.com/victor-nazario/kube-agent/internal/user"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/go-chi/chi/v5"
	"github.com/victor-nazario/kube-agent/graph"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	agt, err := agent.NewAgent()
	if err != nil {
		log.Fatal(err)
	}

	authorizer, err := authz.NewAuthorizer(user.ActiveUsers())
	if err != nil {
		log.Fatal(err)
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Agent: agt,
	}}))

	router := chi.NewRouter()
	router.Use(auth.Authentication())
	router.Use(authz.HasPermissions(authorizer))

	router.Handle("/query", srv)

	log.Fatal(http.ListenAndServe(":"+port, router))
}
