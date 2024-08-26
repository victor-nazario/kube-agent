package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.49

import (
	"context"

	"github.com/victor-nazario/kube-agent/graph/model"
)

// DeployJob is the resolver for the deployJob field.
func (r *mutationResolver) DeployJob(ctx context.Context, input model.DeployJobInput) (string, error) {
	return r.Agent.Deploy(ctx, input.Name, input.Image, input.Command, input.NameSpace, int32(input.BackOffLimit))
}

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
