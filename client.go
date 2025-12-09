package nakama

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gwaylib/errors"
	api "github.com/heroiclabs/nakama-common/api"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
)

// Default configuration values
const (
	DefaultHost              = "127.0.0.1"
	DefaultPort              = "7350"
	DefaultServerKey         = "defaultkey"
	DefaultTimeoutMs         = 7000
	DefaultExpiredTimespanMs = 5 * 60 * 1000 // 5 minutes in milliseconds
)

// Client represents a client for the Nakama server.
type Client struct {
	ExpiredTimespanMs  int64      // The expired timespan used to check session lifetime.
	ApiClient          *NakamaApi // The low-level API client for Nakama server.
	ServerKey          string
	Host               string
	Port               string
	UseSSL             bool
	Timeout            int
	AutoRefreshSession bool
}

// NewClient creates a new instance of Client with the specified configuration.
func NewClient(
	serverKey string,
	host string,
	port string,
	useSSL bool,
	timeout int,
	autoRefreshSession bool,
) *Client {
	// Default values if not provided
	if serverKey == "" {
		serverKey = DefaultServerKey
	}
	if host == "" {
		host = DefaultHost
	}
	if port == "" {
		port = DefaultPort
	}
	if timeout == 0 {
		timeout = DefaultTimeoutMs
	}

	scheme := "http://"
	if useSSL {
		scheme = "https://"
	}
	basePath := scheme + host + ":" + port

	return &Client{
		ExpiredTimespanMs:  DefaultExpiredTimespanMs,
		ApiClient:          &NakamaApi{serverKey, basePath, timeout},
		ServerKey:          serverKey,
		Host:               host,
		Port:               port,
		UseSSL:             useSSL,
		Timeout:            timeout,
		AutoRefreshSession: autoRefreshSession,
	}
}

func (c *Client) refreshSession(session *Session) error {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().UnixMilli()+c.ExpiredTimespanMs)/1000) {
		if _, err := c.SessionRefresh(session, nil); err != nil {
			return errors.As(err)
		}
	}
	return nil
}

// AddGroupUsers adds users to a group, or accepts their join requests.
func (c *Client) AddGroupUsers(session *Session, groupId *string, ids []string) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.AddGroupUsers(&session.Token, groupId, ids, make(map[string]string))
}

// AddFriends adds friends by ID or username to a user's account.
func (c *Client) AddFriends(session *Session, ids []string, usernames []string) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.AddFriends(&session.Token, ids, usernames, make(map[string]string))
}

// AuthenticateApple authenticates a user with an Apple ID against the server.
func (c *Client) AuthenticateApple(token string, create *bool, username string, vars map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := &api.AccountApple{
		Token: token,
		Vars:  vars,
	}

	// Call the API client to authenticate with Apple
	apiSession, err := c.ApiClient.AuthenticateApple(c.ServerKey, "", request, create, username, make(map[string]string))
	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        apiSession.Token,
		RefreshToken: apiSession.RefreshToken,
		Created:      apiSession.Created,
	}, nil
}

// AuthenticateCustom authenticates a user with a custom ID against the server.
func (c *Client) AuthenticateCustom(id string, create *bool, username *string, vars map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := &api.AccountCustom{
		Id:   id,
		Vars: vars,
	}

	// Call the API client to authenticate with a custom ID
	apiSession, err := c.ApiClient.AuthenticateCustom(c.ServerKey, "", request, create, username, make(map[string]string))
	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        apiSession.Token,
		RefreshToken: apiSession.RefreshToken,
		Created:      apiSession.Created,
	}, nil
}

// AuthenticateDevice authenticates a user with a device ID against the server.
func (c *Client) AuthenticateDevice(id string, create *bool, username string, vars map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := &api.AccountDevice{
		Id:   id,
		Vars: vars,
	}

	// Call the API client to authenticate with a device ID
	apiSession, err := c.ApiClient.AuthenticateDevice(c.ServerKey, "", request, create, username, make(map[string]string))
	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        apiSession.Token,
		RefreshToken: apiSession.RefreshToken,
		Created:      apiSession.Created,
	}, nil
}

