package response

type PingResponse struct {
	Status          string `json:"status"`
	ServerTimestamp int64  `json:"server_timestamp"`
	AppName         string `json:"app_name"`
	Environment     string `json:"environment"`
}
