package main

import (
	"context"
	"fmt"

	pixai_client "github.com/pixai-art/pixai-client-go"
)

func main() {
	client := pixai_client.NewPixAIClient().
		SetApiKey("replace to your api keys").
		Init()

	ctx := context.Background()

	params := pixai_client.JSONObject{
		"width":    512,
		"height":   512,
		"prompts":  "miku",
		"modelId":  "1648918127446573124",
		"priority": 1000,
	}

	var taskId string

	{
		task, err := client.GenerateImage(ctx, params, func(task *pixai_client.TaskBase) {
			fmt.Printf("Task: %+v\n", task)
		})
		fmt.Printf("Task: %+v\n", task)
		fmt.Printf("Error: %+v\n", err)
		taskId = task.Id
	}

	{
		task, err := client.GetTaskById(ctx, taskId)
		fmt.Printf("Task: %+v\n", task)
		fmt.Printf("Error: %+v\n", err)
	}

}
