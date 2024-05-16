package pixai_client

import "time"

type MediaBase struct {
	Id     string `json:"id"`
	Type   string `json:"type"`
	Width  *int   `json:"width"`
	Height *int   `json:"height"`
	Urls   []struct {
		Variant string `json:"variant"`
		Url     string `json:"url"`
	} `json:"urls"`
	ImageType    *string  `json:"imageType"`
	FileUrl      *string  `json:"fileUrl"`
	Duration     *float64 `json:"duration"`
	ThumbnailUrl *string  `json:"thumbnailUrl"`
	HlsUrl       *string  `json:"hlsUrl"`
	Size         *int     `json:"size"`
}

type TaskBase struct {
	Id         string         `json:"id"`
	UserId     string         `json:"userId"`
	Parameters map[string]any `json:"parameters" scalar:"true"`
	Outputs    map[string]any `json:"outputs" scalar:"true"`
	Status     *string        `json:"status"`
	StartedAt  *time.Time     `json:"startedAt"`
	EndAt      *time.Time     `json:"endAt"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	Media      *MediaBase     `json:"media"`
}
