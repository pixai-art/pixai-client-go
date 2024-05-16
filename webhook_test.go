package pixai_client

import (
	"encoding/json"
	"testing"
)

func TestWebhookParsePayload(t *testing.T) {
	p := &PixAIClient{
		unmarshalJSON: json.Unmarshal,
		marshalJSON:   json.Marshal,
	}
	event, err := p.ParseWebhookPayload([]byte(`{
		"id": "1",
		"userId": "2",
		"type": "task_completed",
		"data": {
			"id": "3"
		}
	}`))
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("event: %v", event)
}
