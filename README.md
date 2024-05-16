# PixAI Client Go

### Installation

```bash
go get github.com/pixai-art/pixai-client-go
```

### Usage

```go
package main

import (
	"context"
	"fmt"

	pixai_client "github.com/pixai-art/pixai-client-go"
)

func main() {
	// Initialize the client
	client := pixai_client.NewPixAIClient().
		SetApiKey("replace to your api keys").
		Init()

	ctx := context.Background()

	// Prepare the parameters for the request with JSONObject structure.
	// You can learn the detail about the parameters from our API documentation.
	params := pixai_client.JSONObject{
		"width":   512,
		"height":  512,
		"prompts": "miku",
		"modelId": "1648918127446573124",

		// Specific the priority value as 1000 to make sure the task will be processed immediately.
		// This is a historical legacy issue. If this parameter is not passed, the public queue will be used for queuing. We will set high priority tasks as the default value later.
		"priority": 1000,
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

```