// AuthenticateEmail authenticates a user with an email and password against the server.
func (c *Client) AuthenticateEmail(email string, password string, create *bool, username *string, vars map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := &api.AccountEmail{
		Email:    email,
		Password: password,
		Vars:     vars,
	}

	// Call the API client to authenticate with email and password
	apiSession, err := c.ApiClient.AuthenticateEmail(c.ServerKey, "", request, create, username, make(map[string]string))
	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        apiSession.Token,
		RefreshToken: apiSession.RefreshToken,
		Created:      apiSession.Created,
	}, nil
}

// AuthenticateFacebookInstantGame authenticates a user with a Facebook Instant Game token against the server.
func (c *Client) AuthenticateFacebookInstantGame(signedPlayerInfo string, create *bool, username string, vars map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := &api.AccountFacebookInstantGame{
		SignedPlayerInfo: signedPlayerInfo,
		Vars:             vars,
	}

	// Call the API client to authenticate with Facebook Instant Game
	apiSession, err := c.ApiClient.AuthenticateFacebookInstantGame(c.ServerKey, "", request, create, username, make(map[string]string))
	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        apiSession.Token,
		RefreshToken: apiSession.RefreshToken,
		Created:      apiSession.Created,
	}, nil
}

// AuthenticateFacebook authenticates a user with a Facebook OAuth token against the server.
func (c *Client) AuthenticateFacebook(token string, create *bool, username string, sync *bool, vars map[string]string, options map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := &api.AccountFacebook{
		Token: token,
		Vars:  vars,
	}

	// Call the API client to authenticate with Facebook
	apiSession, err := c.ApiClient.AuthenticateFacebook(c.ServerKey, "", request, create, username, sync, options)
	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        apiSession.Token,
		RefreshToken: apiSession.RefreshToken,
		Created:      apiSession.Created,
	}, nil
}

// AuthenticateGoogle authenticates a user with a Google token against the server.
func (c *Client) AuthenticateGoogle(token string, create *bool, username string, vars map[string]string, options map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := &api.AccountGoogle{
		Token: token,
		Vars:  vars,
	}

	// Call the API client to authenticate with Google
	apiSession, err := c.ApiClient.AuthenticateGoogle(c.ServerKey, "", request, create, username, options)
	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        apiSession.Token,
		RefreshToken: apiSession.RefreshToken,
		Created:      apiSession.Created,
	}, nil
}

// AuthenticateGameCenter authenticates a user with GameCenter against the server.
func (c *Client) AuthenticateGameCenter(bundleId string, playerId string, publicKeyUrl string, salt string, signature string, timestamp int64, create *bool, username string, vars map[string]string, options map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := &api.AccountGameCenter{
		BundleId:         bundleId,
		PlayerId:         playerId,
		PublicKeyUrl:     publicKeyUrl,
		Salt:             salt,
		Signature:        signature,
		TimestampSeconds: timestamp,
		Vars:             vars,
	}

	// Call the API client to authenticate with GameCenter
	apiSession, err := c.ApiClient.AuthenticateGameCenter(c.ServerKey, "", request, create, username, options)
	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        apiSession.Token,
		RefreshToken: apiSession.RefreshToken,
		Created:      apiSession.Created,
	}, nil
}

// AuthenticateSteam authenticates a user with a Steam token against the server.
func (c *Client) AuthenticateSteam(token string, create *bool, username string, vars map[string]string) (*Session, error) {
	// Prepare the authentication request
	request := &api.AccountSteam{
		Token: token,
		Vars:  vars,
	}

	// Call the API client to authenticate with Steam
	apiSession, err := c.ApiClient.AuthenticateSteam(c.ServerKey, "", request, create, username, nil, make(map[string]string))

	if err != nil {
		return nil, err
	}

	// Return a new Session object
	return &Session{
		Token:        apiSession.Token,
		RefreshToken: apiSession.RefreshToken,
		Created:      apiSession.Created,
	}, nil
}

