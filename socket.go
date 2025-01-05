package nakama

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type PromiseExecutor struct {
	Resolve func(value interface{})
	Reject  func(reason error)
}

type Presence struct {
	UserID    string `json:"user_id"`
	SessionID string `json:"session_id"`
	Username  string `json:"username"`
	Node      string `json:"node"`
}

type Channel struct {
	ID        string     `json:"id"`
	Presences []Presence `json:"presences"`
	Self      Presence   `json:"self"`
}

type ChannelJoin struct {
	ChannelJoin struct {
		Target      string `json:"target"`
		Type        int    `json:"type"`
		Persistence bool   `json:"persistence"`
		Hidden      bool   `json:"hidden"`
	} `json:"channel_join"`
}

type ChannelLeave struct {
	ChannelLeave struct {
		ChannelID string `json:"channel_id"`
	} `json:"channel_leave"`
}

type ChannelMessageAck struct {
	ChannelID   string `json:"channel_id"`
	MessageID   string `json:"message_id"`
	Code        int    `json:"code"`
	Username    string `json:"username"`
	CreateTime  string `json:"create_time"`
	UpdateTime  string `json:"update_time"`
	Persistence bool   `json:"persistence"`
}

type ChannelMessageSend struct {
	ChannelMessageSend struct {
		ChannelID string      `json:"channel_id"`
		Content   interface{} `json:"content"`
	} `json:"channel_message_send"`
}

type ChannelMessageUpdate struct {
	ChannelMessageUpdate struct {
		ChannelID string      `json:"channel_id"`
		MessageID string      `json:"message_id"`
		Content   interface{} `json:"content"`
	} `json:"channel_message_update"`
}

type ChannelMessageRemove struct {
	ChannelMessageRemove struct {
		ChannelID string `json:"channel_id"`
		MessageID string `json:"message_id"`
	} `json:"channel_message_remove"`
}

type ChannelPresenceEvent struct {
	ChannelID string     `json:"channel_id"`
	Joins     []Presence `json:"joins"`
	Leaves    []Presence `json:"leaves"`
}

type StreamId struct {
	Mode       int    `json:"mode"`
	Subject    string `json:"subject"`
	Subcontext string `json:"subcontext"`
	Label      string `json:"label"`
}

type StreamData struct {
	Stream   StreamId  `json:"stream"`
	Sender   *Presence `json:"sender,omitempty"`
	Data     string    `json:"data"`
	Reliable *bool     `json:"reliable,omitempty"`
}

type StreamPresenceEvent struct {
	Stream StreamId   `json:"stream"`
	Joins  []Presence `json:"joins"`
	Leaves []Presence `json:"leaves"`
}

type MatchPresenceEvent struct {
	MatchID string     `json:"match_id"`
	Joins   []Presence `json:"joins"`
	Leaves  []Presence `json:"leaves"`
}

type MatchmakerAdd struct {
	MatchmakerAdd struct {
		MinCount          int                `json:"min_count"`
		MaxCount          int                `json:"max_count"`
		Query             string             `json:"query"`
		StringProperties  map[string]string  `json:"string_properties,omitempty"`
		NumericProperties map[string]float64 `json:"numeric_properties,omitempty"`
	} `json:"matchmaker_add"`
}

type MatchmakerTicket struct {
	Ticket string `json:"ticket"`
}

type MatchmakerRemove struct {
	MatchmakerRemove struct {
		Ticket string `json:"ticket"`
	} `json:"matchmaker_remove"`
}

type MatchmakerUser struct {
	Presence          Presence           `json:"presence"`
	PartyID           string             `json:"party_id"`
	StringProperties  map[string]string  `json:"string_properties,omitempty"`
	NumericProperties map[string]float64 `json:"numeric_properties,omitempty"`
}

type MatchmakerMatched struct {
	Ticket  string           `json:"ticket"`
	MatchID string           `json:"match_id"`
	Token   string           `json:"token"`
	Users   []MatchmakerUser `json:"users"`
	Self    MatchmakerUser   `json:"self"`
}

type Match struct {
	MatchID       string     `json:"match_id"`
	Authoritative bool       `json:"authoritative"`
	Label         *string    `json:"label,omitempty"`
	Size          int        `json:"size"`
	Presences     []Presence `json:"presences"`
	Self          Presence   `json:"self"`
}

