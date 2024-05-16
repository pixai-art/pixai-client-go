package pixai_client

import (
	"time"
)

type WebhookEventType string

const (
	WebhookEventType_TaskUpdated   WebhookEventType = "task_updated"
	WebhookEventType_TaskRunning   WebhookEventType = "task_running"
	WebhookEventType_TaskCanceled  WebhookEventType = "task_canceled"
	WebhookEventType_TaskCompleted WebhookEventType = "task_completed"
	WebhookEventType_TaskFailed    WebhookEventType = "task_failed"
)

type WebhookEvent struct {
	Id         string
	UserId     string
	WebhookId  string
	WebhookUrl string
	Type       WebhookEventType
	RetryCount int
	Data       any
	CreatedAt  *time.Time

	Task *TaskBase
}

func (p *PixAIClient) mapToStruct(m any, s interface{}) error {
	data, err := p.marshalJSON(m)
	if err != nil {
		return err
	}
	return p.unmarshalJSON(data, s)
}

func (p *PixAIClient) ParseWebhookPayload(data []byte) (*WebhookEvent, error) {
	var webhookEvent WebhookEvent
	err := p.unmarshalJSON(data, &webhookEvent)
	if err != nil {
		return nil, err
	}

	switch webhookEvent.Type {
	case WebhookEventType_TaskUpdated,
		WebhookEventType_TaskRunning,
		WebhookEventType_TaskCanceled,
		WebhookEventType_TaskCompleted,
		WebhookEventType_TaskFailed:
		var task TaskBase
		err := p.mapToStruct(webhookEvent.Data, &task)
		if err != nil {
			return nil, err
		}
		webhookEvent.Task = &task
	}

	return &webhookEvent, nil
}