// BanGroupUsers bans users from a group.
func (c *Client) BanGroupUsers(session *Session, groupId string, ids []string) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.BanGroupUsers(&session.Token, &groupId, ids, make(map[string]string))
}

// BlockFriends blocks one or more users by ID or username.
func (c *Client) BlockFriends(session *Session, ids []string, usernames []string) error {
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return errors.As(err)
		}
	}

	return c.ApiClient.BlockFriends(&session.Token, ids, usernames, make(map[string]string))
}

// CreateGroup creates a new group with the current user as the creator and superadmin.
func (c *Client) CreateGroup(session *Session, request api.CreateGroupRequest) (*api.Group, error) {
	// Check if the session requires refresh
	if c.AutoRefreshSession && session.RefreshToken != "" &&
		session.IsExpired((time.Now().Unix()+c.ExpiredTimespanMs)/1000) {
		_, err := c.SessionRefresh(session, nil)
		if err != nil {
			return nil, err
		}
	}

	// Call the API client to create the group
	return c.ApiClient.CreateGroup(&session.Token, &request, make(map[string]string))
}

// CreateSocket creates a socket using the client's configuration.
func (c *Client) CreateSocket(eventHandle EventHandler, token string, useSSL bool, verbose bool, sendTimeoutMs *int, createStatus *bool) *DefaultSocket {
	return NewDefaultSocket(eventHandle, c.Host, c.Port, token, useSSL, verbose, sendTimeoutMs, createStatus)
}

// DeleteAccount deletes the current user's account.
func (c *Client) DeleteAccount(session *Session) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.DeleteAccount(session.Token, make(map[string]string))
}

// DeleteFriends deletes one or more users by ID or username.
func (c *Client) DeleteFriends(session *Session, ids []string, usernames []string) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}
	return c.ApiClient.DeleteFriends(&session.Token, ids, usernames, make(map[string]string))
}

// DeleteGroup deletes a group the user is part of and has permissions to delete.
func (c *Client) DeleteGroup(session *Session, groupId string) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.DeleteGroup(&session.Token, &groupId, make(map[string]string))
}

// DeleteNotifications deletes one or more notifications.
func (c *Client) DeleteNotifications(session *Session, ids []string) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.DeleteNotifications(session.Token, ids, make(map[string]string))
}

// DeleteStorageObjects deletes one or more storage objects.
func (c *Client) DeleteStorageObjects(session *Session, request *api.DeleteStorageObjectsRequest) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.DeleteStorageObjects(session.Token, request, make(map[string]string))
}

// DeleteTournamentRecord deletes a tournament record.
func (c *Client) DeleteTournamentRecord(session *Session, tournamentId string) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.DeleteTournamentRecord(session.Token, tournamentId, make(map[string]string))
}

// DemoteGroupUsers demotes a set of users in a group to the next role down.
func (c *Client) DemoteGroupUsers(session *Session, groupId *string, ids []string) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.DemoteGroupUsers(&session.Token, groupId, ids, make(map[string]string))
}

// EmitEvent submits an event for processing in the server's registered runtime custom events handler.
func (c *Client) EmitEvent(session *Session, request *api.Event) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.Event(&session.Token, request, make(map[string]string))
}

// GetAccount fetches the current user's account.
func (c *Client) GetAccount(session *Session) (*api.Account, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	return c.ApiClient.GetAccount(session.Token, make(map[string]string))
}

// GetSubscription fetches a subscription by product ID.
func (c *Client) GetSubscription(session *Session, productId *string) (*api.ValidatedSubscription, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	return c.ApiClient.GetSubscription(&session.Token, productId, make(map[string]string))
}