type CreateMatch struct {
	MatchCreate struct {
		Name string `json:"name,omitempty"`
	} `json:"match_create"`
}

type JoinMatch struct {
	MatchJoin struct {
		MatchID  *string                `json:"match_id,omitempty"`
		Token    *string                `json:"token,omitempty"`
		Metadata map[string]interface{} `json:"metadata,omitempty"`
	} `json:"match_join"`
}

type LeaveMatch struct {
	MatchLeave struct {
		MatchID string `json:"match_id"`
	} `json:"match_leave"`
}

type MatchData struct {
	MatchID  string    `json:"match_id"`
	OpCode   int       `json:"op_code"`
	Data     []byte    `json:"data"`
	Presence *Presence `json:"presence,omitempty"`
	Reliable *bool     `json:"reliable,omitempty"`
}

type MatchDataSend struct {
	MatchDataSend struct {
		MatchID   string      `json:"match_id"`
		OpCode    int         `json:"op_code"`
		Data      interface{} `json:"data"`
		Presences []Presence  `json:"presences"`
		Reliable  *bool       `json:"reliable,omitempty"`
	} `json:"match_data_send"`
}

type Party struct {
	PartyID   string     `json:"party_id"`
	Open      bool       `json:"open"`
	MaxSize   int        `json:"max_size"`
	Self      Presence   `json:"self"`
	Leader    Presence   `json:"leader"`
	Presences []Presence `json:"presences"`
}

type PartyCreate struct {
	PartyCreate struct {
		Open    bool `json:"open"`
		MaxSize int  `json:"max_size"`
	} `json:"party_create"`
}

type PartyJoin struct {
	PartyJoin struct {
		PartyID string `json:"party_id"`
	} `json:"party_join"`
}

type PartyLeave struct {
	PartyLeave struct {
		PartyID string `json:"party_id"`
	} `json:"party_leave"`
}

type PartyPromote struct {
	PartyPromote struct {
		PartyID  string   `json:"party_id"`
		Presence Presence `json:"presence"`
	} `json:"party_promote"`
}

type PartyLeader struct {
	PartyID  string   `json:"party_id"`
	Presence Presence `json:"presence"`
}

type PartyAccept struct {
	PartyAccept struct {
		PartyID  string   `json:"party_id"`
		Presence Presence `json:"presence"`
	} `json:"party_accept"`
}

type PartyClose struct {
	PartyClose struct {
		PartyID string `json:"party_id"`
	} `json:"party_close"`
}

type PartyData struct {
	PartyID  string   `json:"party_id"`
	Presence Presence `json:"presence"`
	OpCode   int      `json:"op_code"`
	Data     []byte   `json:"data"`
}

type PartyDataSend struct {
	PartyDataSend struct {
		PartyID string      `json:"party_id"`
		OpCode  int         `json:"op_code"`
		Data    interface{} `json:"data"`
	} `json:"party_data_send"`
}

type PartyJoinRequest struct {
	PartyID   string     `json:"party_id"`
	Presences []Presence `json:"presences"`
}

type PartyJoinRequestList struct {
	PartyJoinRequestList struct {
		PartyID string `json:"party_id"`
	} `json:"party_join_request_list"`
}

type PartyMatchmakerAdd struct {
	PartyMatchmakerAdd struct {
		PartyID           string             `json:"party_id"`
		MinCount          int                `json:"min_count"`
		MaxCount          int                `json:"max_count"`
		Query             string             `json:"query"`
		StringProperties  map[string]string  `json:"string_properties,omitempty"`
		NumericProperties map[string]float64 `json:"numeric_properties,omitempty"`
	} `json:"party_matchmaker_add"`
}

type PartyMatchmakerRemove struct {
	PartyMatchmakerRemove struct {
		PartyID string `json:"party_id"`
		Ticket  string `json:"ticket"`
	} `json:"party_matchmaker_remove"`
}

type PartyMatchmakerTicket struct {
	PartyID string `json:"party_id"`
	Ticket  string `json:"ticket"`
}

type PartyPresenceEvent struct {
	PartyID string     `json:"party_id"`
	Joins   []Presence `json:"joins"`
	Leaves  []Presence `json:"leaves"`
}

