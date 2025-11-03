package a2a

// JSON-RPC 2.0 Error Codes
const (
	// ParseError indicates invalid JSON was received by the server
	ParseError = -32700

	// InvalidRequest indicates the JSON sent is not a valid Request object
	InvalidRequest = -32600

	// MethodNotFound indicates the method does not exist or is not available
	MethodNotFound = -32601

	// InvalidParams indicates invalid method parameter(s)
	InvalidParams = -32602

	// InternalError indicates an internal JSON-RPC error
	InternalError = -32603
)
