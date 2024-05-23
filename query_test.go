package pixai_client_test

import (
	"context"
	"os"
	"testing"

	pixai_client "github.com/pixai-art/pixai-client-go"
)

func TestGetMediaById(t *testing.T) {
	client := pixai_client.NewPixAIClient().
		SetApiKey(os.Getenv("PIXAI_API_KEY")).
		Init()

	ctx := context.Background()
	media, err := client.GetMediaById(ctx, "1")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("media: %v", media)
}