type PartyRemove struct {
	PartyRemove struct {
		PartyID  string   `json:"party_id"`
		Presence Presence `json:"presence"`
	} `json:"party_remove"`
}

type Rpc struct {
	Rpc ApiRpc `json:"rpc"`
}

type Ping struct {
	// No fields needed for Ping
}

type Status struct {
	Presences []Presence `json:"presences"`
}

type StatusFollow struct {
	StatusFollow struct {
		UserIDs []string `json:"user_ids"`
	} `json:"status_follow"`
}

type StatusPresenceEvent struct {
	Joins  []Presence `json:"joins"`
	Leaves []Presence `json:"leaves"`
}

type StatusUnfollow struct {
	StatusUnfollow struct {
		UserIDs []string `json:"user_ids"`
	} `json:"status_unfollow"`
}

type StatusUpdate struct {
	StatusUpdate struct {
		Status *string `json:"status,omitempty"`
	} `json:"status_update"`
}

// Socket defines the Go struct with corresponding methods.
type Socket struct {
	OnDisconnect       func(evt error)                `json:"-"`
	OnError            func(evt error)                `json:"-"`
	OnNotification     func(notification string)      `json:"-"`
	OnMatchData        func(matchData MatchData)      `json:"-"`
	OnParty            func(party Party)              `json:"-"`
	OnPartyJoinRequest func(request PartyJoinRequest) `json:"-"`
	OnStreamData       func(data StreamData)          `json:"-"`
	OnHeartbeatTimeout func()                         `json:"-"`
	HeartbeatTimeoutMs int                            `json:"-"`
}

// SocketError represents an error received from a socket message.
type SocketError struct {
	Code    int    `json:"code"`    // The error code
	Message string `json:"message"` // A message in English to help developers debug the response
}

// DefaultSocket constants
const (
	DefaultHeartbeatTimeoutMs = 10000
	DefaultSendTimeoutMs      = 10000
	DefaultConnectTimeoutMs   = 30000
)

// DefaultSocket represents a WebSocket connection to the Nakama server
type DefaultSocket struct {
	Host               string
	Port               string
	UseSSL             bool
	Verbose            bool
	Adapter            WebSocketAdapter
	SendTimeoutMs      int
	HeartbeatTimeoutMs int
	cIds               map[string]*PromiseExecutor
	nextCid            int
}

// NewDefaultSocket creates an instance of DefaultSocket.
func NewDefaultSocket(host, port string, useSSL, verbose bool, adapter WebSocketAdapter) *DefaultSocket {
	return &DefaultSocket{
		Host:               host,
		Port:               port,
		UseSSL:             useSSL,
		Verbose:            verbose,
		Adapter:            adapter,
		SendTimeoutMs:      DefaultSendTimeoutMs,
		HeartbeatTimeoutMs: DefaultHeartbeatTimeoutMs,
		cIds:               make(map[string]*PromiseExecutor),
		nextCid:            1,
	}
}

// GenerateCID generates a unique client ID for requests.
func (socket *DefaultSocket) GenerateCID() string {
	cid := fmt.Sprintf("%d", socket.nextCid)
	socket.nextCid++
	return cid
}

// Connect establishes the WebSocket connection with optional timeouts.
func (socket *DefaultSocket) Connect(session Session, createStatus bool, timeoutMs int) (*Session, error) {
	if socket.Adapter.IsOpen() {
		return &session, nil
	}

	scheme := "ws://"
	if socket.UseSSL {
		scheme = "wss://"
	}

	err := socket.Adapter.Connect(scheme, socket.Host, socket.Port, createStatus, session.Token)
	if err != nil {
		return nil, err
	}

	socket.Adapter.SetOnClose(func(evt error) {
		socket.OnDisconnect(evt)
	})
	socket.Adapter.SetOnError(func(evt error) {
		socket.OnError(evt)
	})
	socket.Adapter.SetOnMessage(func(message []byte) {
		socket.HandleMessage(message)
	})

	// Set a timeout for the connection process
	resChan := make(chan error, 1)
	go func() {
		time.Sleep(time.Duration(timeoutMs) * time.Millisecond)
		resChan <- errors.New("Socket connection timed out")
	}()

	select {
	case err := <-resChan:
		if err != nil {
			socket.Adapter.Close()
			return nil, err
		}
	}

	return &session, nil
}

