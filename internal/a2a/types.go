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

// A2AResult contains the result of an A2A request
type A2AResult struct {
	Message A2AMessageResult `json:"message"`
}

// A2AMessageResult represents a message result in A2A responses
type A2AMessageResult struct {
	Kind      string                 `json:"kind"`
	Role      string                 `json:"role"`
	Parts     []A2APart              `json:"parts"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	MessageID string                 `json:"messageId"`
}

// A2AError represents an error in A2A responses
type A2AError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ChatResponse represents the internal chat response
type ChatResponse struct {
	Response  string `json:"response"`
	MessageID string `json:"message_id"`
}
