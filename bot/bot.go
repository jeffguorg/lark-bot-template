package bot

type ChatRequest struct {
	Message   string
	SessionID string
}

type ChatResponse struct {
	Message   string
	SessionID string
}

type Bot interface {
	Chat(req ChatRequest) (*ChatResponse, error)
}
