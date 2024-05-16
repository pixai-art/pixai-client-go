package pixai_client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hasura/go-graphql-client"
)

type MarshalJSON = func(v interface{}) ([]byte, error)
type UnmarshalJSON = func(data []byte, v interface{}) error

type PixAIClient struct {
	initiated          bool
	httpClient         *http.Client
	graphqlClient      *graphql.Client
	subscriptionClient *graphql.SubscriptionClient

	personalEvents *eventSource

	apiKey           string
	apiBaseUrl       string
	websocketBaseUrl string

	marshalJSON   MarshalJSON
	unmarshalJSON UnmarshalJSON
}

func (p *PixAIClient) SetApiKey(apiKey string) *PixAIClient {
	p.apiKey = apiKey
	return p
}

func (p *PixAIClient) SetApiBaseUrl(apiBaseUrl string) *PixAIClient {
	p.apiBaseUrl = apiBaseUrl
	return p
}

func (p *PixAIClient) SetWebSocketBaseUrl(websocketBaseUrl string) *PixAIClient {
	p.websocketBaseUrl = websocketBaseUrl
	return p
}

func (p *PixAIClient) SetMarshalJSON(marshal MarshalJSON) *PixAIClient {
	p.marshalJSON = marshal
	return p
}

func (p *PixAIClient) SetUnmarshalJSON(unmarshal UnmarshalJSON) *PixAIClient {
	p.unmarshalJSON = unmarshal
	return p
}

func (p *PixAIClient) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.apiKey))
	return p.httpClient.Do(req)
}

func (p *PixAIClient) Init() *PixAIClient {
	if p.initiated {
		return p
	}
	p.initiated = true

	p.httpClient = &http.Client{}

	p.graphqlClient = graphql.NewClient(fmt.Sprintf("%s/graphql", p.apiBaseUrl), p)
	p.subscriptionClient = graphql.NewSubscriptionClient(fmt.Sprintf("%s/graphql", p.websocketBaseUrl)).
		WithConnectionParams(map[string]interface{}{
			"token": p.apiKey,
		}).
		WithProtocol(graphql.GraphQLWS).
		OnDisconnected(func() {
			if p.personalEvents != nil {
				p.personalEvents.Close()
				p.personalEvents = nil
			}
		})

	return p
}

func NewPixAIClient() *PixAIClient {
	client := &PixAIClient{
		apiBaseUrl:       ApiBaseUrl,
		websocketBaseUrl: WebsocketBaseUrl,

		marshalJSON:   json.Marshal,
		unmarshalJSON: json.Unmarshal,
	}

	return client
}

func (p *PixAIClient) Close() {
	p.subscriptionClient.Close()
}

func (p *PixAIClient) GenerateImage(ctx context.Context, parameters JSONObject, onUpdate func(*TaskBase)) (*TaskBase, error) {
	task, err := p.CreateGenerationTask(ctx, parameters)
	if err != nil {
		return nil, err
	}

	personalEvents, err := p.PersonalEvents()
	if err != nil {
		return task, err
	}

	completed := make(chan *TaskBase)

	sid := personalEvents.AddEventListener(func(e *AllEvents) {
		if e != nil && e.TaskUpdated != nil && e.TaskUpdated.Id == task.Id {
			if onUpdate != nil {
				onUpdate(e.TaskUpdated)
			}
			if e.TaskUpdated.Status != nil {
				switch TaskStatus(*e.TaskUpdated.Status) {
				case TaskStatus_Waiting, TaskStatus_Running:
					break
				default:
					completed <- e.TaskUpdated
				}
			}
		}
	})

	defer personalEvents.RemoveEventListener(sid)

	return <-completed, nil
}
