package pixai_client

import (
	"context"

	"github.com/hasura/go-graphql-client"
)

func (p *PixAIClient) GetTaskById(ctx context.Context, id string) (*TaskBase, error) {
	var q struct {
		Task *TaskBase `graphql:"task(id: $id)"`
	}
	err := p.graphqlClient.Query(ctx, &q, map[string]any{
		"id": graphql.ID(id),
	}, graphql.OperationName("getTaskById"))
	if err != nil {
		return nil, err
	}
	return q.Task, nil
}