// Disconnect terminates the WebSocket connection.
func (socket *DefaultSocket) Disconnect(fireDisconnectEvent bool) {
	if socket.Adapter.IsOpen() {
		socket.Adapter.Close()
	}
	if fireDisconnectEvent {
		socket.OnDisconnect(fmt.Errorf("Socket disconnected"))
	}
}

// SetHeartbeatTimeoutMs sets the timeout for heartbeat pings.
func (socket *DefaultSocket) SetHeartbeatTimeoutMs(ms int) {
	socket.HeartbeatTimeoutMs = ms
}

// GetHeartbeatTimeoutMs gets the timeout for heartbeat pings.
func (socket *DefaultSocket) GetHeartbeatTimeoutMs() int {
	return socket.HeartbeatTimeoutMs
}

// OnDisconnect handles WebSocket disconnections.
func (socket *DefaultSocket) OnDisconnect(evt error) {
	if socket.Verbose {
		fmt.Println("OnDisconnect:", evt)
	}
}

// OnError handles WebSocket errors.
func (socket *DefaultSocket) OnError(evt error) {
	if socket.Verbose {
		fmt.Println("OnError:", evt)
	}
}

// HandleMessage processes incoming WebSocket messages.
func (socket *DefaultSocket) HandleMessage(message []byte) {
	var msg map[string]interface{}
	if err := json.Unmarshal(message, &msg); err != nil {
		if socket.Verbose {
			fmt.Println("Failed to parse message:", err)
		}
		return
	}

	if cid, ok := msg["cid"].(string); ok {
		executor, exists := socket.cIds[cid]
		if exists {
			delete(socket.cIds, cid)
			if _, hasError := msg["error"]; hasError {
				executor.Reject(errors.New(msg["error"].(string)))
			} else {
				executor.Resolve(msg)
			}
		} else {
			if socket.Verbose {
				fmt.Println("No promise executor for message CID:", cid)
			}
		}
	} else {
		// Handle different message types here (notifications, match data, etc.)
		if socket.Verbose {
			fmt.Println("Message received:", string(message))
		}
	}
}

// Send sends a message to the WebSocket server with optional timeout.
func (socket *DefaultSocket) Send(message interface{}, sendTimeout int) error {
	if !socket.Adapter.IsOpen() {
		return errors.New("Socket connection is not established")
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("Failed to encode message: %v", err)
	}

	cid := socket.GenerateCID()
	socket.cIds[cid] = &PromiseExecutor{
		Resolve: func(result interface{}) {
			if socket.Verbose {
				fmt.Println("Message sent successfully:", string(data))
			}
		},
		Reject: func(e error) {
			if socket.Verbose {
				fmt.Println("Message failed:", e)
			}
		},
	}

	socket.Adapter.Send(data)

	// Set a timeout for the send operation
	go func(cid string) {
		time.Sleep(time.Duration(sendTimeout) * time.Millisecond)
		delete(socket.cIds, cid)
	}(cid)

	return nil
}

// CreateParty Example methods for handling specific socket calls
func (socket *DefaultSocket) CreateParty(open bool, maxSize int) (*Party, error) {
	request := map[string]interface{}{
		"party_create": map[string]interface{}{
			"open":     open,
			"max_size": maxSize,
		},
	}

	err := socket.Send(request, DefaultSendTimeoutMs)
	if err != nil {
		return nil, err
	}

	return &Party{Open: open, MaxSize: maxSize}, nil
}

// PingPong does a periodic ping-pong check with the server.
func (socket *DefaultSocket) PingPong() {
	ticker := time.NewTicker(time.Duration(socket.HeartbeatTimeoutMs) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ping := map[string]interface{}{"ping": struct{}{}}
			if err := socket.Send(ping, socket.HeartbeatTimeoutMs); err != nil {
				if socket.Adapter.IsOpen() {
					socket.OnHeartbeatTimeout()
					socket.Adapter.Close()
				}
				return
			}
		}
	}
}

// OnHeartbeatTimeout handles heartbeat timeouts.
func (socket *DefaultSocket) OnHeartbeatTimeout() {
	if socket.Verbose {
		fmt.Println("Heartbeat timeout")
	}
}
