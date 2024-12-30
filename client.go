package nakama

// Default configuration values
const (
	DefaultHost              = "127.0.0.1"
	DefaultPort              = "7350"
	DefaultServerKey         = "defaultkey"
	DefaultTimeoutMs         = 7000
	DefaultExpiredTimespanMs = 5 * 60 * 1000 // 5 minutes in milliseconds
)

// RpcResponse defines the response for an RPC function executed on the server.
type RpcResponse struct {
	// ID is the identifier of the function.
	ID string

	// Payload is the payload of the function, which must be a JSON object.
	Payload map[string]interface{}
}
