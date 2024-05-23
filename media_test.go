package pixai_client_test

import (
	"context"
	"os"
	"testing"

	pixai_client "github.com/pixai-art/pixai-client-go"
)

func TestUploadMedia(t *testing.T) {
	client := pixai_client.NewPixAIClient().
		SetApiKey(os.Getenv("PIXAI_API_KEY")).
		Init()

	ctx := context.Background()
	res, err := client.UploadMediaUrl(ctx, "https://i.imgur.com/n2fiOiC.png")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Logf("upload media result: %v", res)
}
