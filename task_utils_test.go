package pixai_client_test

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	pixai_client "github.com/pixai-art/pixai-client-go"
)

func TestGetMediaFromTask(t *testing.T) {
	client := pixai_client.NewPixAIClient().
		SetApiKey(os.Getenv("PIXAI_API_KEY")).
		Init()

	ctx := context.Background()

	params := pixai_client.JSONObject{
		"width":     512,
		"height":    512,
		"prompts":   "miku",
		"modelId":   "1648918127446573124",
		"batchSize": 4,
	}

	task, err := client.GenerateImage(ctx, params, func(task *pixai_client.TaskBase) {
		t.Logf("Task: %+v\n", task)
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	taskData, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if err := json.Unmarshal(taskData, &task); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	media, err := client.GetMediaFromTask(ctx, task)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	t.Logf("Media: %+v\n", media)
}
