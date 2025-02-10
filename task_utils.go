package pixai_client

import (
	"context"
	"fmt"
	"sync"
)

func (p *PixAIClient) getMediaFromBatchItem(ctx context.Context, item any) (*MediaBase, error) {
	o := ToJSONObject(item)
	if o == nil {
		return nil, fmt.Errorf("unexpected item format")
	}

	mediaId := o.GetString("mediaId")

	if mediaId == "" {
		return nil, fmt.Errorf("unexpected item mediaId format")
	}

	return p.GetMediaById(ctx, mediaId)
}

func (p *PixAIClient) GetMediaFromTask(ctx context.Context, task *TaskBase) ([]*MediaBase, error) {
	if task == nil {
		return nil, fmt.Errorf("task is nil")
	}

	if task.Status == nil || TaskStatus(*task.Status) != TaskStatus_Completed {
		return nil, fmt.Errorf("task is not completed")
	}

	var outputs JSONObject = task.Outputs

	if outputs == nil {
		return nil, fmt.Errorf("task outputs is nil")
	}

	var lastErr error

	if batch := outputs.GetArray("batch"); batch != nil {
		if len(batch) == 0 {
			return nil, fmt.Errorf("task outputs batch is empty")
		}

		mediaList := make([]*MediaBase, len(batch))

		var wg sync.WaitGroup
		wg.Add(len(batch))

		for idx, item := range batch {
			go func(i int, m any) {
				defer wg.Done()

				media, err := p.getMediaFromBatchItem(ctx, m)

				if err != nil {
					lastErr = err
					return
				}

				mediaList[i] = media
			}(idx, item)
		}

		wg.Wait()

		return mediaList, lastErr
	}

	mediaId := outputs.GetString("mediaId")

	if mediaId != "" {
		media, err := p.GetMediaById(ctx, mediaId)

		if err != nil {
			return nil, err
		}

		return []*MediaBase{media}, nil
	}

	return nil, nil
}
