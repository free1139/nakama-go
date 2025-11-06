package nakama

import (
	"context"
	"encoding/base64"
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gwaylib/errors"
	"github.com/gwaylib/log"
	api "github.com/heroiclabs/nakama-common/api"
	"github.com/heroiclabs/nakama-common/rtapi"
	"google.golang.org/protobuf/encoding/protojson"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

const (
	ChannelMessageTypeChat int32 = iota
	ChannelMessageTypeChatUpdate
	ChannelMessageTypeChatRemove
	ChannelMessageTypeGroupJoin
	ChannelMessageTypeGroupAdd
	ChannelMessageTypeGroupLeave
	ChannelMessageTypeGroupKick
	ChannelMessageTypeGroupPromote
	ChannelMessageTypeGroupBan
	ChannelMessageTypeGroupDemote
)

const (
	EventTypeConnect   = 0
	EventTypeMessage   = 1
	EventTypeReconnect = 2
	EventTypeConnected = 3
	EventTypePingPong  = 4
)

type RspResult struct {
	Decoded *rtapi.Envelope // try parse, maybe nil
	Message any             // origin data
}

type EventHandler func(event int, data *RspResult)

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
	verbose            bool
	adapter            *WebSocketAdapter
	sendTimeoutMs      int
	heartbeatTimeoutMs int
	eventHandle        EventHandler

	cIds    sync.Map // string:chan any
	nextCid int

	pingUsed      time.Duration
	userClosed    atomic.Bool
	joinChatStack sync.Map
}

// NewDefaultSocket creates an instance of DefaultSocket.
func NewDefaultSocket(eventHandle EventHandler, host, port, token string, useSSL, verbose bool, sendTimeoutMs *int, createStatus *bool) *DefaultSocket {
	if sendTimeoutMs == nil {
		defaultTimeout := DefaultSendTimeoutMs
		sendTimeoutMs = &defaultTimeout
	}

	if createStatus == nil {
		defaultStatus := false
		createStatus = &defaultStatus
	}

	scheme := "ws://"
	if useSSL {
		scheme = "wss://"
	}

	socket := &DefaultSocket{
		verbose:            verbose,
		sendTimeoutMs:      *sendTimeoutMs,
		heartbeatTimeoutMs: DefaultHeartbeatTimeoutMs,
		eventHandle:        eventHandle,
		cIds:               sync.Map{},
		nextCid:            1,
	}
	adapter := NewWebSocketAdapterText(scheme, host, port, *createStatus, token)
	adapter.onError = socket.onError
	adapter.onMessage = func(mType int, message []byte) {
		if err := socket.handleMessage(mType, message); err != nil {
			log.Warn(errors.As(err))
		}
	}
	socket.adapter = adapter
	return socket
}

// GenerateCID generates a unique client ID for requests.
func (socket *DefaultSocket) GenerateCID() string {
	cid := fmt.Sprintf("%d", socket.nextCid)
	socket.nextCid++
	return cid
}

// Connect establishes the WebSocket connection with optional timeouts.
func (socket *DefaultSocket) Connect() error {
	socket.eventHandle(EventTypeConnect, nil)
	if err := socket.adapter.Connect(); err != nil {
		return errors.As(err)
	}
	go socket.pingPong(context.TODO())

	socket.eventHandle(EventTypeConnected, nil)
	return nil

}

// Disconnect terminates the WebSocket connection.
func (socket *DefaultSocket) Disconnect() {
	socket.userClosed.Store(true)
	if socket.adapter.IsOpen() {
		socket.adapter.Close()
	}
}

// SetHeartbeatTimeoutMs sets the timeout for heartbeat pings.
func (socket *DefaultSocket) SetHeartbeatTimeoutMs(ms int) {
	socket.heartbeatTimeoutMs = ms
}

// GetHeartbeatTimeoutMs gets the timeout for heartbeat pings.
func (socket *DefaultSocket) GetHeartbeatTimeoutMs() int {
	return socket.heartbeatTimeoutMs
}