// ImportFacebookFriends imports Facebook friends and adds them to a user's account.
func (c *Client) ImportFacebookFriends(session *Session, request *api.AccountFacebook) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.ImportFacebookFriends(&session.Token, request, nil, make(map[string]string))
}

// ImportSteamFriends imports Steam friends and adds them to a user's account.
func (c *Client) ImportSteamFriends(session *Session, request *api.AccountSteam, reset bool) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.ImportSteamFriends(&session.Token, request, &reset, make(map[string]string))
}

// FetchUsers fetches zero or more users by ID and/or username.
func (c *Client) FetchUsers(session *Session, ids []string, usernames []string, facebookIds []string) (*api.Users, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	return c.ApiClient.GetUsers(&session.Token, ids, usernames, facebookIds, make(map[string]string))
}

// JoinGroup either joins a group that's open or sends a request to join a group that's closed.
func (c *Client) JoinGroup(session *Session, groupId string) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.JoinGroup(&session.Token, &groupId, make(map[string]string))
}

// JoinTournament allows a user to join a tournament by its ID.
func (c *Client) JoinTournament(session *Session, tournamentId string) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.JoinTournament(session.Token, tournamentId, make(map[string]string))
}

// KickGroupUsers kicks users from a group or declines their join requests.
func (c *Client) KickGroupUsers(session *Session, groupId string, ids []string) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.KickGroupUsers(&session.Token, &groupId, ids, make(map[string]string))
}

// LeaveGroup allows a user to leave a group they are part of.
func (c *Client) LeaveGroup(session *Session, groupId string) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.LeaveGroup(&session.Token, &groupId, make(map[string]string))
}

// ListChannelMessages retrieves a channel's message history.
func (c *Client) ListChannelMessages(session *Session, channelId string, limit *int, forward *bool, cursor *string) (*api.ChannelMessageList, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	return c.ApiClient.ListChannelMessages(&session.Token, &channelId, limit, forward, cursor, make(map[string]string))
}

// ListGroupUsers retrieves a group's users with optional state, limit, and cursor parameters.
func (c *Client) ListGroupUsers(session *Session, groupId string, state *int, limit *int, cursor *string) (*api.GroupUserList, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	return c.ApiClient.ListGroupUsers(&session.Token, &groupId, state, limit, cursor, make(map[string]string))
}

// ListUserGroups lists a user's groups.
func (c *Client) ListUserGroups(session *Session, userId string, state *int, limit int, cursor string) (*api.UserGroupList, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	return c.ApiClient.ListUserGroups(session.Token, userId, state, limit, cursor, make(map[string]string))
}

// ListGroups retrieves a list of groups based on the given filters.
func (c *Client) ListGroups(session *Session, name *string, cursor *string, limit *int) (*api.GroupList, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	return c.ApiClient.ListGroups(&session.Token, name, cursor, limit, nil, nil, nil, make(map[string]string))
}

// LinkApple adds an Apple ID to the social profiles on the current user's account.
func (c *Client) LinkApple(session *Session, request *api.AccountApple) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.LinkApple(session.Token, request, make(map[string]string))
}

// LinkCustom adds a custom ID to the social profiles on the current user's account.
func (c *Client) LinkCustom(session *Session, request *api.AccountCustom) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.LinkCustom(session.Token, request, make(map[string]string))
}

// LinkDevice adds a device ID to the social profiles on the current user's account.
func (c *Client) LinkDevice(session *Session, request *api.AccountDevice) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.LinkDevice(session.Token, request, make(map[string]string))
}

// LinkEmail adds an email and password to the social profiles on the current user's account.
func (c *Client) LinkEmail(session *Session, request *api.AccountEmail) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.LinkEmail(session.Token, request, make(map[string]string))
}

