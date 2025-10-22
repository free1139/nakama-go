package nakama

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/eweb/jsonp"
	"github.com/gwaylib/log"
)

const (
	_msg_key_todo                = "todo" // need confirm the key
	_msg_key_channel             = "channel"
	_msg_key_match               = "match"
	_msg_key_party_join_request  = "party_join_request"
	_msg_key_channel_message_ack = "channel_message_ack"
	_msg_key_party_leader        = "party_leader"
	_msg_key_rpc                 = "rpc"
	_msg_key_ping                = "ping"
)

type Cid string

type RspResult struct {
	Decoded jsonp.Params // try parse, maybe nil
	Message []byte       // origin data
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

type ChannelRsp struct {
	Cid     `json:"cid"`
	Channel Channel `json:"channel"`
}

type ChannelJoin struct {
	Cid         `json:"cid"`
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

type ChannelMessageAckRsp struct {
	Cid               string            `json:"cid"`
	ChannelMessageAck ChannelMessageAck `json:"channel_message_ack"`
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

type MatchCreate struct {
	Name *string `json:"name,omitempty"`
}
type MatchReq struct {
	MatchCreate MatchCreate `json:"match_create"`
}
type MatchRsp struct {
	Cid   string `json:"match"`
	Match Match  `json:"match"`
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

type PartyLeaderRsp struct {
	Cid         string      `json:"cid"`
	PartyLeader PartyLeader `json:"party_leader"`
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

type PartyJoinRequestRsp struct {
	Cid              string           `json:"cid"`
	PartyJoinRequest PartyJoinRequest `party_join_request`
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

type StatusRsp struct {
	Cid    string `json:"cid"`
	Status Status `json:"status"`
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
type Socket interface {
	OnHeartbeatTimeout()
}

// SocketError represents an error received from a socket message.
type SocketError struct {
	Code    int    `json:"code"`    // The error code
	Message string `json:"message"` // A message in English to help developers debug the response
}

type Message struct {
	Cid           *string         `json:"cid"`
	Error         *error          `json:"error"`
	Notifications *[]Notification `json:"notifications"`
	Payload       interface{}     `json:"payload"`
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
	Adapter            *WebSocketAdapter
	SendTimeoutMs      int
	HeartbeatTimeoutMs int
	cIds               sync.Map // string:chan any
	nextCid            int

	userClosed     atomic.Bool
	reconnectTimes atomic.Int32
}

// NewDefaultSocket creates an instance of DefaultSocket.
func NewDefaultSocket(host, port string, useSSL, verbose bool, adapter *WebSocketAdapter, sendTimeoutMs *int) *DefaultSocket {
	if adapter == nil {
		adapter = NewWebSocketAdapterText()
	}
	if sendTimeoutMs == nil {
		defaultTimeout := DefaultSendTimeoutMs
		sendTimeoutMs = &defaultTimeout
	}

	return &DefaultSocket{
		Host:               host,
		Port:               port,
		UseSSL:             useSSL,
		Verbose:            verbose,
		Adapter:            adapter,
		SendTimeoutMs:      *sendTimeoutMs,
		HeartbeatTimeoutMs: DefaultHeartbeatTimeoutMs,
		cIds:               sync.Map{},
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
func (socket *DefaultSocket) Connect(session *Session, createStatus *bool, timeoutMs *int, userHandle func(*RspResult) bool) error {
	if createStatus == nil {
		defaultStatus := false
		createStatus = &defaultStatus
	}

	if timeoutMs == nil {
		defaultTimeout := DefaultConnectTimeoutMs
		timeoutMs = &defaultTimeout
	}

	scheme := "ws://"
	if socket.UseSSL {
		scheme = "wss://"
	}
	if !checkStr(session.Token) {
		return errors.New("Invalid token")
	}

	if err := socket.Adapter.Connect(scheme, socket.Host, socket.Port, *createStatus, *session.Token); err != nil {
		return errors.As(err)
	}

	socket.Adapter.onClose = socket.onDisconnect

	socket.Adapter.onError = socket.onError

	socket.Adapter.onMessage = func(mType int, message []byte) {
		if err := socket.handleMessage(mType, message, userHandle); err != nil {
			log.Warn(errors.As(err))
		}
	}

	socket.Adapter.onOpen = func(event interface{}) error {
		log.Printf("Socket opened: %v\n", event)
		go socket.pingPong(context.TODO())
		return nil
	}

	return nil

}

// Disconnect terminates the WebSocket connection.
func (socket *DefaultSocket) Disconnect(fireDisconnectEvent bool) {
	socket.userClosed.Store(true)
	if socket.Adapter.IsOpen() {
		socket.Adapter.Close()
	}
	if fireDisconnectEvent {
		socket.onDisconnect(fmt.Errorf("socket disconnected"))
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
func (socket *DefaultSocket) onDisconnect(evt error) {
	if socket.Verbose {
		log.Info("OnDisconnect:", evt)
	}
	if socket.userClosed.Load() {
		return
	}
	// TODO: try reconnect
}

// OnError handles WebSocket errors.
func (socket *DefaultSocket) onError(evt error) {
	if socket.Verbose {
		log.Info("OnError:", evt)
	}
	if socket.userClosed.Load() {
		return
	}
	// TODO: try reconnect
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

// HandleMessage processes incoming WebSocket messages.
func (socket *DefaultSocket) handleMessage(mType int, message []byte, handle func(*RspResult) bool) error {
	//log.Debugf("message_type:%d, message:%s", mType, string(message))

	result := &RspResult{Message: message}
	// try find the request cid
	decoded, err := jsonp.ParseParams(message)
	if err != nil {
		handle(result)
		return nil
	}
	result.Decoded = decoded

	// Handle specific decoding logic for match_data and party_data
	decodeReceivedData(decoded, "match_data")
	decodeReceivedData(decoded, "party_data")

	cid := decoded.String("cid")
	if len(cid) > 0 {
		rsp, ok := socket.cIds.Load(cid)
		if ok {
			rsp.(chan any) <- result
			return nil
		}
	}

	if val := decoded.Any(_msg_key_match); val != nil {
		rsp, exists := socket.cIds.Load(_msg_key_match)
		if exists {
			rsp.(chan any) <- result
			return nil
		}
	}

	if val := decoded.Any(_msg_key_channel); val != nil {
		rsp, exists := socket.cIds.Load(_msg_key_channel)
		if exists {
			rsp.(chan any) <- result
			return nil
		}
	}
	if val := decoded.Any(_msg_key_party_join_request); val != nil {
		rsp, exists := socket.cIds.Load(_msg_key_party_join_request)
		if exists {
			rsp.(chan any) <- result
			return nil
		}
	}
	if val := decoded.Any(_msg_key_channel_message_ack); val != nil {
		rsp, exists := socket.cIds.Load(_msg_key_channel_message_ack)
		if exists {
			rsp.(chan any) <- result
			return nil
		}
	}
	if val := decoded.Any(_msg_key_party_leader); val != nil {
		rsp, exists := socket.cIds.Load(_msg_key_channel_message_ack)
		if exists {
			rsp.(chan any) <- result
			return nil
		}
	}

	// unknow message, notify to caller
	handle(result)
	return nil

}

// Send sends a message to the WebSocket server with optional timeout.
// any should be error or []byte or Rsp pointer
func (socket *DefaultSocket) Send(mainKey, methodKey string, message map[string]any, sendTimeout *int) any {
	if !socket.Adapter.IsOpen() {
		return errors.New("socket connection is not established")
	}

	rsp := make(chan any, 1)
	defer close(rsp)

	cid := socket.GenerateCID()
	message["cid"] = cid // write a seq number

	socket.cIds.Store(cid, rsp)
	defer socket.cIds.Delete(cid)
	socket.cIds.Store(mainKey, rsp)
	defer socket.cIds.Delete(mainKey)

	//// Handle specific cases of match_data_send and party_data_send
	//if msgMap, ok := message.(map[string]interface{}); ok {
	//	handleEncodedData(msgMap, "match_data_send")
	//	handleEncodedData(msgMap, "party_data_send")
	//}

	if err := socket.Adapter.Send(message); err != nil {
		return errors.As(err)
	}

	if sendTimeout == nil {
		sendTimeout = new(int)
		*sendTimeout = DefaultTimeoutMs
	}

	t := time.NewTimer(time.Duration(*sendTimeout) * time.Millisecond)
	select {
	case <-t.C:
		return errors.New("timeout").As(mainKey, methodKey)
	case data := <-rsp:
		return data
	}
}

// CreateMatch sends a request to create a match and returns the created Match.
func (socket *DefaultSocket) CreateMatch(name *string) (*Match, error) {
	req := map[string]any{
		"match_create": MatchCreate{Name: name},
	}

	result := socket.Send("match", "CreateMatch", req, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}
	rsp, ok := result.(*RspResult)
	if !ok {
		return nil, errors.New("unknow protocal").As(result)
	}

	var match MatchRsp
	if err := json.Unmarshal(rsp.Message, &match); err != nil {
		return nil, fmt.Errorf("failed to deserialize match data into Match struct: %w", err)
	}

	return &match.Match, nil
}

// CreateParty Example methods for handling specific socket calls
func (socket *DefaultSocket) CreateParty(open bool, maxSize int) (*Party, error) {
	request := map[string]interface{}{
		"party_create": map[string]interface{}{
			"open":     open,
			"max_size": maxSize,
		},
	}

	result := socket.Send(_msg_key_todo, "CreateParty", request, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}

	return &Party{Open: open, MaxSize: maxSize}, nil
}

// FollowUsers sends a request to follow a list of user IDs and returns the status.
func (socket *DefaultSocket) FollowUsers(userIds []string) (*Status, error) {
	request := map[string]interface{}{
		"status_follow": map[string]interface{}{
			"user_ids": userIds,
		},
	}

	result := socket.Send(_msg_key_todo, "FollowUsers", request, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}

	var response StatusRsp
	if err := json.Unmarshal(result.(*RspResult).Message, &response); err != nil {
		return nil, errors.As(err)
	}
	return &response.Status, nil
}

// JoinChat sends a request to join a chat and returns the joined Channel.
func (socket *DefaultSocket) JoinChat(target string, chatType int, persistence, hidden bool) (*Channel, error) {
	request := map[string]interface{}{
		"channel_join": map[string]interface{}{
			"target":      target,
			"type":        chatType,
			"persistence": persistence,
			"hidden":      hidden,
		},
	}

	result := socket.Send(_msg_key_channel, "JoinChat", request, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}

	// TODO: some other server push when sent
	data := ChannelRsp{}
	if err := json.Unmarshal(result.(*RspResult).Message, &data); err != nil {
		return nil, errors.As(err, result)
	}
	return &data.Channel, nil
}

// JoinMatch sends a request to join a match and returns the joined Match.
func (socket *DefaultSocket) JoinMatch(matchID, token *string, metadata *map[string]interface{}) (*Match, error) {
	request := map[string]interface{}{
		"match_join": map[string]interface{}{
			"metadata": metadata,
		},
	}

	if token != nil && *token != "" {
		request["match_join"].(map[string]interface{})["token"] = token
	} else {
		request["match_join"].(map[string]interface{})["match_id"] = matchID
	}

	result := socket.Send(_msg_key_match, "JoinMatch", request, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}

	matchRsp := &MatchRsp{}
	if err := json.Unmarshal(result.(*RspResult).Message, matchRsp); err != nil {
		return nil, errors.As(err)
	}
	return &matchRsp.Match, nil
}

// JoinParty sends a request to join a party.
func (socket *DefaultSocket) JoinParty(partyID string) error {
	request := map[string]interface{}{
		"party_join": map[string]interface{}{
			"party_id": partyID,
		},
	}

	result := socket.Send(_msg_key_todo, "JoinPatry", request, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	// TODO: need decode?

	return nil
}

// LeaveChat sends a request to leave a chat channel.
func (socket *DefaultSocket) LeaveChat(channelID string) error {
	request := map[string]interface{}{
		"channel_leave": map[string]interface{}{
			"channel_id": channelID,
		},
	}

	result := socket.Send(_msg_key_todo, "LeaveChat", request, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	// TODO: decode?

	return nil
}

// LeaveMatch sends a request to leave a match.
func (socket *DefaultSocket) LeaveMatch(matchID string) error {
	request := map[string]interface{}{
		"match_leave": map[string]interface{}{
			"match_id": matchID,
		},
	}

	result := socket.Send(_msg_key_todo, "LeaveMatch", request, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	// TODO: decode?

	return nil
}

// LeaveParty sends a request to leave a party.
func (socket *DefaultSocket) LeaveParty(partyID string) error {
	request := map[string]interface{}{
		"party_leave": map[string]interface{}{
			"party_id": partyID,
		},
	}

	result := socket.Send(_msg_key_todo, "LeaveParty", request, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	// TODO: decode

	return nil
}

// ListPartyJoinRequests fetches the list of join requests for a given party ID.
func (socket *DefaultSocket) ListPartyJoinRequests(partyID string) (*PartyJoinRequest, error) {
	request := map[string]interface{}{
		"party_join_request_list": map[string]interface{}{
			"party_id": partyID,
		},
	}

	// TODO: confirm the main key is channel.
	result := socket.Send(_msg_key_party_join_request, "ListPartyJoinRequests", request, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}

	response := &PartyJoinRequestRsp{}
	if err := json.Unmarshal(result.(*RspResult).Message, &response); err != nil {
		return nil, errors.As(err)
	}

	return &response.PartyJoinRequest, nil
}

// RemoveChatMessage sends a request to remove a chat message and returns the ChannelMessageAck.
func (socket *DefaultSocket) RemoveChatMessage(channelID, messageID string) (*ChannelMessageAck, error) {
	request := map[string]interface{}{
		"channel_message_remove": map[string]interface{}{
			"channel_id": channelID,
			"message_id": messageID,
		},
	}

	result := socket.Send(_msg_key_channel_message_ack, "RemoveChatMessage", request, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}

	response := &ChannelMessageAckRsp{}
	if err := json.Unmarshal(result.(*RspResult).Message, response); err != nil {
		return nil, errors.As(err, result)
	}
	return &response.ChannelMessageAck, nil
}

// PromotePartyMember promotes a party member to party leader and returns the new PartyLeader.
func (socket *DefaultSocket) PromotePartyMember(partyID string, partyMember Presence) (*PartyLeader, error) {
	request := map[string]interface{}{
		"party_promote": map[string]interface{}{
			"party_id": partyID,
			"presence": partyMember,
		},
	}

	result := socket.Send(_msg_key_party_leader, "PromotePartyMember", request, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}

	response := &PartyLeaderRsp{}
	if err := json.Unmarshal(result.(*RspResult).Message, &response); err != nil {
		return nil, errors.As(err, result)
	}
	return &response.PartyLeader, nil
}

// RemoveMatchmaker sends a request to remove a matchmaker ticket.
func (socket *DefaultSocket) RemoveMatchmaker(ticket string) error {
	request := map[string]interface{}{
		"matchmaker_remove": map[string]interface{}{
			"ticket": ticket,
		},
	}

	result := socket.Send(_msg_key_todo, "RemoveMatchmaker", request, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	return nil
}

// RemoveMatchmakerParty sends a request to remove a matchmaker ticket from a party.
func (socket *DefaultSocket) RemoveMatchmakerParty(partyID, ticket string) error {
	request := map[string]interface{}{
		"party_matchmaker_remove": map[string]interface{}{
			"party_id": partyID,
			"ticket":   ticket,
		},
	}

	result := socket.Send(_msg_key_todo, "RemoveMatchmakerParty", request, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	return nil
}

// RemovePartyMember sends a request to remove a member from a party.
func (socket *DefaultSocket) RemovePartyMember(partyID string, member Presence) error {
	request := map[string]interface{}{
		"party_remove": map[string]interface{}{
			"party_id": partyID,
			"presence": member,
		},
	}

	result := socket.Send(_msg_key_todo, "RemovePartyMemeber", request, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	return nil
}

// Rpc sends an RPC request and returns an ApiRpc response.
func (socket *DefaultSocket) Rpc(id, payload, httpKey string) (*ApiRpc, error) {
	request := map[string]interface{}{
		"rpc": map[string]interface{}{
			"id":       id,
			"payload":  payload,
			"http_key": httpKey,
		},
	}

	result := socket.Send(_msg_key_rpc, "Rpc", request, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}

	response := &ApiRpcRsp{}
	if err := json.Unmarshal(result.(*RspResult).Message, response); err != nil {
		return nil, errors.As(err, result)
	}
	return &response.ApiRpc, nil
}

// SendMatchState sends match state updates to the server.
func (socket *DefaultSocket) SendMatchState(matchID string, opCode int, data interface{}, presences []Presence, reliable bool) error {
	request := map[string]interface{}{
		"match_data_send": map[string]interface{}{
			"match_id":  matchID,
			"op_code":   opCode,
			"data":      data,
			"presences": presences,
			"reliable":  reliable,
		},
	}

	// TODO: confirm the msg_key
	result := socket.Send(_msg_key_match, "SendMatchState", request, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	return nil
}

// SendPartyData sends party data updates to the server.
func (socket *DefaultSocket) SendPartyData(partyID string, opCode int, data interface{}) error {
	request := map[string]interface{}{
		"party_data_send": map[string]interface{}{
			"party_id": partyID,
			"op_code":  opCode,
			"data":     data,
		},
	}

	result := socket.Send(_msg_key_todo, "SendPartyData", request, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	return nil
}

// UnfollowUsers sends a request to unfollow the specified users.
func (socket *DefaultSocket) UnfollowUsers(userIDs []string) error {
	request := map[string]interface{}{
		"status_unfollow": map[string]interface{}{
			"user_ids": userIDs,
		},
	}

	result := socket.Send(_msg_key_todo, "UnfollowUsers", request, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	return nil
}

// UpdateChatMessage sends a request to update a chat message and returns the ChannelMessageAck.
func (socket *DefaultSocket) UpdateChatMessage(channelID, messageID string, content interface{}) (*ChannelMessageAck, error) {
	request := map[string]interface{}{
		"channel_message_update": map[string]interface{}{
			"channel_id": channelID,
			"message_id": messageID,
			"content":    content,
		},
	}

	result := socket.Send(_msg_key_channel_message_ack, "UpdateChatMessage", request, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}

	response := &ChannelMessageAckRsp{}
	if err := json.Unmarshal(result.(*RspResult).Message, &response); err != nil {
		return nil, errors.As(err)
	}
	return &response.ChannelMessageAck, nil
}

// UpdateStatus sends a status update to the server.
func (socket *DefaultSocket) UpdateStatus(status *string) error {
	request := map[string]interface{}{
		"status_update": map[string]interface{}{
			"status": status,
		},
	}

	result := socket.Send(_msg_key_todo, "UpdateStatus", request, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	return nil
}

// WriteChatMessage sends a chat message and returns the ChannelMessageAck.
func (socket *DefaultSocket) WriteChatMessage(channelID string, content interface{}) (*ChannelMessageAck, error) {
	// const response = await this.send({channel_message_send: {channel_id: channel_id, content: content}});
	request := map[string]interface{}{
		"channel_message_send": map[string]interface{}{
			"channel_id": channelID,
			"content":    content,
		},
	}

	result := socket.Send(_msg_key_channel_message_ack, "WriteChatMessage", request, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}

	response := &ChannelMessageAckRsp{}
	if err := json.Unmarshal(result.(*RspResult).Message, &response); err != nil {
		return nil, errors.As(err)
	}
	return &response.ChannelMessageAck, nil
}

// pingPong does a periodic ping-pong check with the server.
func (socket *DefaultSocket) pingPong(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(socket.HeartbeatTimeoutMs) * time.Millisecond)
	defer ticker.Stop()
	log.Println("before pingpong socket is nil:", socket.Adapter.socket == nil)

	for {
		select {
		case <-ticker.C:
			ping := map[string]interface{}{"ping": struct{}{}}
			result := socket.Send(_msg_key_ping, "pingPong", ping, &socket.HeartbeatTimeoutMs)
			if err, ok := result.(error); ok {
				log.Println("after pingpong socket is nil:", socket.Adapter.socket == nil)
				log.Println("Failed to send ping:", err)
				if socket.Adapter.IsOpen() {
					socket.OnHeartbeatTimeout()
					socket.Adapter.Close()
				}
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// OnHeartbeatTimeout handles heartbeat timeouts.
func (socket *DefaultSocket) OnHeartbeatTimeout() {
	if socket.Verbose {
		fmt.Println("Heartbeat timeout")
	}
}