func (socket *DefaultSocket) reconnect(tryTimes int) error {
	if socket.eventHandle != nil {
		socket.eventHandle(EventTypeReconnect, nil)
	}
	for i := tryTimes; i > 0; i-- {
		if socket.userClosed.Load() {
			return errors.New("user has closed the connection")
		}
		if socket.adapter.IsOpen() {
			return nil
		}

		if err := socket.adapter.Connect(); err != nil {
			log.Warn("retry failed", errors.As(err, i))
			time.Sleep(3e9)
			continue
		}

		// restore the chats
		socket.joinChatStack.Range(func(k, v any) bool {
			log.Debugf("reconnect talk, key:%+v", k)
			originJoin := v.(*rtapi.ChannelJoin)
			rejoinChannel, err := socket.joinChat(originJoin)
			if err != nil {
				log.Warn(errors.As(err))
				return true // continue range
			}
			if rejoinChannel.Id != k.(string) {
				log.Error(errors.New("Auto reconnect the channel not match").As(k, rejoinChannel.Id))
			}
			return true // continue range
		})
		return nil
	}
	return errors.New("reconnection failed")
}

// OnError handles WebSocket errors.
func (socket *DefaultSocket) onError(evt error) {
	if socket.verbose {
		log.Info("OnError:", evt)
	}
	socket.reconnect(math.MaxInt)
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
func (socket *DefaultSocket) handleMessage(mType int, message []byte) error {
	//log.Debugf("message_type:%d, message:%s", mType, string(message))

	result := &RspResult{Message: message}
	// try find the request cid
	decoded := &rtapi.Envelope{}
	if err := protojson.Unmarshal(message, decoded); err != nil {
		if socket.eventHandle != nil {
			socket.eventHandle(EventTypeMessage, result)
			return nil
		}
		return errors.As(err)
	}
	result.Decoded = decoded

	// Handle specific decoding logic for match_data and party_data
	// decodeReceivedData(decoded, "match_data")
	//decodeReceivedData(decoded, "party_data")

	cid := decoded.Cid
	if len(cid) > 0 {
		rsp, ok := socket.cIds.Load(cid)
		if ok {
			err, ok := decoded.GetMessage().(*rtapi.Envelope_Error)
			if ok {
				rsp.(chan any) <- errors.New(err.Error.Message).As(err.Error)
			} else {
				rsp.(chan any) <- result
			}

			return nil
		}
	}

	// deal the kick channel event and more.
	if event, ok := decoded.GetMessage().(*rtapi.Envelope_ChannelMessage); ok {
		// {"channel_message":{"channel_id":"3.4f634582-8cd0-4fd8-a71c-3e093ae30cf2..", "message_id":"1bdac70c-9d0e-4c6b-ad1d-bd969f95e1bf", "code":6, "sender_id":"f5996e0c-37da-421f-a7e6-df78eb4c79ad", "username":"z1", "content":"{}", "create_time":"2025-10-24T09:00:46.033690893Z", "update_time":"2025-10-24T09:00:46.033690893Z", "persistent":true, "group_id":"4f634582-8cd0-4fd8-a71c-3e093ae30cf2"}}
		msg := event.ChannelMessage
		if msg.Code != nil && msg.Code.Value == ChannelMessageTypeGroupKick {
			socket.joinChatStack.Delete(msg.ChannelId)
		}
	}

	// unknow message, notify to caller
	if socket.eventHandle != nil {
		socket.eventHandle(EventTypeMessage, result)
	}
	return nil

}

// Send sends a message to the WebSocket server with optional timeout.
// any should be error or []byte or Rsp pointer
func (socket *DefaultSocket) Send(message *rtapi.Envelope, sendTimeout *int) any {
	if !socket.adapter.IsOpen() {
		if err := socket.reconnect(3); err != nil {
			return errors.As(err)
		}
	}

	rsp := make(chan any, 1)
	defer close(rsp)

	cid := socket.GenerateCID()
	message.Cid = cid // write a seq number

	socket.cIds.Store(cid, rsp)
	defer socket.cIds.Delete(cid)

	//// Handle specific cases of match_data_send and party_data_send
	//if msgMap, ok := message.(map[string]interface{}); ok {
	//	handleEncodedData(msgMap, "match_data_send")
	//	handleEncodedData(msgMap, "party_data_send")
	//}

	if err := socket.adapter.Send(message); err != nil {
		return errors.As(err)
	}

	if sendTimeout == nil {
		sendTimeout = new(int)
		*sendTimeout = DefaultTimeoutMs
	}

	t := time.NewTimer(time.Duration(*sendTimeout) * time.Millisecond)
	select {
	case <-t.C:
		return errors.New("timeout")
	case data := <-rsp:
		return data
	}
}

// CreateMatch sends a request to create a match and returns the created Match.
func (socket *DefaultSocket) CreateMatch(name *string) (*rtapi.Match, error) {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchCreate{
			MatchCreate: &rtapi.MatchCreate{Name: *name},
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}
	rsp, ok := result.(*RspResult)
	if !ok {
		return nil, errors.New("unknow protocal").As(result)
	}

	return rsp.Decoded.GetMessage().(*rtapi.Envelope_Match).Match, nil
}

// CreateParty Example methods for handling specific socket calls
func (socket *DefaultSocket) CreateParty(open bool, maxSize int32) (*rtapi.Party, error) {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyCreate{
			PartyCreate: &rtapi.PartyCreate{Open: open, MaxSize: maxSize},
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}

	return result.(*RspResult).Decoded.GetMessage().(*rtapi.Envelope_Party).Party, nil
}

// FollowUsers sends a request to follow a list of user IDs and returns the status.
func (socket *DefaultSocket) FollowUsers(userIds []string) (*rtapi.Status, error) {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_StatusFollow{
			StatusFollow: &rtapi.StatusFollow{
				UserIds: userIds,
			},
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}

	return result.(*RspResult).Decoded.GetMessage().(*rtapi.Envelope_Status).Status, nil
}

func (socket *DefaultSocket) joinChat(target *rtapi.ChannelJoin) (*rtapi.Channel, error) {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_ChannelJoin{
			ChannelJoin: target,
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}
	return result.(*RspResult).Decoded.GetMessage().(*rtapi.Envelope_Channel).Channel, nil
}

// JoinChat sends a request to join a chat and returns the joined Channel.
func (socket *DefaultSocket) JoinChat(target string, chatType int32, persistence, hidden bool) (*rtapi.Channel, error) {
	targetChannel := &rtapi.ChannelJoin{
		Target:      target,
		Type:        chatType,
		Persistence: wrapperspb.Bool(persistence),
		Hidden:      wrapperspb.Bool(hidden),
	}
	channel, err := socket.joinChat(targetChannel)
	if err != nil {
		return nil, errors.As(err)
	}
	// stash the the chat, and restore when socket reconnected.
	socket.joinChatStack.Store(channel.Id, targetChannel)
	return channel, nil
}

// JoinMatch sends a request to join a match and returns the joined Match.
func (socket *DefaultSocket) JoinMatch(matchID, token *string, metadata map[string]string) (*rtapi.Match, error) {
	matchId := &rtapi.MatchJoin_MatchId{MatchId: *matchID}
	matchToken := &rtapi.MatchJoin_Token{Token: *token}
	matchJoin := &rtapi.MatchJoin{
		Metadata: metadata,
	}
	if token != nil && *token != "" {
		matchJoin.Id = matchToken
	} else {
		matchJoin.Id = matchId
	}
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchJoin{
			MatchJoin: matchJoin,
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}

	return result.(*RspResult).Decoded.GetMessage().(*rtapi.Envelope_Match).Match, nil
}

// JoinParty sends a request to join a party.
func (socket *DefaultSocket) JoinParty(partyID string) error {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyJoin{
			PartyJoin: &rtapi.PartyJoin{
				PartyId: partyID,
			},
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	// TODO: need response?

	return nil
}

// LeaveChat sends a request to leave a chat channel.
func (socket *DefaultSocket) LeaveChat(channelID string) error {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_ChannelLeave{
			ChannelLeave: &rtapi.ChannelLeave{
				ChannelId: channelID,
			},
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	// TODO: decode?

	socket.joinChatStack.Delete(channelID)
	return nil
}

// LeaveMatch sends a request to leave a match.
func (socket *DefaultSocket) LeaveMatch(matchID string) error {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchLeave{
			MatchLeave: &rtapi.MatchLeave{
				MatchId: matchID,
			},
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	// TODO: decode?

	return nil
}

// LeaveParty sends a request to leave a party.
func (socket *DefaultSocket) LeaveParty(partyID string) error {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyLeave{
			PartyLeave: &rtapi.PartyLeave{
				PartyId: partyID,
			},
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	// TODO: decode

	return nil
}

// ListPartyJoinRequests fetches the list of join requests for a given party ID.
func (socket *DefaultSocket) ListPartyJoinRequests(partyID string) (*rtapi.PartyJoinRequest, error) {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyJoinRequestList{
			PartyJoinRequestList: &rtapi.PartyJoinRequestList{
				PartyId: partyID,
			},
		},
	}

	// TODO: confirm the main key is channel.
	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}

	return result.(*RspResult).Decoded.GetMessage().(*rtapi.Envelope_PartyJoinRequest).PartyJoinRequest, nil
}

// RemoveChatMessage sends a request to remove a chat message and returns the ChannelMessageAck.
func (socket *DefaultSocket) RemoveChatMessage(channelID, messageID string) (*rtapi.ChannelMessageAck, error) {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_ChannelMessageRemove{
			ChannelMessageRemove: &rtapi.ChannelMessageRemove{
				ChannelId: channelID,
				MessageId: messageID,
			},
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}
	return result.(*RspResult).Decoded.GetMessage().(*rtapi.Envelope_ChannelMessageAck).ChannelMessageAck, nil
}

// PromotePartyMember promotes a party member to party leader and returns the new PartyLeader.
func (socket *DefaultSocket) PromotePartyMember(partyID string, partyMember *rtapi.UserPresence) (*rtapi.PartyLeader, error) {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyPromote{
			PartyPromote: &rtapi.PartyPromote{
				PartyId:  partyID,
				Presence: partyMember,
			},
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}
	return result.(*RspResult).Decoded.GetMessage().(*rtapi.Envelope_PartyLeader).PartyLeader, nil
}

// RemoveMatchmaker sends a request to remove a matchmaker ticket.
func (socket *DefaultSocket) RemoveMatchmaker(ticket string) error {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchmakerRemove{
			MatchmakerRemove: &rtapi.MatchmakerRemove{
				Ticket: ticket,
			},
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	return nil
}

// RemoveMatchmakerParty sends a request to remove a matchmaker ticket from a party.
func (socket *DefaultSocket) RemoveMatchmakerParty(partyID, ticket string) error {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyMatchmakerRemove{
			PartyMatchmakerRemove: &rtapi.PartyMatchmakerRemove{
				PartyId: partyID,
				Ticket:  ticket,
			},
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	return nil
}

// RemovePartyMember sends a request to remove a member from a party.
func (socket *DefaultSocket) RemovePartyMember(partyID string, member *rtapi.UserPresence) error {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyRemove{
			PartyRemove: &rtapi.PartyRemove{
				PartyId:  partyID,
				Presence: member,
			},
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	return nil
}

// Rpc sends an RPC request and returns an ApiRpc response.
func (socket *DefaultSocket) Rpc(id, payload, httpKey string) (*api.Rpc, error) {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_Rpc{
			Rpc: &api.Rpc{
				Id:      id,
				Payload: payload,
				HttpKey: httpKey,
			},
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}
	return result.(*RspResult).Decoded.GetMessage().(*rtapi.Envelope_Rpc).Rpc, nil
}

// SendMatchState sends match state updates to the server.
func (socket *DefaultSocket) SendMatchState(matchID string, opCode int64, data []byte, presences []*rtapi.UserPresence, reliable bool) error {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_MatchDataSend{
			MatchDataSend: &rtapi.MatchDataSend{
				MatchId:   matchID,
				OpCode:    opCode,
				Data:      data,
				Presences: presences,
				Reliable:  reliable,
			},
		},
	}

	// TODO: confirm the msg_key
	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	return nil
}

// SendPartyData sends party data updates to the server.
func (socket *DefaultSocket) SendPartyData(partyID string, opCode int64, data []byte) error {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_PartyDataSend{
			PartyDataSend: &rtapi.PartyDataSend{
				PartyId: partyID,
				OpCode:  opCode,
				Data:    data,
			},
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	return nil
}

// UnfollowUsers sends a request to unfollow the specified users.
func (socket *DefaultSocket) UnfollowUsers(userIDs []string) error {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_StatusUnfollow{
			StatusUnfollow: &rtapi.StatusUnfollow{
				UserIds: userIDs,
			},
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	return nil
}

// UpdateChatMessage sends a request to update a chat message and returns the ChannelMessageAck.
func (socket *DefaultSocket) UpdateChatMessage(channelID, messageID string, content string) (*rtapi.ChannelMessageAck, error) {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_ChannelMessageUpdate{
			ChannelMessageUpdate: &rtapi.ChannelMessageUpdate{
				ChannelId: channelID,
				MessageId: messageID,
				Content:   content,
			},
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}

	return result.(*RspResult).Decoded.GetMessage().(*rtapi.Envelope_ChannelMessageAck).ChannelMessageAck, nil
}

// UpdateStatus sends a status update to the server.
func (socket *DefaultSocket) UpdateStatus(status *string) error {
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_StatusUpdate{
			StatusUpdate: &rtapi.StatusUpdate{
				Status: wrapperspb.String(*status),
			},
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return errors.As(err)
	}

	return nil
}

// WriteChatMessage sends a chat message and returns the ChannelMessageAck.
func (socket *DefaultSocket) WriteChatMessage(channelID, content string) (*rtapi.ChannelMessageAck, error) {
	// const response = await this.send({channel_message_send: {channel_id: channel_id, content: content}});
	req := &rtapi.Envelope{
		Message: &rtapi.Envelope_ChannelMessageSend{
			ChannelMessageSend: &rtapi.ChannelMessageSend{
				ChannelId: channelID,
				Content:   content,
			},
		},
	}

	result := socket.Send(req, nil)
	if err, ok := result.(error); ok {
		return nil, errors.As(err)
	}
	return result.(*RspResult).Decoded.GetMessage().(*rtapi.Envelope_ChannelMessageAck).ChannelMessageAck, nil
}

// pingPong does a periodic ping-pong check with the server.
func (socket *DefaultSocket) pingPong(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(socket.heartbeatTimeoutMs) * time.Millisecond)
	defer ticker.Stop()
	log.Println("before pingpong socket is nil:", socket.adapter.socket == nil)

	pingReq := &rtapi.Envelope{
		Message: &rtapi.Envelope_Ping{
			Ping: &rtapi.Ping{},
		},
	}

	for {
		select {
		case <-ticker.C:
			starTime := time.Now()
			result := socket.Send(pingReq, &socket.heartbeatTimeoutMs)
			if err, ok := result.(error); ok {
				log.Println("Failed to send ping:", err)
				continue
			}
			socket.pingUsed = time.Now().Sub(starTime)
			socket.eventHandle(EventTypePingPong, &RspResult{Message: socket.pingUsed})
		case <-ctx.Done():
			return
		}
	}
}

// OnHeartbeatTimeout handles heartbeat timeouts.
func (socket *DefaultSocket) OnHeartbeatTimeout() {
	if socket.verbose {
		fmt.Println("Heartbeat timeout")
	}
}