// LinkFacebook adds a Facebook ID to the social profiles on the current user's account.
func (c *Client) LinkFacebook(session *Session, request *api.AccountFacebook) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.LinkFacebook(session.Token, request, nil, make(map[string]string))
}

// LinkFacebookInstant adds Facebook Instant to the social profiles on the current user's account.
func (c *Client) LinkFacebookInstant(session *Session, request *api.AccountFacebookInstantGame) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.LinkFacebookInstantGame(session.Token, request, make(map[string]string))
}

// LinkGoogle adds a Google account to the social profiles on the current user's account.
func (c *Client) LinkGoogle(session *Session, request *api.AccountGoogle) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.LinkGoogle(session.Token, request, make(map[string]string))
}

// LinkGameCenter adds GameCenter to the social profiles on the current user's account.
func (c *Client) LinkGameCenter(session *Session, request *api.AccountGameCenter) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}
	return c.ApiClient.LinkGameCenter(session.Token, request, make(map[string]string))
}

// LinkSteam adds Steam to the social profiles on the current user's account.
func (c *Client) LinkSteam(session *Session, request *api.LinkSteamRequest) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.LinkSteam(session.Token, request, make(map[string]string))
}

// ListFriends lists all friends for the current user.
func (c *Client) ListFriends(session *Session, state *int, limit *int, cursor *string) (*api.FriendList, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	return c.ApiClient.ListFriends(&session.Token, limit, state, cursor, make(map[string]string))
}

// ListFriendsOfFriends lists the friends of friends for the current user.
func (c *Client) ListFriendsOfFriends(session *Session, limit *int, cursor *string) (*api.FriendsOfFriendsList, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	return c.ApiClient.ListFriendsOfFriends(&session.Token, limit, cursor, make(map[string]string))
}

// ListLeaderboardRecords lists the leaderboard records with optional ownerIds, pagination, and expiry filters.
func (c *Client) ListLeaderboardRecords(session *Session, leaderboardId string, ownerIds []string, limit *int, cursor *string, expiry *string) (*api.LeaderboardRecordList, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	return c.ApiClient.ListLeaderboardRecords(&session.Token, &leaderboardId, ownerIds, limit, cursor, expiry, make(map[string]string))
}

func (c *Client) ListLeaderboardRecordsAroundOwner(session *Session, leaderboardId string, ownerId string, limit int, expiry string, cursor string) (*api.LeaderboardRecordList, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	return c.ApiClient.ListLeaderboardRecordsAroundOwner(session.Token, leaderboardId, ownerId, limit, expiry, cursor, make(map[string]string))
}

// ListMatches fetches a list of running matches.
func (c *Client) ListMatches(session *Session, limit int, authoritative *bool, label string, minSize int, maxSize int, query string) (*api.MatchList, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	return c.ApiClient.ListMatches(session.Token, limit, authoritative, label, minSize, maxSize, query, make(map[string]string))
}

// ListNotifications fetches a list of notifications.
func (c *Client) ListNotifications(session *Session, limit int, cacheableCursor string) (*api.NotificationList, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	return c.ApiClient.ListNotifications(session.Token, limit, cacheableCursor, make(map[string]string))
}

// ListStorageObjects retrieves a list of storage objects.
func (c *Client) ListStorageObjects(session *Session, collection string, userID string, limit int, cursor string) (*api.StorageObjectList, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	return c.ApiClient.ListStorageObjects(session.Token, collection, userID, limit, cursor, make(map[string]string))
}

// ListTournaments retrieves a list of current or upcoming tournaments.
func (c *Client) ListTournaments(session *Session, categoryStart *int, categoryEnd *int, startTime *int64, endTime *int64, limit int, cursor string) (*api.TournamentList, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	return c.ApiClient.ListTournaments(session.Token, categoryStart, categoryEnd, startTime, endTime, limit, cursor, make(map[string]string))
}

