package nakama

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/coder/websocket"
	"github.com/gwaylib/errors"
	"github.com/gwaylib/log"
	"github.com/heroiclabs/nakama-common/rtapi"
	"google.golang.org/protobuf/encoding/protojson"
)

// WebSocketAdapter is a text-based WebSocket adapter for transmitting payloads over UTF-8.
type WebSocketAdapter struct {
	uri       string
	socket    *websocket.Conn
	onClose   func(err error)
	onError   func(err error)
	onMessage func(mType int, message []byte)
	onOpen    func(event interface{}) error
	mu        sync.Mutex // To guard websocket connection reference
}

// NewWebSocketAdapterText creates a new instance of WebSocketAdapter.
func NewWebSocketAdapterText() *WebSocketAdapter {
	return &WebSocketAdapter{}
}

// IsOpen determines if the WebSocket connection is open.
func (w *WebSocketAdapter) IsOpen() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.socket != nil
}

// Close closes the WebSocket connection.
func (w *WebSocketAdapter) Close() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.socket != nil {
		_ = w.socket.Close(websocket.StatusNormalClosure, "Client closed connection")
		w.socket = nil
	}
}

// Connect connects to the WebSocket using the specified arguments.
func (w *WebSocketAdapter) Connect(scheme, host, port string, createStatus bool, token string) error {
	w.uri = fmt.Sprintf("%s%s:%s/ws?lang=en&status=%s&token=%s",
		scheme,
		host,
		port,
		url.QueryEscape(fmt.Sprintf("%v", createStatus)),
		url.QueryEscape(token),
	)
	return w.connect()
}

func (w *WebSocketAdapter) connect() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	var err error

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	w.socket, _, err = websocket.Dial(ctx, w.uri, nil)
	if err != nil {
		return err
	}

	go w.listen()

	return nil
}

// Send sends a message through the WebSocket connection.
func (w *WebSocketAdapter) Send(message *rtapi.Envelope) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.socket == nil {
		return fmt.Errorf("WebSocket is not connected")
	}

	//msgBytes, err := json.Marshal(message)
	msgBytes, err := protojson.Marshal(message)
	if err != nil {
		return errors.As(err)
	}

	// ctx, cancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := w.socket.Write(ctx, websocket.MessageText, msgBytes); err != nil {
		return errors.As(err)
	}

	return nil
}

// ReadSocketResponse reads a single response message from the WebSocket connection.
func (w *WebSocketAdapter) Read() ([]byte, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.socket == nil {
		return nil, fmt.Errorf("WebSocket is not connected")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, message, err := w.socket.Read(ctx)
	if err != nil {
		return nil, err
	}

	return message, nil
}

// listen listens for messages or errors from the WebSocket server.
func (w *WebSocketAdapter) listen() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for {
		mType, message, err := w.socket.Read(ctx)
		if err != nil {
			w.mu.Lock()
			socket := w.socket
			w.mu.Unlock()

			closeStatus := websocket.CloseStatus(err)

			if socket != nil {
				w.Close()
			}
			if w.onError != nil {
				w.onError(errors.As(err, closeStatus))
			} else {
				log.Infof("WebSocket closed with status: %d, cause:%s", closeStatus, err.Error())
			}
			break
		}
		if w.onMessage == nil {
			// message handler not set
			continue
		}
		w.onMessage(int(mType), message)
		continue
	}
}
