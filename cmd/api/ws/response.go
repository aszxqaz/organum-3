package ws

type WsResponse struct {
	Method  string `json:"method"`
	Success bool   `json:"success"`
	Error   string `json:"error,omitempty"`
}