// ListSubscriptions lists user subscriptions.
func (c *Client) ListSubscriptions(session *Session, cursor string, limit int32) (*api.SubscriptionList, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	return c.ApiClient.ListSubscriptions(
		&session.Token, &api.ListSubscriptionsRequest{
			Cursor: cursor,
			Limit:  wrapperspb.Int32(limit),
		},
		make(map[string]string),
	)
}

// ListTournamentRecords lists tournament records from a given tournament.
func (c *Client) ListTournamentRecords(
	session *Session,
	tournamentId string,
	ownerIds []string,
	limit int,
	cursor string,
	expiry string,
) (*api.TournamentRecordList, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	// Call the API to list tournament records.
	return c.ApiClient.ListTournamentRecords(
		session.Token,
		tournamentId,
		ownerIds,
		limit,
		cursor,
		expiry,
		make(map[string]string),
	)
}

// ListTournamentRecordsAroundOwner lists tournament records around a specific owner.
func (c *Client) ListTournamentRecordsAroundOwner(
	session *Session,
	tournamentId string,
	ownerId string,
	limit int,
	expiry string,
	cursor string,
) (*api.TournamentRecordList, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	// Call the API to get tournament records around owner.
	return c.ApiClient.ListTournamentRecordsAroundOwner(
		session.Token,
		tournamentId,
		ownerId,
		limit,
		expiry,
		cursor,
		make(map[string]string),
	)
}

// PromoteGroupUsers promotes the users in a group to the next role up.
func (c *Client) PromoteGroupUsers(session *Session, groupId string, ids []string) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.PromoteGroupUsers(session.Token, groupId, ids, make(map[string]string))
}

// ReadStorageObjects fetches storage objects.
func (c *Client) ReadStorageObjects(session *Session, request *api.ReadStorageObjectsRequest) (*api.StorageObjects, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	return c.ApiClient.ReadStorageObjects(session.Token, request, make(map[string]string))
}

// Rpc executes an RPC function on the server.
func (c *Client) Rpc(session *Session, id string, input map[string]interface{}) (*api.Rpc, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	// Serialize the input to JSON
	inputJson, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize input to JSON: %w", err)
	}
	jsonStr := string(inputJson)

	// Execute the RPC function on the API client
	return c.ApiClient.RpcFunc(session.Token, id, jsonStr, "", make(map[string]string))
}

// RpcHttpKey executes an RPC function on the server using an HTTP key.
func (c *Client) RpcHttpKey(httpKey, id string, input map[string]interface{}) (*api.Rpc, error) {
	// Serialize the input to JSON
	var inputJson string
	if input != nil {
		jsonBytes, err := json.Marshal(input)
		if err != nil {
			return nil, fmt.Errorf("failed to serialize input to JSON: %w", err)
		}
		inputJson = string(jsonBytes)
	}

	// Execute the RPC function on the API client
	return c.ApiClient.RpcFunc2("", id, inputJson, httpKey, make(map[string]string))
}

// SessionLogout logs out a session, invalidates a refresh token, or logs out all sessions/refresh tokens for a user.
func (c *Client) SessionLogout(session *Session, token, refreshToken string) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	// Create request payload for logout
	logoutRequest := api.SessionLogoutRequest{
		Token:        token,
		RefreshToken: refreshToken,
	}

	// Call the API client's session logout function
	return c.ApiClient.SessionLogout(session.Token, &logoutRequest, make(map[string]string))
}

