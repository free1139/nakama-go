package nakama

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"sync"

	"github.com/gorilla/websocket"
)

// WebSocketAdapter defines an interface for a websocket adapter.
type WebSocketAdapter interface {
	IsOpen() bool
	Close()
	Connect(scheme, host, port string, createStatus bool, token string) error
	Send(message interface{}) error

	// Event handlers
	SetOnClose(handler func(evt error))
	SetOnError(handler func(event error))
	SetOnMessage(handler func(message []byte))
	SetOnOpen(handler func(event interface{}))
}

// WebSocketAdapterText is a text-based WebSocket adapter for transmitting payloads over UTF-8.
type WebSocketAdapterText struct {
	socket    *websocket.Conn
	onClose   func(event interface{})
	onError   func(event error)
	onMessage func(message interface{})
	onOpen    func(event interface{})
	mu        sync.Mutex // To guard websocket connection reference
}

// NewWebSocketAdapterText creates a new instance of WebSocketAdapterText.
func NewWebSocketAdapterText() *WebSocketAdapterText {
	return &WebSocketAdapterText{}
}

// IsOpen determines if the WebSocket connection is open.
func (w *WebSocketAdapterText) IsOpen() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.socket != nil && w.socket.CloseHandler() == nil
}

// Close closes the WebSocket connection.
func (w *WebSocketAdapterText) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.socket != nil {
		_ = w.socket.Close()
		w.socket = nil
	}
}

// Connect connects to the WebSocket using the specified arguments.
func (w *WebSocketAdapterText) Connect(scheme, host, port string, createStatus bool, token string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	urlStr := fmt.Sprintf("%s%s:%s/ws?lang=en&status=%s&token=%s",
		scheme,
		host,
		port,
		url.QueryEscape(fmt.Sprintf("%v", createStatus)),
		url.QueryEscape(token),
	)

	var err error
	w.socket, _, err = websocket.DefaultDialer.Dial(urlStr, nil)
	if err != nil {
		return err
	}

	if w.onOpen != nil {
		w.onOpen(nil)
	}

	go w.listen()

	return nil
}

// Send sends a message through the WebSocket connection.
func (w *WebSocketAdapterText) Send(message interface{}) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.socket == nil {
		return fmt.Errorf("WebSocket is not connected")
	}

	// Handle specific cases of match_data_send and party_data_send
	if msgMap, ok := message.(map[string]interface{}); ok {
		handleEncodedData(msgMap, "match_data_send")
		handleEncodedData(msgMap, "party_data_send")
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return w.socket.WriteMessage(websocket.TextMessage, msgBytes)
}

// SetOnClose sets the handler for WebSocket close events.
func (w *WebSocketAdapterText) SetOnClose(handler func(event interface{})) {
	w.onClose = handler
}

// SetOnError sets the handler for WebSocket error events.
func (w *WebSocketAdapterText) SetOnError(handler func(event error)) {
	w.onError = handler
}

// SetOnMessage sets the handler for WebSocket message events.
func (w *WebSocketAdapterText) SetOnMessage(handler func(message interface{})) {
	w.onMessage = handler
}

// SetOnOpen sets the handler for WebSocket open events.
func (w *WebSocketAdapterText) SetOnOpen(handler func(event interface{})) {
	w.onOpen = handler
}

// listen listens for messages or errors from the WebSocket server.
func (w *WebSocketAdapterText) listen() {
	for {
		_, message, err := w.socket.ReadMessage()
		if err != nil {
			w.mu.Lock()
			socket := w.socket
			w.mu.Unlock()

			if socket != nil {
				if w.onError != nil {
					w.onError(err)
				}
				if websocket.IsUnexpectedCloseError(err) && w.onClose != nil {
					w.onClose(nil)
				}
				w.Close()
			}
			break
		}

		var decodedMessage map[string]interface{}
		if err := json.Unmarshal(message, &decodedMessage); err != nil {
			if w.onError != nil {
				w.onError(err)
			}
			continue
		}

		// Handle specific decoding logic for match_data and party_data
		decodeReceivedData(decodedMessage, "match_data")
		decodeReceivedData(decodedMessage, "party_data")

		if w.onMessage != nil {
			w.onMessage(decodedMessage)
		}
	}
}

// handleEncodedData handles encoding of match_data_send and party_data_send fields.
func handleEncodedData(msg map[string]interface{}, field string) {
	if sendData, exists := msg[field]; exists {
		if sendMap, ok := sendData.(map[string]interface{}); ok {
			// Convert op_code to string
			if opCode, ok := sendMap["op_code"]; ok {
				sendMap["op_code"] = fmt.Sprintf("%v", opCode)
			}

			// Encode data
			if payload, exists := sendMap["data"]; exists {
				switch v := payload.(type) {
				case []byte:
					sendMap["data"] = base64.StdEncoding.EncodeToString(v)
				case string:
					sendMap["data"] = base64.StdEncoding.EncodeToString([]byte(v))
				}
			}
		}
	}
}

// decodeReceivedData decodes the match_data and party_data fields in messages received from the server.
func decodeReceivedData(msg map[string]interface{}, field string) {
	if data, exists := msg[field]; exists {
		if dataMap, ok := data.(map[string]interface{}); ok {
			if encoded, exists := dataMap["data"]; exists {
				if encodedStr, ok := encoded.(string); ok {
					decodedBytes, err := base64.StdEncoding.DecodeString(encodedStr)
					if err == nil {
						dataMap["data"] = decodedBytes
					}
				}
			}
		}
	}
}
