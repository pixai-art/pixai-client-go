package pixai_client

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

func (p *PixAIClient) CreateGenerationTask(ctx context.Context, parameters JSONObject) (*TaskBase, error) {
	var m struct {
		CreateGenerationTask *TaskBase `graphql:"createGenerationTask(parameters: $parameters)"`
	}
	err := p.graphqlClient.Mutate(ctx, &m, map[string]any{
		"parameters": parameters,
	}, graphql.OperationName("createGenerationTask"))
	if err != nil {
		return nil, err
	}
	return m.CreateGenerationTask, nil
}

func (p *PixAIClient) CancelGenerationTask(ctx context.Context, id string) (*TaskBase, error) {
	var m struct {
		CancelGenerationTask *TaskBase `graphql:"cancelGenerationTask(id: $id)"`
	}
	err := p.graphqlClient.Mutate(ctx, &m, map[string]any{
		"id": graphql.ID(id),
	}, graphql.OperationName("cancelGenerationTask"))
	if err != nil {
		return nil, err
	}
	return m.CancelGenerationTask, nil
}