// SessionRefresh refreshes a user's session using a refresh token retrieved from a previous authentication request.
func (c *Client) SessionRefresh(session *Session, vars map[string]string) (*Session, error) {
	if session == nil {
		return nil, fmt.Errorf("cannot refresh a null session")
	}

	if session.ExpiresAt > 0 && session.CreatedAt > 0 && session.ExpiresAt-session.CreatedAt < 70 {
		log.Println("Session lifetime too short, please set '--session.token_expiry_sec' option. See the documentation for more info: https://heroiclabs.com/docs/nakama/getting-started/configuration/#session")
	}

	if session.RefreshExpiresAt > 0 && session.CreatedAt > 0 && session.RefreshExpiresAt-session.CreatedAt < 3700 {
		log.Println("Session refresh lifetime too short, please set '--session.refresh_token_expiry_sec' option. See the documentation for more info: https://heroiclabs.com/docs/nakama/getting-started/configuration/#session")
	}

	apiSession, err := c.ApiClient.SessionRefresh(c.ServerKey, "", &api.SessionRefreshRequest{
		Token: session.RefreshToken,
		Vars:  vars,
	}, make(map[string]string))

	if err != nil {
		return nil, err
	}

	session.Update(apiSession.Token, apiSession.RefreshToken)
	return session, nil
}

// UnlinkApple removes the Apple ID from the social profiles on the current user's account.
func (c *Client) UnlinkApple(session *Session, request *api.AccountApple) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.UnlinkApple(session.Token, request, make(map[string]string))
}

// UnlinkCustom removes a custom ID from the social profiles on the current user's account.
func (c *Client) UnlinkCustom(session *Session, request *api.AccountCustom) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.UnlinkCustom(session.Token, request, make(map[string]string))
}

// UnlinkDevice removes a device ID from the social profiles on the current user's account.
func (c *Client) UnlinkDevice(session *Session, request *api.AccountDevice) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.UnlinkDevice(session.Token, request, make(map[string]string))
}

// UnlinkEmail removes an email+password from the social profiles on the current user's account.
func (c *Client) UnlinkEmail(session *Session, request *api.AccountEmail) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.UnlinkEmail(session.Token, request, make(map[string]string))
}

// UnlinkFacebook removes the Facebook ID from the social profiles on the current user's account.
func (c *Client) UnlinkFacebook(session *Session, request *api.AccountFacebook) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}
	return c.ApiClient.UnlinkFacebook(session.Token, request, make(map[string]string))
}

// UnlinkFacebookInstantGame removes Facebook Instant social profiles from the current user's account.
func (c *Client) UnlinkFacebookInstantGame(session *Session, request *api.AccountFacebookInstantGame) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.UnlinkFacebookInstantGame(session.Token, request, make(map[string]string))
}

// UnlinkGoogle removes the Google ID from the social profiles on the current user's account.
func (c *Client) UnlinkGoogle(session *Session, request *api.AccountGoogle) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.UnlinkGoogle(session.Token, request, make(map[string]string))
}

// UnlinkGameCenter removes GameCenter from the social profiles on the current user's account.
func (c *Client) UnlinkGameCenter(session *Session, request *api.AccountGameCenter) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.UnlinkGameCenter(session.Token, request, make(map[string]string))
}

// UnlinkSteam removes Steam from the social profiles on the current user's account.
func (c *Client) UnlinkSteam(session *Session, request *api.AccountSteam) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.UnlinkSteam(session.Token, request, make(map[string]string))
}

// UpdateAccount updates fields in the current user's account.
func (c *Client) UpdateAccount(session *Session, request *api.UpdateAccountRequest) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.UpdateAccount(session.Token, request, make(map[string]string))
}

// UpdateGroup updates a group the user is part of and has permissions to update.
func (c *Client) UpdateGroup(session *Session, groupId string, request *api.UpdateGroupRequest) error {
	if err := c.refreshSession(session); err != nil {
		return errors.As(err)
	}

	return c.ApiClient.UpdateGroup(session.Token, &groupId, request, make(map[string]string))
}

