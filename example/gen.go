package main

import (
	"context"
	"fmt"
	"mime"
	"os"

	pixai_client "github.com/pixai-art/pixai-client-go"
)

func main() {
	// Initialize the client
	client := pixai_client.NewPixAIClient().
		SetApiKey(os.Getenv("API_KEY")).
		Init()

	ctx := context.Background()

	// Prepare the parameters for the request with JSONObject structure.
	// You can learn the detail about the parameters from our API documentation.
	params := pixai_client.JSONObject{
		"width":   512,
		"height":  512,
		"prompts": "miku",
		"modelId": "1648918127446573124",
	}

	var taskId string

	{
		// Here is an example of how to generate an image and wait for the result
		task, err := client.GenerateImage(ctx, params, func(task *pixai_client.TaskBase) {
			fmt.Printf("Task: %+v\n", task)
		})
		fmt.Printf("Task: %+v\n", task)
		fmt.Printf("Error: %+v\n", err)
		taskId = task.Id

		media, err := client.GetMediaFromTask(ctx, task)
		if err != nil {
			fmt.Printf("Error: %+v\n", err)
			return
		}
		file, mimeType, err := client.DownloadMedia(ctx, media[0])
		if err != nil {
			fmt.Printf("Error: %+v\n", err)
			return
		}
		ext, _ := mime.ExtensionsByType(mimeType)
		err = os.WriteFile("output"+ext[0], file, 0644)
		if err != nil {
			fmt.Printf("Error: %+v\n", err)
			return
		}
	}

	{
		// Here is an example of how to get the result of a task
		task, err := client.GetTaskById(ctx, taskId)
		fmt.Printf("Task: %+v\n", task)
		fmt.Printf("Error: %+v\n", err)
	}

	{
		// If you don't want to wait for the result, you can use the following code to generate an image and get the task id.
		// Then you can poll the task status by task id or use our webhook to get the result.
		task, err := client.CreateGenerationTask(ctx, params)
		if err != nil {
			fmt.Printf("Error: %+v\n", err)
		}
		fmt.Printf("Task: %+v\n", task)
	}
}
