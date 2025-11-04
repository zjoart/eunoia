package a2a

// A2ARequest represents an incoming A2A message request
type A2ARequest struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      string    `json:"id"`
	Method  string    `json:"method"`
	Params  A2AParams `json:"params"`
}

// A2AParams contains the parameters for an A2A request
type A2AParams struct {
	Message       A2AMessage `json:"message"`
	Configuration A2AConfig  `json:"configuration"`
}

// A2AMessage represents a message in the A2A protocol
type A2AMessage struct {
	Kind      string                 `json:"kind"`
	Role      string                 `json:"role"`
	Parts     []A2APart              `json:"parts"`
	Metadata  map[string]interface{} `json:"metadata"`
	MessageID string                 `json:"messageId"`
}

// A2APart represents a part of an A2A message
type A2APart struct {
	Kind string    `json:"kind"`
	Text string    `json:"text,omitempty"`
	Data []A2APart `json:"data,omitempty"`
}

// A2AConfig represents configuration for A2A requests
type A2AConfig struct {
	AcceptedOutputModes    []string            `json:"acceptedOutputModes"`
	HistoryLength          int                 `json:"historyLength"`
	PushNotificationConfig A2APushNotification `json:"pushNotificationConfig"`
	Blocking               bool                `json:"blocking"`
}

// A2APushNotification represents push notification configuration
type A2APushNotification struct {
	URL            string                 `json:"url"`
	Token          string                 `json:"token"`
	Authentication map[string]interface{} `json:"authentication"`
}

// A2AResponse represents an A2A response
type A2AResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	ID      string    `json:"id"`
	Result  A2AResult `json:"result,omitempty"`
	Error   *A2AError `json:"error,omitempty"`
}

// A2AResult represents the task result (the top-level result object)
type A2AResult struct {
	ID        string             `json:"id"`
	ContextID string             `json:"contextId,omitempty"`
	Status    A2ATaskStatus      `json:"status"`
	Artifacts []A2AArtifact      `json:"artifacts,omitempty"`
	History   []A2AMessageResult `json:"history,omitempty"`
	Kind      string             `json:"kind"`
}

// A2ATaskStatus represents the status of a task with embedded message
type A2ATaskStatus struct {
	State     string           `json:"state"`
	Timestamp string           `json:"timestamp"`
	Message   A2AMessageResult `json:"message"`
}

// A2AMessageResult represents a message in A2A responses
type A2AMessageResult struct {
	MessageID string                 `json:"messageId"`
	Role      string                 `json:"role"`
	Parts     []A2APart              `json:"parts"`
	Kind      string                 `json:"kind"`
	TaskID    string                 `json:"taskId"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// A2AArtifact represents an artifact in A2A responses
type A2AArtifact struct {
	ArtifactID string    `json:"artifactId"`
	Name       string    `json:"name"`
	Parts      []A2APart `json:"parts"`
}

// A2AError represents an error in A2A responses
type A2AError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ChatResponse represents the internal chat response
type ChatResponse struct {
	Response string `json:"response"`
}
