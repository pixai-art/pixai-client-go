package pixai_client

import (
	"github.com/hasura/go-graphql-client"
	"github.com/hasura/go-graphql-client/pkg/jsonutil"
)

type AllEvents struct {
	TaskUpdated *TaskBase `json:"taskUpdated"`
}

type subscribePersonalEvents struct {
	PersonalEvents *AllEvents `graphql:"personalEvents"`
}

func (p *PixAIClient) PersonalEvents() (*eventSource, error) {
	if p.personalEvents == nil {
		ch := make(chan *AllEvents)

		sid, err := p.subscriptionClient.Subscribe(&subscribePersonalEvents{}, nil, func(message []byte, err error) error {
			if err != nil {
				return nil
			}
			var data subscribePersonalEvents
			if err := jsonutil.UnmarshalGraphQL(message, &data); err != nil {
				return nil
			}

			ch <- data.PersonalEvents

			return nil
		}, graphql.OperationName("subscribePersonalEvents"))

		go p.subscriptionClient.Run()

		if err != nil {
			return nil, err
		}
		p.personalEvents = &eventSource{
			sid:      sid,
			ch:       ch,
			handlers: make(map[int]*handler),
		}
		p.personalEvents.Start()
	}
	return p.personalEvents, nil
}