// ValidatePurchaseApple validates an Apple IAP receipt.
func (c *Client) ValidatePurchaseApple(session *Session, receipt string, persist bool) (*api.ValidatePurchaseResponse, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}
	response, err := c.ApiClient.ValidatePurchaseApple(&session.Token, &api.ValidatePurchaseAppleRequest{
		Receipt: receipt,
		Persist: wrapperspb.Bool(persist),
	}, make(map[string]string))
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ValidatePurchaseFacebookInstant validates a Facebook Instant IAP receipt.
func (c *Client) ValidatePurchaseFacebookInstant(session *Session, signedRequest string, persist bool) (*api.ValidatePurchaseResponse, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	response, err := c.ApiClient.ValidatePurchaseFacebookInstant(&session.Token, &api.ValidatePurchaseFacebookInstantRequest{
		SignedRequest: signedRequest,
		Persist:       wrapperspb.Bool(persist),
	}, make(map[string]string))
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ValidatePurchaseGoogle validates a Google IAP receipt.
func (c *Client) ValidatePurchaseGoogle(session *Session, purchase string, persist bool) (*api.ValidatePurchaseResponse, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	response, err := c.ApiClient.ValidatePurchaseGoogle(&session.Token, &api.ValidatePurchaseGoogleRequest{
		Purchase: purchase,
		Persist:  wrapperspb.Bool(persist),
	}, make(map[string]string))
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ValidatePurchaseHuawei validates a Huawei IAP receipt.
func (c *Client) ValidatePurchaseHuawei(session *Session, purchase string, signature string, persist bool) (*api.ValidatePurchaseResponse, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	response, err := c.ApiClient.ValidatePurchaseHuawei(&session.Token, &api.ValidatePurchaseHuaweiRequest{
		Purchase:  purchase,
		Signature: signature,
		Persist:   wrapperspb.Bool(persist),
	}, make(map[string]string))
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ValidateSubscriptionApple validates an Apple subscription receipt.
func (c *Client) ValidateSubscriptionApple(session *Session, receipt string, persist bool) (*api.ValidateSubscriptionResponse, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	response, err := c.ApiClient.ValidateSubscriptionApple(&session.Token, &api.ValidateSubscriptionAppleRequest{
		Receipt: receipt,
		Persist: wrapperspb.Bool(persist),
	}, make(map[string]string))
	if err != nil {
		return nil, err
	}

	return response, nil
}

// ValidateSubscriptionGoogle validates a Google subscription receipt.
func (c *Client) ValidateSubscriptionGoogle(session *Session, receipt string, persist bool) (*api.ValidateSubscriptionResponse, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	response, err := c.ApiClient.ValidateSubscriptionGoogle(&session.Token, &api.ValidateSubscriptionGoogleRequest{
		Receipt: receipt,
		Persist: wrapperspb.Bool(persist),
	}, make(map[string]string))
	if err != nil {
		return nil, err
	}

	return response, nil
}

// WriteLeaderboardRecord writes a record to a leaderboard.
func (c *Client) WriteLeaderboardRecord(session *Session, leaderboardId string, request *api.WriteLeaderboardRecordRequest_LeaderboardRecordWrite) (*api.LeaderboardRecord, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	return c.ApiClient.WriteLeaderboardRecord(
		session.Token,
		leaderboardId,
		request,
		make(map[string]string),
	)
}

// WriteStorageObjects writes storage objects.
func (c *Client) WriteStorageObjects(session *Session, objects []*api.WriteStorageObject) (*api.StorageObjectAcks, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	request := api.WriteStorageObjectsRequest{Objects: objects}
	storageObjects, err := c.ApiClient.WriteStorageObjects(session.Token, &request, make(map[string]string))
	if err != nil {
		return nil, err
	}

	return storageObjects, nil
}

// WriteTournamentRecord writes a record to a tournament.
func (c *Client) WriteTournamentRecord(session *Session, tournamentId string, request *api.WriteTournamentRecordRequest_TournamentRecordWrite) (*api.LeaderboardRecord, error) {
	if err := c.refreshSession(session); err != nil {
		return nil, errors.As(err)
	}

	return c.ApiClient.WriteTournamentRecord(
		session.Token,
		tournamentId,
		request,
		make(map[string]string),
	)
}
