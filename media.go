package pixai_client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/hasura/go-graphql-client"
)

type MediaType string

const (
	MediaType_Image MediaType = "IMAGE"
)

type MediaProvider string

const (
	MediaProvider_S3 MediaProvider = "S3"
)

const (
	defaultMediaType = MediaType_Image
	defaultProvider  = MediaProvider_S3
)

type UploadMediaInput struct {
	Type       MediaType     `json:"type"`
	Provider   MediaProvider `json:"provider"`
	ExternalId *string       `json:"externalId,omitempty"`
	ImageType  *string       `json:"imageType,omitempty"`
}

type UploadMediaResult struct {
	UploadUrl  *string    `json:"uploadUrl,omitempty"`
	ExternalId *string    `json:"externalId,omitempty"`
	MediaId    *string    `json:"mediaId,omitempty"`
	Media      *MediaBase `json:"media,omitempty"`
}

func (p *PixAIClient) uploadMedia(ctx context.Context, input *UploadMediaInput) (*UploadMediaResult, error) {
	var m struct {
		UploadMedia *UploadMediaResult `graphql:"uploadMedia(input: $input)"`
	}
	err := p.graphqlClient.Mutate(ctx, &m, map[string]any{
		"input": *input,
	}, graphql.OperationName("uploadMedia"))
	if err != nil {
		return nil, err
	}
	if m.UploadMedia == nil {
		return nil, fmt.Errorf("unexpected response is nil")
	}
	return m.UploadMedia, nil
}

func (p *PixAIClient) getUploadUrl(ctx context.Context) (uploadUrl string, externalId *string, err error) {
	res, err := p.uploadMedia(ctx, &UploadMediaInput{
		Type:     defaultMediaType,
		Provider: defaultProvider,
	})
	if err != nil {
		err = fmt.Errorf("failed to get upload url: %v", err)
		return
	}
	if res.UploadUrl == nil {
		err = fmt.Errorf("unexpected upload url is nil")
		return
	}
	uploadUrl = *res.UploadUrl
	externalId = res.ExternalId
	return
}

func (p *PixAIClient) registerMedia(ctx context.Context, externalId *string) (*MediaBase, error) {
	res, err := p.uploadMedia(ctx, &UploadMediaInput{
		Type:       defaultMediaType,
		Provider:   defaultProvider,
		ExternalId: externalId,
	})
	if err != nil {
		return nil, err
	}
	if res.Media == nil {
		return nil, fmt.Errorf("unexpected media is nil")
	}
	return res.Media, nil
}

func (p *PixAIClient) UploadMediaFile(ctx context.Context, file io.Reader, mimeType string) (*MediaBase, error) {
	uploadUrl, externalId, err := p.getUploadUrl(ctx)
	if err != nil {
		return nil, err
	}
	c, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	req, err := http.NewRequestWithContext(c, http.MethodPut, uploadUrl, file)
	req.Header.Set("Content-Type", mimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	if res, err := p.httpClient.Do(req); err != nil {
		return nil, fmt.Errorf("failed to upload file: %v", err)
	} else {
		defer res.Body.Close()
		if res.StatusCode > 299 {
			body, _ := io.ReadAll(res.Body)
			return nil, fmt.Errorf("unexpected status code on upload image. status=%d body=%s", res.StatusCode, string(body))
		}
	}

	return p.registerMedia(ctx, externalId)
}

func (p *PixAIClient) UploadMediaUrl(ctx context.Context, url string) (*MediaBase, error) {
	c, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	req, err := http.NewRequestWithContext(c, http.MethodGet, url, nil)
	req.Header.Set("User-Agent", "PixAIClientGo/1.0.0")
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}
	res, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %v", err)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if res.StatusCode > 299 {
		return nil, fmt.Errorf("unexpected status code on download media. status=%d", res.StatusCode)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %v", err)
	}
	return p.UploadMediaFile(ctx, bytes.NewReader(body), res.Header.Get("Content-Type"))
}
