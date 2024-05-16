package pixai_client

type TaskStatus string

const (
	TaskStatus_Waiting   TaskStatus = "waiting"
	TaskStatus_Running   TaskStatus = "running"
	TaskStatus_Completed TaskStatus = "completed"
	TaskStatus_Failed    TaskStatus = "failed"
	TaskStatus_Cancelled TaskStatus = "cancelled"
)

const (
	ApiBaseUrl       = "https://api.pixai.art"
	WebsocketBaseUrl = "wss://gw.pixai.art"
)
