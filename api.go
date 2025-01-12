package nakama

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ApiOperator int

const (
	ApiOperatorNoOverride ApiOperator = iota
	ApiOperatorBest
	ApiOperatorSet
	ApiOperatorIncrement
	ApiOperatorDecrement
)

type ApiStoreEnvironment int

const (
	ApiStoreEnvironmentUnknown ApiStoreEnvironment = iota
	ApiStoreEnvironmentSandbox
	ApiStoreEnvironmentProduction
)

type ApiStoreProvider int

const (
	ApiStoreProviderAppleAppStore ApiStoreProvider = iota
	ApiStoreProviderGooglePlayStore
	ApiStoreProviderHuaweiAppGallery
	ApiStoreProviderFacebookInstantStore
)

// A friend of a friend
type FriendsOfFriendsListFriendOfFriend struct {
	Referrer *string  `json:"referrer,omitempty"`
	User     *ApiUser `json:"user,omitempty"`
}

// A single user-role pair
type GroupUserListGroupUser struct {
	State *int     `json:"state,omitempty"`
	User  *ApiUser `json:"user,omitempty"`
}

// A single group-role pair
type UserGroupListUserGroup struct {
	Group *ApiGroup `json:"group,omitempty"`
	State *int      `json:"state,omitempty"`
}

// Record values to write
type WriteLeaderboardRecordRequestLeaderboardRecordWrite struct {
	Metadata *string      `json:"metadata,omitempty"`
	Operator *ApiOperator `json:"operator,omitempty"`
	Score    *string      `json:"score,omitempty"`
	Subscore *string      `json:"subscore,omitempty"`
}

type WriteTournamentRecordRequestTournamentRecordWrite struct {
	Metadata *string      `json:"metadata,omitempty"`
	Operator *ApiOperator `json:"operator,omitempty"`
	Score    *string      `json:"score,omitempty"`
	Subscore *string      `json:"subscore,omitempty"`
}

// A user with additional account details
type ApiAccount struct {
	CustomID    *string            `json:"custom_id,omitempty"`
	Devices     []ApiAccountDevice `json:"devices,omitempty"`
	DisableTime *time.Time         `json:"disable_time,omitempty"`
	Email       *string            `json:"email,omitempty"`
	User        *ApiUser           `json:"user,omitempty"`
	VerifyTime  *time.Time         `json:"verify_time,omitempty"`
	Wallet      *string            `json:"wallet,omitempty"`
}

type ApiAccountApple struct {
	Token *string           `json:"token,omitempty"`
	Vars  map[string]string `json:"vars,omitempty"`
}

type ApiAccountCustom struct {
	ID   *string           `json:"id,omitempty"`
	Vars map[string]string `json:"vars,omitempty"`
}

type ApiAccountDevice struct {
	ID   *string           `json:"id,omitempty"`
	Vars map[string]string `json:"vars,omitempty"`
}

type ApiAccountEmail struct {
	Email    *string           `json:"email,omitempty"`
	Password *string           `json:"password,omitempty"`
	Vars     map[string]string `json:"vars,omitempty"`
}

type ApiAccountFacebook struct {
	Token *string           `json:"token,omitempty"`
	Vars  map[string]string `json:"vars,omitempty"`
}

type ApiAccountFacebookInstantGame struct {
	SignedPlayerInfo *string           `json:"signed_player_info,omitempty"`
	Vars             map[string]string `json:"vars,omitempty"`
}

type ApiAccountGameCenter struct {
	BundleID     *string           `json:"bundle_id,omitempty"`
	PlayerID     *string           `json:"player_id,omitempty"`
	PublicKeyURL *string           `json:"public_key_url,omitempty"`
	Salt         *string           `json:"salt,omitempty"`
	Signature    *string           `json:"signature,omitempty"`
	Timestamp    *string           `json:"timestamp_seconds,omitempty"`
	Vars         map[string]string `json:"vars,omitempty"`
}

type ApiAccountGoogle struct {
	Token *string           `json:"token,omitempty"`
	Vars  map[string]string `json:"vars,omitempty"`
}

type ApiAccountSteam struct {
	Token *string           `json:"token,omitempty"`
	Vars  map[string]string `json:"vars,omitempty"`
	Sync  *bool             `json:"sync,omitempty"`
}

type ApiChannelMessage struct {
	ChannelID  *string    `json:"channel_id,omitempty"`
	Code       *int       `json:"code,omitempty"`
	Content    *string    `json:"content,omitempty"`
	CreateTime *time.Time `json:"create_time,omitempty"`
	GroupID    *string    `json:"group_id,omitempty"`
	MessageID  *string    `json:"message_id,omitempty"`
	Persistent *bool      `json:"persistent,omitempty"`
	RoomName   *string    `json:"room_name,omitempty"`
	SenderID   *string    `json:"sender_id,omitempty"`
	UpdateTime *time.Time `json:"update_time,omitempty"`
	UserIDOne  *string    `json:"user_id_one,omitempty"`
	UserIDTwo  *string    `json:"user_id_two,omitempty"`
	Username   *string    `json:"username,omitempty"`
}

type ApiChannelMessageList struct {
	CacheableCursor *string             `json:"cacheable_cursor,omitempty"`
	Messages        []ApiChannelMessage `json:"messages,omitempty"`
	NextCursor      *string             `json:"next_cursor,omitempty"`
	PrevCursor      *string             `json:"prev_cursor,omitempty"`
}

type ApiCreateGroupRequest struct {
	AvatarURL   *string `json:"avatar_url,omitempty"`
	Description *string `json:"description,omitempty"`
	LangTag     *string `json:"lang_tag,omitempty"`
	MaxCount    *int    `json:"max_count,omitempty"`
	Name        *string `json:"name,omitempty"`
	Open        *bool   `json:"open,omitempty"`
}

type ApiDeleteStorageObjectId struct {
	Collection *string `json:"collection,omitempty"`
	Key        *string `json:"key,omitempty"`
	Version    *string `json:"version,omitempty"`
}

type ApiDeleteStorageObjectsRequest struct {
	ObjectIDs []ApiDeleteStorageObjectId `json:"object_ids,omitempty"`
}

type ApiEvent struct {
	External   *bool             `json:"external,omitempty"`
	Name       *string           `json:"name,omitempty"`
	Properties map[string]string `json:"properties,omitempty"`
	Timestamp  *time.Time        `json:"timestamp,omitempty"`
}

type ApiFriend struct {
	State      *int       `json:"state,omitempty"`
	UpdateTime *time.Time `json:"update_time,omitempty"`
	User       *ApiUser   `json:"user,omitempty"`
}

type ApiFriendList struct {
	Cursor  *string     `json:"cursor,omitempty"`
	Friends []ApiFriend `json:"friends,omitempty"`
}

type ApiFriendsOfFriendsList struct {
	Cursor           *string                              `json:"cursor,omitempty"`
	FriendsOfFriends []FriendsOfFriendsListFriendOfFriend `json:"friends_of_friends,omitempty"`
}

type ApiGroup struct {
	AvatarURL   *string    `json:"avatar_url,omitempty"`
	CreateTime  *time.Time `json:"create_time,omitempty"`
	CreatorID   *string    `json:"creator_id,omitempty"`
	Description *string    `json:"description,omitempty"`
	EdgeCount   *int       `json:"edge_count,omitempty"`
	ID          *string    `json:"id,omitempty"`
	LangTag     *string    `json:"lang_tag,omitempty"`
	MaxCount    *int       `json:"max_count,omitempty"`
	Metadata    *string    `json:"metadata,omitempty"`
	Name        *string    `json:"name,omitempty"`
	Open        *bool      `json:"open,omitempty"`
	UpdateTime  *time.Time `json:"update_time,omitempty"`
}

type ApiGroupList struct {
	Cursor *string    `json:"cursor,omitempty"`
	Groups []ApiGroup `json:"groups,omitempty"`
}

type ApiGroupUserList struct {
	Cursor     *string                  `json:"cursor,omitempty"`
	GroupUsers []GroupUserListGroupUser `json:"group_users,omitempty"`
}

type ApiLeaderboardRecord struct {
	CreateTime    *time.Time `json:"create_time,omitempty"`
	ExpiryTime    *time.Time `json:"expiry_time,omitempty"`
	LeaderboardID *string    `json:"leaderboard_id,omitempty"`
	MaxNumScore   *int       `json:"max_num_score,omitempty"`
	Metadata      *string    `json:"metadata,omitempty"`
	NumScore      *int       `json:"num_score,omitempty"`
	OwnerID       *string    `json:"owner_id,omitempty"`
	Rank          *string    `json:"rank,omitempty"`
	Score         *string    `json:"score,omitempty"`
	Subscore      *string    `json:"subscore,omitempty"`
	UpdateTime    *time.Time `json:"update_time,omitempty"`
	Username      *string    `json:"username,omitempty"`
}

type ApiLeaderboardRecordList struct {
	NextCursor   *string                `json:"next_cursor,omitempty"`
	OwnerRecords []ApiLeaderboardRecord `json:"owner_records,omitempty"`
	PrevCursor   *string                `json:"prev_cursor,omitempty"`
	RankCount    *string                `json:"rank_count,omitempty"`
	Records      []ApiLeaderboardRecord `json:"records,omitempty"`
}

type ApiLinkSteamRequest struct {
	Account *ApiAccountSteam `json:"account,omitempty"`
	Sync    *bool            `json:"sync,omitempty"`
}

type ApiListSubscriptionsRequest struct {
	Cursor *string `json:"cursor,omitempty"`
	Limit  *int    `json:"limit,omitempty"`
}

type ApiMatch struct {
	Authoritative *bool   `json:"authoritative,omitempty"`
	HandlerName   *string `json:"handler_name,omitempty"`
	Label         *string `json:"label,omitempty"`
	MatchID       *string `json:"match_id,omitempty"`
	Size          *int    `json:"size,omitempty"`
	TickRate      *int    `json:"tick_rate,omitempty"`
}

type ApiMatchList struct {
	Matches []ApiMatch `json:"matches,omitempty"`
}

type ApiNotification struct {
	Code       *int       `json:"code,omitempty"`
	Content    *string    `json:"content,omitempty"`
	CreateTime *time.Time `json:"create_time,omitempty"`
	ID         *string    `json:"id,omitempty"`
	Persistent *bool      `json:"persistent,omitempty"`
	SenderID   *string    `json:"sender_id,omitempty"`
	Subject    *string    `json:"subject,omitempty"`
}

type ApiNotificationList struct {
	CacheableCursor *string           `json:"cacheable_cursor,omitempty"`
	Notifications   []ApiNotification `json:"notifications,omitempty"`
}

type ApiReadStorageObjectId struct {
	Collection *string `json:"collection,omitempty"`
	Key        *string `json:"key,omitempty"`
	UserID     *string `json:"user_id,omitempty"`
}

type ApiReadStorageObjectsRequest struct {
	ObjectIDs []ApiReadStorageObjectId `json:"object_ids,omitempty"`
}

type ApiRpc struct {
	HttpKey *string `json:"http_key,omitempty"`
	ID      *string `json:"id,omitempty"`
	Payload *string `json:"payload,omitempty"`
}

type ApiSession struct {
	Created      *bool   `json:"created,omitempty"`
	RefreshToken *string `json:"refresh_token,omitempty"`
	Token        *string `json:"token,omitempty"`
}

type ApiSessionLogoutRequest struct {
	RefreshToken *string `json:"refresh_token,omitempty"`
	Token        *string `json:"token,omitempty"`
}

type ApiSessionRefreshRequest struct {
	Token *string           `json:"token,omitempty"`
	Vars  map[string]string `json:"vars,omitempty"`
}

type ApiStorageObject struct {
	Collection      *string    `json:"collection,omitempty"`
	CreateTime      *time.Time `json:"create_time,omitempty"`
	Key             *string    `json:"key,omitempty"`
	PermissionRead  *int       `json:"permission_read,omitempty"`
	PermissionWrite *int       `json:"permission_write,omitempty"`
	UpdateTime      *time.Time `json:"update_time,omitempty"`
	UserID          *string    `json:"user_id,omitempty"`
	Value           *string    `json:"value,omitempty"`
	Version         *string    `json:"version,omitempty"`
}

type ApiStorageObjectAck struct {
	Collection *string    `json:"collection,omitempty"`
	CreateTime *time.Time `json:"create_time,omitempty"`
	Key        *string    `json:"key,omitempty"`
	UpdateTime *time.Time `json:"update_time,omitempty"`
	UserID     *string    `json:"user_id,omitempty"`
	Version    *string    `json:"version,omitempty"`
}

type ApiStorageObjectAcks struct {
	Acks []ApiStorageObjectAck `json:"acks,omitempty"`
}

type ApiStorageObjectList struct {
	Cursor  *string            `json:"cursor,omitempty"`
	Objects []ApiStorageObject `json:"objects,omitempty"`
}

type ApiStorageObjects struct {
	Objects []ApiStorageObject `json:"objects,omitempty"`
}

type ApiSubscriptionList struct {
	Cursor                 *string                    `json:"cursor,omitempty"`
	PrevCursor             *string                    `json:"prev_cursor,omitempty"`
	ValidatedSubscriptions []ApiValidatedSubscription `json:"validated_subscriptions,omitempty"`
}

type ApiTournament struct {
	Authoritative *bool        `json:"authoritative,omitempty"`
	CanEnter      *bool        `json:"can_enter,omitempty"`
	Category      *int         `json:"category,omitempty"`
	CreateTime    *time.Time   `json:"create_time,omitempty"`
	Description   *string      `json:"description,omitempty"`
	Duration      *int         `json:"duration,omitempty"`
	EndActive     *int64       `json:"end_active,omitempty"`
	EndTime       *time.Time   `json:"end_time,omitempty"`
	ID            *string      `json:"id,omitempty"`
	MaxNumScore   *int         `json:"max_num_score,omitempty"`
	MaxSize       *int         `json:"max_size,omitempty"`
	Metadata      *string      `json:"metadata,omitempty"`
	NextReset     *int64       `json:"next_reset,omitempty"`
	Operator      *ApiOperator `json:"operator,omitempty"`
	PrevReset     *int64       `json:"prev_reset,omitempty"`
	Size          *int         `json:"size,omitempty"`
	SortOrder     *int         `json:"sort_order,omitempty"`
	StartActive   *int64       `json:"start_active,omitempty"`
	StartTime     *time.Time   `json:"start_time,omitempty"`
	Title         *string      `json:"title,omitempty"`
}

type ApiTournamentList struct {
	Cursor      *string         `json:"cursor,omitempty"`
	Tournaments []ApiTournament `json:"tournaments,omitempty"`
}

type ApiTournamentRecordList struct {
	NextCursor   *string                `json:"next_cursor,omitempty"`
	OwnerRecords []ApiLeaderboardRecord `json:"owner_records,omitempty"`
	PrevCursor   *string                `json:"prev_cursor,omitempty"`
	RankCount    *string                `json:"rank_count,omitempty"`
	Records      []ApiLeaderboardRecord `json:"records,omitempty"`
}

type ApiUpdateAccountRequest struct {
	AvatarURL   *string `json:"avatar_url,omitempty"`
	DisplayName *string `json:"display_name,omitempty"`
	LangTag     *string `json:"lang_tag,omitempty"`
	Location    *string `json:"location,omitempty"`
	Timezone    *string `json:"timezone,omitempty"`
	Username    *string `json:"username,omitempty"`
}

type ApiUpdateGroupRequest struct {
	AvatarURL   *string `json:"avatar_url,omitempty"`
	Description *string `json:"description,omitempty"`
	GroupID     *string `json:"group_id,omitempty"`
	LangTag     *string `json:"lang_tag,omitempty"`
	Name        *string `json:"name,omitempty"`
	Open        *bool   `json:"open,omitempty"`
}

type ApiUser struct {
	AppleID               *string    `json:"apple_id,omitempty"`
	AvatarURL             *string    `json:"avatar_url,omitempty"`
	CreateTime            *time.Time `json:"create_time,omitempty"`
	DisplayName           *string    `json:"display_name,omitempty"`
	EdgeCount             *int       `json:"edge_count,omitempty"`
	FacebookID            *string    `json:"facebook_id,omitempty"`
	FacebookInstantGameID *string    `json:"facebook_instant_game_id,omitempty"`
	GamecenterID          *string    `json:"gamecenter_id,omitempty"`
	GoogleID              *string    `json:"google_id,omitempty"`
	ID                    *string    `json:"id,omitempty"`
	LangTag               *string    `json:"lang_tag,omitempty"`
	Location              *string    `json:"location,omitempty"`
	Metadata              *string    `json:"metadata,omitempty"`
	Online                *bool      `json:"online,omitempty"`
	SteamID               *string    `json:"steam_id,omitempty"`
	Timezone              *string    `json:"timezone,omitempty"`
	UpdateTime            *time.Time `json:"update_time,omitempty"`
	Username              *string    `json:"username,omitempty"`
}

// ApiUserGroupList A list of groups belonging to a user, along with the user's role in each group.
type ApiUserGroupList struct {
	// Cursor for the next page of results, if any.
	Cursor     *string                   `json:"cursor,omitempty"`
	UserGroups *[]UserGroupListUserGroup `json:"user_groups,omitempty"`
}

// ApiUsers A collection of zero or more users.
type ApiUsers struct {
	// The User objects.
	Users *[]ApiUser `json:"users,omitempty"`
}

// ApiValidatePurchaseAppleRequest Request to validate an Apple in-app purchase.
type ApiValidatePurchaseAppleRequest struct {
	Persist *bool   `json:"persist,omitempty"`
	Receipt *string `json:"receipt,omitempty"` // Base64 encoded Apple receipt data payload.
}

// ApiValidatePurchaseFacebookInstantRequest Request to validate a Facebook Instant in-app purchase.
type ApiValidatePurchaseFacebookInstantRequest struct {
	Persist       *bool   `json:"persist,omitempty"`
	SignedRequest *string `json:"signed_request,omitempty"` // Base64 encoded Facebook Instant signedRequest receipt data payload.
}

// ApiValidatePurchaseGoogleRequest Request to validate a Google Play Store in-app purchase.
type ApiValidatePurchaseGoogleRequest struct {
	Persist  *bool   `json:"persist,omitempty"`
	Purchase *string `json:"purchase,omitempty"` // JSON encoded Google purchase payload.
}

// ApiValidatePurchaseHuaweiRequest Request to validate a Huawei AppGallery in-app purchase.
type ApiValidatePurchaseHuaweiRequest struct {
	Persist   *bool   `json:"persist,omitempty"`
	Purchase  *string `json:"purchase,omitempty"`  // JSON encoded Huawei InAppPurchaseData.
	Signature *string `json:"signature,omitempty"` // InAppPurchaseData signature.
}

// ApiValidatePurchaseResponse Response for validated in-app purchases.
type ApiValidatePurchaseResponse struct {
	ValidatedPurchases *[]ApiValidatedPurchase `json:"validated_purchases,omitempty"`
}

// ApiValidateSubscriptionAppleRequest Request to validate an Apple subscription.
type ApiValidateSubscriptionAppleRequest struct {
	Persist *bool   `json:"persist,omitempty"` // Persist the subscription.
	Receipt *string `json:"receipt,omitempty"` // Base64 encoded Apple receipt data payload.
}

// ApiValidateSubscriptionGoogleRequest Request to validate a Google Play subscription.
type ApiValidateSubscriptionGoogleRequest struct {
	Persist *bool   `json:"persist,omitempty"` // Persist the subscription.
	Receipt *string `json:"receipt,omitempty"` // JSON encoded Google purchase payload.
}

// ApiValidateSubscriptionResponse Response for validated subscriptions.
type ApiValidateSubscriptionResponse struct {
	ValidatedSubscription *ApiValidatedSubscription `json:"validated_subscription,omitempty"`
}

// ApiValidatedPurchase Validated Purchase stored by the backend system.
type ApiValidatedPurchase struct {
	CreateTime       *string              `json:"create_time,omitempty"`       // Timestamp when the receipt validation was stored in DB.
	Environment      *ApiStoreEnvironment `json:"environment,omitempty"`       // Whether the purchase was done in production or sandbox environment.
	ProductID        *string              `json:"product_id,omitempty"`        // Purchase Product ID.
	ProviderResponse *string              `json:"provider_response,omitempty"` // Raw provider validation response.
	PurchaseTime     *string              `json:"purchase_time,omitempty"`     // Timestamp when the purchase was done.
	RefundTime       *string              `json:"refund_time,omitempty"`
	SeenBefore       *bool                `json:"seen_before,omitempty"` // Whether the purchase had already been validated before.
	Store            *ApiStoreProvider    `json:"store,omitempty"`
	TransactionID    *string              `json:"transaction_id,omitempty"` // Purchase Transaction ID.
	UpdateTime       *string              `json:"update_time,omitempty"`    // Timestamp when the receipt validation was updated.
	UserID           *string              `json:"user_id,omitempty"`        // Purchase User ID.
}

// ApiValidatedSubscription Validated Subscription stored by the backend system.
type ApiValidatedSubscription struct {
	Active                *bool                `json:"active,omitempty"`                  // Whether the subscription is currently active or not.
	CreateTime            *string              `json:"create_time,omitempty"`             // UNIX Timestamp when the receipt validation was stored in DB.
	Environment           *ApiStoreEnvironment `json:"environment,omitempty"`             // Whether the purchase was done in production or sandbox environment.
	ExpiryTime            *string              `json:"expiry_time,omitempty"`             // Subscription expiration time.
	OriginalTransactionID *string              `json:"original_transaction_id,omitempty"` // Purchase Original transaction ID.
	ProductID             *string              `json:"product_id,omitempty"`              // Purchase Product ID.
	ProviderNotification  *string              `json:"provider_notification,omitempty"`   // Raw provider notification body.
	ProviderResponse      *string              `json:"provider_response,omitempty"`       // Raw provider validation response body.
	PurchaseTime          *string              `json:"purchase_time,omitempty"`           // UNIX Timestamp when the purchase was done.
	RefundTime            *string              `json:"refund_time,omitempty"`             // Subscription refund time.
	Store                 *ApiStoreProvider    `json:"store,omitempty"`
	UpdateTime            *string              `json:"update_time,omitempty"` // UNIX Timestamp when the receipt validation was updated.
	UserID                *string              `json:"user_id,omitempty"`     // Subscription User ID.
}

// ApiWriteStorageObject The object to store in the database or storage engine.
type ApiWriteStorageObject struct {
	Collection      *string `json:"collection,omitempty"`       // The collection to store the object.
	Key             *string `json:"key,omitempty"`              // The key for the object within the collection.
	PermissionRead  *int    `json:"permission_read,omitempty"`  // Read access permissions for the object.
	PermissionWrite *int    `json:"permission_write,omitempty"` // Write access permissions for the object.
	Value           *string `json:"value,omitempty"`            // The value of the object.
	Version         *string `json:"version,omitempty"`          // Version hash for optimistic concurrency control.
}

// ApiWriteStorageObjectsRequest Request to write objects to the storage engine.
type ApiWriteStorageObjectsRequest struct {
	Objects *[]ApiWriteStorageObject `json:"objects,omitempty"` // The objects to store on the server.
}

type NakamaApi struct {
	ServerKey string
	BasePath  string
	TimeoutMs int
}

// Healthcheck is a healthcheck function that load balancers can use to check the service.
func (api *NakamaApi) Healthcheck(bearerToken string, options map[string]string) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/healthcheck"
	queryParams := url.Values{}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}
	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// DeleteAccount deletes the current user's account.
func (api *NakamaApi) DeleteAccount(bearerToken string, options map[string]string) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account"
	queryParams := url.Values{}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("DELETE", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}
	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// GetAccount fetches the current user's account.
func (api *NakamaApi) GetAccount(bearerToken string, options map[string]string) (*ApiAccount, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account"
	queryParams := url.Values{}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}
	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result *ApiAccount
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// UpdateAccount updates fields in the current user's account.
func (api *NakamaApi) UpdateAccount(bearerToken string, body ApiUpdateAccountRequest, options map[string]string) (any, error) {
	// Check if the body is nil
	if body == (ApiUpdateAccountRequest{}) {
		return nil, errors.New("'body' is a required parameter but is null or undefined")
	}

	// Define the URL path and query parameters
	urlPath := "/v2/account"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("PUT", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// AuthenticateApple authenticates a user with an Apple ID against the server.
func (api *NakamaApi) AuthenticateApple(basicAuthUsername string, basicAuthPassword string, account ApiAccountApple, create *bool, username *string, options map[string]string) (*ApiSession, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/authenticate/apple"
	queryParams := url.Values{}
	if create != nil {
		queryParams.Set("create", fmt.Sprintf("%v", *create))
	}
	if username != nil {
		queryParams.Set("username", *username)
	}

	// Convert the account to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Authorization header
	if basicAuthUsername != "" {
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(basicAuthUsername + ":" + basicAuthPassword))
		req.Header.Set("Authorization", "Basic "+encodedAuth)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result *ApiSession
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// AuthenticateCustom authenticates a user with a custom ID against the server.
func (api *NakamaApi) AuthenticateCustom(
	basicAuthUsername string,
	basicAuthPassword string,
	account ApiAccountCustom,
	create *bool,
	username *string,
	options map[string]string,
) (*ApiSession, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/authenticate/custom"
	queryParams := url.Values{}
	if create != nil {
		queryParams.Set("create", fmt.Sprintf("%v", *create))
	}
	if username != nil {
		queryParams.Set("username", *username)
	}

	// Convert the account to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Authorization header
	if basicAuthUsername != "" {
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(basicAuthUsername + ":" + basicAuthPassword))
		req.Header.Set("Authorization", "Basic "+encodedAuth)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result *ApiSession
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// AuthenticateDevice authenticates a user with a device ID against the server.
func (api *NakamaApi) AuthenticateDevice(
	basicAuthUsername string,
	basicAuthPassword string,
	account ApiAccountDevice,
	create *bool,
	username *string,
	options map[string]string,
) (*ApiSession, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/authenticate/device"
	queryParams := url.Values{}
	if create != nil {
		queryParams.Set("create", fmt.Sprintf("%v", *create))
	}
	if username != nil {
		queryParams.Set("username", *username)
	}

	// Convert the account to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Authorization header
	if basicAuthUsername != "" {
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(basicAuthUsername + ":" + basicAuthPassword))
		req.Header.Set("Authorization", "Basic "+encodedAuth)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result *ApiSession
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// AuthenticateEmail authenticates a user with an email and password against the server.
func (api *NakamaApi) AuthenticateEmail(
	basicAuthUsername string,
	basicAuthPassword string,
	account ApiAccountEmail,
	create *bool,
	username *string,
	options map[string]string,
) (*ApiSession, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/authenticate/email"
	queryParams := url.Values{}
	if create != nil {
		queryParams.Set("create", fmt.Sprintf("%v", *create))
	}
	if username != nil {
		queryParams.Set("username", *username)
	}

	// Convert the account to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Authorization header
	if basicAuthUsername != "" {
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(basicAuthUsername + ":" + basicAuthPassword))
		req.Header.Set("Authorization", "Basic "+encodedAuth)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result *ApiSession
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// AuthenticateFacebook authenticates a user with a Facebook OAuth token against the server.
func (api *NakamaApi) AuthenticateFacebook(
	basicAuthUsername string,
	basicAuthPassword string,
	account ApiAccountFacebook,
	create *bool,
	username *string,
	sync *bool,
	options map[string]string,
) (*ApiSession, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/authenticate/facebook"
	queryParams := url.Values{}
	if create != nil {
		queryParams.Set("create", fmt.Sprintf("%v", *create))
	}
	if username != nil {
		queryParams.Set("username", *username)
	}
	if sync != nil {
		queryParams.Set("sync", fmt.Sprintf("%v", *sync))
	}

	// Convert the account to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Authorization header
	if basicAuthUsername != "" {
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(basicAuthUsername + ":" + basicAuthPassword))
		req.Header.Set("Authorization", "Basic "+encodedAuth)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result *ApiSession
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// AuthenticateFacebookInstantGame authenticates a user with a Facebook Instant Game token against the server.
func (api *NakamaApi) AuthenticateFacebookInstantGame(
	basicAuthUsername string,
	basicAuthPassword string,
	account ApiAccountFacebookInstantGame,
	create *bool,
	username *string,
	options map[string]string,
) (*ApiSession, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/authenticate/facebookinstantgame"
	queryParams := url.Values{}
	if create != nil {
		queryParams.Set("create", fmt.Sprintf("%v", *create))
	}
	if username != nil {
		queryParams.Set("username", *username)
	}

	// Convert the account to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Authorization header
	if basicAuthUsername != "" {
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(basicAuthUsername + ":" + basicAuthPassword))
		req.Header.Set("Authorization", "Basic "+encodedAuth)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result *ApiSession
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// AuthenticateGameCenter authenticates a user with Apple's GameCenter against the server.
func (api *NakamaApi) AuthenticateGameCenter(
	basicAuthUsername string,
	basicAuthPassword string,
	account ApiAccountGameCenter,
	create *bool,
	username *string,
	options map[string]string,
) (*ApiSession, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/authenticate/gamecenter"
	queryParams := url.Values{}
	if create != nil {
		queryParams.Set("create", fmt.Sprintf("%v", *create))
	}
	if username != nil {
		queryParams.Set("username", *username)
	}

	// Convert the account to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Authorization header
	if basicAuthUsername != "" {
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(basicAuthUsername + ":" + basicAuthPassword))
		req.Header.Set("Authorization", "Basic "+encodedAuth)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result *ApiSession
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// AuthenticateGoogle authenticates a user with Google against the server.
func (api *NakamaApi) AuthenticateGoogle(
	basicAuthUsername string,
	basicAuthPassword string,
	account ApiAccountGoogle,
	create *bool,
	username *string,
	options map[string]string,
) (*ApiSession, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/authenticate/google"
	queryParams := url.Values{}
	if create != nil {
		queryParams.Set("create", fmt.Sprintf("%v", *create))
	}
	if username != nil {
		queryParams.Set("username", *username)
	}

	// Convert the account to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Authorization header
	if basicAuthUsername != "" {
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(basicAuthUsername + ":" + basicAuthPassword))
		req.Header.Set("Authorization", "Basic "+encodedAuth)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result *ApiSession
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// AuthenticateSteam authenticates a user with Steam against the server.
func (api *NakamaApi) AuthenticateSteam(
	basicAuthUsername string,
	basicAuthPassword string,
	account ApiAccountSteam,
	create *bool,
	username *string,
	sync *bool,
	options map[string]string,
) (*ApiSession, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/authenticate/steam"
	queryParams := url.Values{}
	if create != nil {
		queryParams.Set("create", fmt.Sprintf("%v", *create))
	}
	if username != nil {
		queryParams.Set("username", *username)
	}
	if sync != nil {
		queryParams.Set("sync", fmt.Sprintf("%v", *sync))
	}

	// Convert the account to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Authorization header
	if basicAuthUsername != "" {
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(basicAuthUsername + ":" + basicAuthPassword))
		req.Header.Set("Authorization", "Basic "+encodedAuth)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result *ApiSession
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// LinkApple adds an Apple ID to the social profiles on the current user's account.
func (api *NakamaApi) LinkApple(
	bearerToken string,
	body ApiAccountApple,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/link/apple"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// LinkCustom adds a custom ID to the social profiles on the current user's account.
func (api *NakamaApi) LinkCustom(
	bearerToken string,
	body ApiAccountCustom,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/link/custom"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// LinkDevice adds a device ID to the social profiles on the current user's account.
func (api *NakamaApi) LinkDevice(
	bearerToken string,
	body ApiAccountDevice,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/link/device"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// LinkEmail adds an email and password to the social profiles on the current user's account.
func (api *NakamaApi) LinkEmail(
	bearerToken string,
	body ApiAccountEmail,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/link/email"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// LinkFacebook adds a Facebook account to the social profiles on the current user's account.
func (api *NakamaApi) LinkFacebook(
	bearerToken string,
	account ApiAccountFacebook,
	sync *bool,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/link/facebook"
	queryParams := url.Values{}
	if sync != nil {
		queryParams.Set("sync", fmt.Sprintf("%t", *sync))
	}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// LinkFacebookInstantGame adds a Facebook Instant Game account to the social profiles on the current user's account.
func (api *NakamaApi) LinkFacebookInstantGame(
	bearerToken string,
	body ApiAccountFacebookInstantGame,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/link/facebookinstantgame"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// LinkGameCenter adds Apple's GameCenter to the social profiles on the current user's account.
func (api *NakamaApi) LinkGameCenter(
	bearerToken string,
	body ApiAccountGameCenter,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/link/gamecenter"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// LinkGoogle adds a Google account to the social profiles on the current user's account.
func (api *NakamaApi) LinkGoogle(
	bearerToken string,
	body ApiAccountGoogle,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/link/google"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// LinkSteam adds a Steam account to the social profiles on the current user's account.
func (api *NakamaApi) LinkSteam(
	bearerToken string,
	body ApiLinkSteamRequest,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/link/steam"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// SessionRefresh refreshes a user's session using a refresh token retrieved from a previous authentication request.
func (api *NakamaApi) SessionRefresh(
	basicAuthUsername string,
	basicAuthPassword string,
	body ApiSessionRefreshRequest,
	options map[string]string,
) (*ApiSession, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/session/refresh"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Auth header
	if basicAuthUsername != "" && basicAuthPassword != "" {
		auth := basicAuthUsername + ":" + basicAuthPassword
		req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(auth)))
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result *ApiSession
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// UnlinkApple removes the Apple ID from the social profiles on the current user's account.
func (api *NakamaApi) UnlinkApple(
	bearerToken string,
	body ApiAccountApple,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/unlink/apple"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// UnlinkCustom removes the custom ID from the social profiles on the current user's account.
func (api *NakamaApi) UnlinkCustom(
	bearerToken string,
	body ApiAccountCustom,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/unlink/custom"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// UnlinkDevice removes the device ID from the social profiles on the current user's account.
func (api *NakamaApi) UnlinkDevice(
	bearerToken string,
	body ApiAccountDevice,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/unlink/device"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// UnlinkEmail removes the email+password from the social profiles on the current user's account.
func (api *NakamaApi) UnlinkEmail(
	bearerToken string,
	body ApiAccountEmail,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/unlink/email"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// UnlinkFacebook removes the Facebook profile from the social profiles on the current user's account.
func (api *NakamaApi) UnlinkFacebook(
	bearerToken string,
	body ApiAccountFacebook,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/unlink/facebook"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// UnlinkFacebookInstantGame removes the Facebook Instant Game profile from the social profiles on the current user's account.
func (api *NakamaApi) UnlinkFacebookInstantGame(
	bearerToken string,
	body ApiAccountFacebookInstantGame,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/unlink/facebookinstantgame"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// UnlinkGameCenter removes the GameCenter profile from the social profiles on the current user's account.
func (api *NakamaApi) UnlinkGameCenter(
	bearerToken string,
	body ApiAccountGameCenter,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/unlink/gamecenter"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// UnlinkGoogle removes the Google profile from the social profiles on the current user's account.
func (api *NakamaApi) UnlinkGoogle(
	bearerToken string,
	body ApiAccountGoogle,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/unlink/google"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// UnlinkSteam removes the Steam profile from the social profiles on the current user's account.
func (api *NakamaApi) UnlinkSteam(
	bearerToken string,
	body ApiAccountSteam,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/unlink/steam"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// ListChannelMessages lists a channel's message history.
func (api *NakamaApi) ListChannelMessages(
	bearerToken string,
	channelId string,
	limit *int,
	forward *bool,
	cursor *string,
	options map[string]string,
) (ApiChannelMessageList, error) {
	if channelId == "" {
		return ApiChannelMessageList{}, errors.New("'channelId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/channel/{channelId}", "{channelId}", url.PathEscape(channelId), 1)
	queryParams := url.Values{}

	if limit != nil {
		queryParams.Set("limit", strconv.Itoa(*limit))
	}
	if forward != nil {
		queryParams.Set("forward", strconv.FormatBool(*forward))
	}
	if cursor != nil {
		queryParams.Set("cursor", *cursor)
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return ApiChannelMessageList{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiChannelMessageList{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiChannelMessageList{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiChannelMessageList{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result ApiChannelMessageList
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiChannelMessageList{}, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiChannelMessageList{}, err
			}
			return result, nil
		} else {
			return ApiChannelMessageList{}, errors.New(resp.Status)
		}
	}
}

// Event submits an event for processing in the server's registered runtime custom events handler.
func (api *NakamaApi) Event(
	bearerToken string,
	body ApiEvent,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/event"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

func (api *NakamaApi) DeleteFriends(
	bearerToken string,
	ids []string,
	usernames []string,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/friend"
	queryParams := url.Values{}

	if len(ids) > 0 {
		queryParams.Set("ids", strings.Join(ids, ","))
	}
	if len(usernames) > 0 {
		queryParams.Set("usernames", strings.Join(usernames, ","))
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("DELETE", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// ListFriends fetches the list of all friends for the current user.
func (api *NakamaApi) ListFriends(
	bearerToken string,
	limit *int,
	state *int,
	cursor *string,
	options map[string]string,
) (ApiFriendList, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/friend"
	queryParams := url.Values{}

	if limit != nil {
		queryParams.Set("limit", strconv.Itoa(*limit))
	}
	if state != nil {
		queryParams.Set("state", strconv.Itoa(*state))
	}
	if cursor != nil {
		queryParams.Set("cursor", *cursor)
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return ApiFriendList{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiFriendList{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiFriendList{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiFriendList{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result ApiFriendList
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiFriendList{}, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiFriendList{}, err
			}
			return result, nil
		} else {
			return ApiFriendList{}, errors.New(resp.Status)
		}
	}
}

func (api *NakamaApi) AddFriends(
	bearerToken string,
	ids []string,
	usernames []string,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/friend"
	queryParams := url.Values{}

	if len(ids) > 0 {
		queryParams.Set("ids", strings.Join(ids, ","))
	}
	if len(usernames) > 0 {
		queryParams.Set("usernames", strings.Join(usernames, ","))
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

func (api *NakamaApi) BlockFriends(
	bearerToken string,
	ids []string,
	usernames []string,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/friend/block"
	queryParams := url.Values{}

	if len(ids) > 0 {
		queryParams.Set("ids", strings.Join(ids, ","))
	}
	if len(usernames) > 0 {
		queryParams.Set("usernames", strings.Join(usernames, ","))
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

func (api *NakamaApi) ImportFacebookFriends(
	bearerToken string,
	account ApiAccountFacebook,
	reset bool,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/friend/facebook"
	queryParams := url.Values{}
	queryParams.Set("reset", strconv.FormatBool(reset))

	// Serialize the account object to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

func (api *NakamaApi) ListFriendsOfFriends(
	bearerToken string,
	limit *int,
	cursor *string,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/friend/friends"
	queryParams := url.Values{}

	if limit != nil {
		queryParams.Set("limit", strconv.Itoa(*limit))
	}
	if cursor != nil {
		queryParams.Set("cursor", *cursor)
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

func (api *NakamaApi) ImportSteamFriends(
	bearerToken string,
	account ApiAccountSteam,
	reset bool,
	options map[string]string,
) (any, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/friend/steam"
	queryParams := url.Values{}
	queryParams.Set("reset", strconv.FormatBool(reset))

	// Serialize the account object to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

func (api *NakamaApi) ListGroups(
	bearerToken string,
	name *string,
	cursor *string,
	limit *int,
	langTag *string,
	members *int,
	open *bool,
	options map[string]string,
) (*ApiGroupList, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/group"
	queryParams := url.Values{}

	if name != nil {
		queryParams.Set("name", *name)
	}
	if cursor != nil {
		queryParams.Set("cursor", *cursor)
	}
	if limit != nil {
		queryParams.Set("limit", strconv.Itoa(*limit))
	}
	if langTag != nil {
		queryParams.Set("lang_tag", *langTag)
	}
	if members != nil {
		queryParams.Set("members", strconv.Itoa(*members))
	}
	if open != nil {
		queryParams.Set("open", strconv.FormatBool(*open))
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result *ApiGroupList
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// CreateGroup creates a new group with the current user as the owner.
func (api *NakamaApi) CreateGroup(
	bearerToken string,
	body ApiCreateGroupRequest,
	options map[string]string,
) (ApiGroup, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/group"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return ApiGroup{}, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return ApiGroup{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiGroup{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiGroup{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiGroup{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result ApiGroup
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiGroup{}, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiGroup{}, err
			}
			return result, nil
		} else {
			return ApiGroup{}, errors.New(resp.Status)
		}
	}
}

// DeleteGroup deletes a group by ID.
func (api *NakamaApi) DeleteGroup(
	bearerToken string,
	groupId string,
	options map[string]string,
) (any, error) {
	if groupId == "" {
		return nil, errors.New("'groupId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}", "{groupId}", url.QueryEscape(groupId), 1)
	queryParams := url.Values{}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("DELETE", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// UpdateGroup updates fields in a given group.
func (api *NakamaApi) UpdateGroup(
	bearerToken string,
	groupId string,
	body ApiUpdateGroupRequest,
	options map[string]string,
) (any, error) {
	// Validate required parameters
	if groupId == "" {
		return nil, errors.New("'groupId' is a required parameter but is empty")
	}
	if (ApiUpdateGroupRequest{}) == body {
		return nil, errors.New("'body' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}", "{groupId}", url.QueryEscape(groupId), 1)
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("PUT", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// AddGroupUsers adds users to a group.
func (api *NakamaApi) AddGroupUsers(
	bearerToken string,
	groupId string,
	userIds []string,
	options map[string]string,
) (any, error) {

	// Check required parameters
	if groupId == "" {
		return nil, errors.New("'groupId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}/add", "{groupId}", url.QueryEscape(groupId), 1)
	queryParams := url.Values{}
	for _, userId := range userIds {
		queryParams.Add("user_ids", userId)
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// BanGroupUsers bans a set of users from a group.
func (api *NakamaApi) BanGroupUsers(
	bearerToken string,
	groupId string,
	userIds []string,
	options map[string]string,
) (any, error) {
	// Check required parameters
	if groupId == "" {
		return nil, errors.New("'groupId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}/ban", "{groupId}", url.QueryEscape(groupId), 1)
	queryParams := url.Values{}
	for _, userId := range userIds {
		queryParams.Add("user_ids", userId)
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// DemoteGroupUsers demotes a set of users in a group to the next role down.
func (api *NakamaApi) DemoteGroupUsers(
	bearerToken string,
	groupId string,
	userIds []string,
	options map[string]string,
) (any, error) {
	// Check required parameters
	if groupId == "" {
		return nil, errors.New("'groupId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}/demote", "{groupId}", url.QueryEscape(groupId), 1)
	queryParams := url.Values{}
	for _, userId := range userIds {
		queryParams.Add("user_ids", userId)
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// JoinGroup immediately joins an open group, or requests to join a closed one.
func (api *NakamaApi) JoinGroup(
	bearerToken string,
	groupId string,
	options map[string]string,
) (any, error) {
	if groupId == "" {
		return nil, errors.New("'groupId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}/join", "{groupId}", url.QueryEscape(groupId), 1)
	queryParams := url.Values{}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// KickGroupUsers kicks a set of users from a group.
func (api *NakamaApi) KickGroupUsers(
	bearerToken string,
	groupId string,
	userIds []string,
	options map[string]string,
) (any, error) {

	// Validate required parameter
	if groupId == "" {
		return nil, errors.New("'groupId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}/kick", "{groupId}", url.QueryEscape(groupId), 1)
	queryParams := url.Values{}
	for _, userId := range userIds {
		queryParams.Add("user_ids", userId)
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// LeaveGroup allows a user to leave a group they are a member of.
func (api *NakamaApi) LeaveGroup(
	bearerToken string,
	groupId string,
	options map[string]string,
) (any, error) {
	// Validate the required parameter
	if groupId == "" {
		return nil, errors.New("'groupId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}/leave", "{groupId}", url.QueryEscape(groupId), 1)
	queryParams := url.Values{}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// PromoteGroupUsers promotes a set of users in a group to the next role up.
func (api *NakamaApi) PromoteGroupUsers(
	bearerToken string,
	groupId string,
	userIds []string,
	options map[string]string,
) (any, error) {
	// Validate required parameter
	if groupId == "" {
		return nil, errors.New("'groupId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}/promote", "{groupId}", url.QueryEscape(groupId), 1)
	queryParams := url.Values{}
	for _, userId := range userIds {
		queryParams.Add("user_ids", userId)
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// ListGroupUsers lists all users that are part of a group.
func (api *NakamaApi) ListGroupUsers(
	bearerToken string,
	groupId string,
	limit *int,
	state *int,
	cursor *string,
	options map[string]string,
) (*ApiGroupUserList, error) {
	// Validate the required parameter
	if groupId == "" {
		return nil, errors.New("'groupId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}/user", "{groupId}", url.QueryEscape(groupId), 1)
	queryParams := url.Values{}
	if limit != nil {
		queryParams.Set("limit", strconv.Itoa(*limit))
	}
	if state != nil {
		queryParams.Set("state", strconv.Itoa(*state))
	}
	if cursor != nil {
		queryParams.Set("cursor", *cursor)
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result *ApiGroupUserList
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

func (api *NakamaApi) ValidatePurchaseApple(
	bearerToken string,
	body ApiValidatePurchaseAppleRequest,
	options map[string]string,
) (any, error) {
	// Define the URL path
	urlPath := "/v2/iap/purchase/apple"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// ValidatePurchaseFacebookInstant validates an Instant IAP receipt from Facebook.
func (api *NakamaApi) ValidatePurchaseFacebookInstant(
	bearerToken string,
	body ApiValidatePurchaseFacebookInstantRequest,
	options map[string]string,
) (any, error) {
	// Validate the required parameter
	if body == (ApiValidatePurchaseFacebookInstantRequest{}) {
		return nil, errors.New("'body' is a required parameter but is null or undefined.")
	}

	// Define the URL path
	urlPath := "/v2/iap/purchase/facebookinstant"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// ValidatePurchaseGoogle validates an IAP receipt from Google.
func (api *NakamaApi) ValidatePurchaseGoogle(
	bearerToken string,
	body ApiValidatePurchaseGoogleRequest,
	options map[string]string,
) (any, error) {
	// Validate the required parameter
	if body == (ApiValidatePurchaseGoogleRequest{}) {
		return nil, errors.New("'body' is a required parameter but is null or undefined.")
	}

	// Define the URL path
	urlPath := "/v2/iap/purchase/google"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// ValidatePurchaseHuawei validates an IAP receipt from Huawei.
func (api *NakamaApi) ValidatePurchaseHuawei(
	bearerToken string,
	body ApiValidatePurchaseHuaweiRequest,
	options map[string]string,
) (any, error) {
	// Validate the required parameter
	if body == (ApiValidatePurchaseHuaweiRequest{}) {
		return nil, errors.New("'body' is a required parameter but is null or undefined.")
	}

	// Define the URL path
	urlPath := "/v2/iap/purchase/huawei"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// ListSubscriptions lists user's subscriptions.
func (api *NakamaApi) ListSubscriptions(
	bearerToken string,
	body ApiListSubscriptionsRequest,
	options map[string]string,
) (ApiSubscriptionList, error) {

	// Validate the required parameter
	if body == (ApiListSubscriptionsRequest{}) {
		return ApiSubscriptionList{}, errors.New("'body' is a required parameter but is null or undefined.")
	}

	// Define the URL path
	urlPath := "/v2/iap/subscription"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return ApiSubscriptionList{}, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(bodyJson))
	if err != nil {
		return ApiSubscriptionList{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiSubscriptionList{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiSubscriptionList{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiSubscriptionList{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result ApiSubscriptionList
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiSubscriptionList{}, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiSubscriptionList{}, err
			}
			return result, nil
		} else {
			return ApiSubscriptionList{}, errors.New(resp.Status)
		}
	}
}

// ValidateSubscriptionApple validates an Apple subscription receipt.
func (api *NakamaApi) ValidateSubscriptionApple(
	bearerToken string,
	body ApiValidateSubscriptionAppleRequest,
	options map[string]string,
) (ApiValidateSubscriptionResponse, error) {

	// Validate the required parameter
	if body == (ApiValidateSubscriptionAppleRequest{}) {
		return ApiValidateSubscriptionResponse{}, errors.New("'body' is a required parameter but is null or undefined.")
	}

	// Define the URL path
	urlPath := "/v2/iap/subscription/apple"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return ApiValidateSubscriptionResponse{}, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(bodyJson))
	if err != nil {
		return ApiValidateSubscriptionResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiValidateSubscriptionResponse{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiValidateSubscriptionResponse{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiValidateSubscriptionResponse{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result ApiValidateSubscriptionResponse
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiValidateSubscriptionResponse{}, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiValidateSubscriptionResponse{}, err
			}
			return result, nil
		} else {
			return ApiValidateSubscriptionResponse{}, errors.New(resp.Status)
		}
	}
}

// ValidateSubscriptionGoogle validates a Google subscription receipt.
func (api *NakamaApi) ValidateSubscriptionGoogle(
	bearerToken string,
	body ApiValidateSubscriptionGoogleRequest,
	options map[string]string,
) (ApiValidateSubscriptionResponse, error) {

	// Validate the required parameter
	if body == (ApiValidateSubscriptionGoogleRequest{}) {
		return ApiValidateSubscriptionResponse{}, errors.New("'body' is a required parameter but is null or undefined.")
	}

	// Define the URL path
	urlPath := "/v2/iap/subscription/google"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return ApiValidateSubscriptionResponse{}, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(bodyJson))
	if err != nil {
		return ApiValidateSubscriptionResponse{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiValidateSubscriptionResponse{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiValidateSubscriptionResponse{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiValidateSubscriptionResponse{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result ApiValidateSubscriptionResponse
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiValidateSubscriptionResponse{}, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiValidateSubscriptionResponse{}, err
			}
			return result, nil
		} else {
			return ApiValidateSubscriptionResponse{}, errors.New(resp.Status)
		}
	}
}

// GetSubscription retrieves a subscription by product ID.
func (api *NakamaApi) GetSubscription(
	bearerToken string,
	productId string,
	options map[string]string,
) (ApiValidatedSubscription, error) {

	// Validate the required parameter
	if productId == "" {
		return ApiValidatedSubscription{}, errors.New("'productId' is a required parameter but is null or empty.")
	}

	// Define the URL path
	urlPath := fmt.Sprintf("/v2/iap/subscription/%s", url.QueryEscape(productId))
	queryParams := url.Values{}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return ApiValidatedSubscription{}, err
	}

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiValidatedSubscription{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiValidatedSubscription{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiValidatedSubscription{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result ApiValidatedSubscription
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiValidatedSubscription{}, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiValidatedSubscription{}, err
			}
			return result, nil
		} else {
			return ApiValidatedSubscription{}, errors.New(resp.Status)
		}
	}
}

// DeleteLeaderboardRecord deletes a leaderboard record.
func (api *NakamaApi) DeleteLeaderboardRecord(
	bearerToken string,
	leaderboardId string,
	options map[string]string,
) error {

	// Validate the required parameter
	if leaderboardId == "" {
		return errors.New("'leaderboardId' is a required parameter but is null or empty.")
	}

	// Define the URL path
	urlPath := fmt.Sprintf("/v2/leaderboard/%s", url.QueryEscape(leaderboardId))
	queryParams := url.Values{}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("DELETE", fullUrl, nil)
	if err != nil {
		return err
	}

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return errors.New("request timed out")
	case err := <-errorChan:
		return err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		} else {
			return errors.New(resp.Status)
		}
	}
}

// ListLeaderboardRecords retrieves a list of leaderboard records.
func (api *NakamaApi) ListLeaderboardRecords(
	bearerToken string,
	leaderboardId string,
	ownerIds []string,
	limit *int,
	cursor *string,
	expiry *string,
	options map[string]string,
) (ApiLeaderboardRecordList, error) {

	// Validate the required parameter
	if leaderboardId == "" {
		return ApiLeaderboardRecordList{}, errors.New("'leaderboardId' is a required parameter but is null or empty.")
	}

	// Define the URL path
	urlPath := fmt.Sprintf("/v2/leaderboard/%s", url.QueryEscape(leaderboardId))
	queryParams := url.Values{}

	// Add query parameters
	if len(ownerIds) > 0 {
		for _, ownerId := range ownerIds {
			queryParams.Add("owner_ids", ownerId)
		}
	}
	if limit != nil {
		queryParams.Set("limit", fmt.Sprintf("%d", *limit))
	}
	if cursor != nil {
		queryParams.Set("cursor", *cursor)
	}
	if expiry != nil {
		queryParams.Set("expiry", *expiry)
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return ApiLeaderboardRecordList{}, err
	}

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiLeaderboardRecordList{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiLeaderboardRecordList{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiLeaderboardRecordList{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result ApiLeaderboardRecordList
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiLeaderboardRecordList{}, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiLeaderboardRecordList{}, err
			}
			return result, nil
		} else {
			return ApiLeaderboardRecordList{}, errors.New(resp.Status)
		}
	}
}

// WriteLeaderboardRecord writes a record to a leaderboard.
func (api *NakamaApi) WriteLeaderboardRecord(
	bearerToken string,
	leaderboardId string,
	record WriteLeaderboardRecordRequestLeaderboardRecordWrite,
	options map[string]string,
) (ApiLeaderboardRecord, error) {

	// Validate the required parameters
	if leaderboardId == "" {
		return ApiLeaderboardRecord{}, errors.New("'leaderboardId' is a required parameter but is null or empty.")
	}
	if record == (WriteLeaderboardRecordRequestLeaderboardRecordWrite{}) {
		return ApiLeaderboardRecord{}, errors.New("'record' is a required parameter but is null or empty.")
	}

	// Define the URL path
	urlPath := fmt.Sprintf("/v2/leaderboard/%s", url.QueryEscape(leaderboardId))
	queryParams := url.Values{}

	// Convert the record to JSON
	bodyJson, err := json.Marshal(record)
	if err != nil {
		return ApiLeaderboardRecord{}, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return ApiLeaderboardRecord{}, err
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiLeaderboardRecord{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiLeaderboardRecord{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiLeaderboardRecord{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result ApiLeaderboardRecord
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiLeaderboardRecord{}, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiLeaderboardRecord{}, err
			}
			return result, nil
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ApiLeaderboardRecord{}, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

// ListLeaderboardRecordsAroundOwner lists leaderboard records that belong to a user.
func (api *NakamaApi) ListLeaderboardRecordsAroundOwner(
	bearerToken string,
	leaderboardId string,
	ownerId string,
	limit *int,
	expiry *string,
	cursor *string,
	options map[string]string,
) (ApiLeaderboardRecordList, error) {

	// Validate the required parameters
	if leaderboardId == "" {
		return ApiLeaderboardRecordList{}, errors.New("'leaderboardId' is a required parameter but is null or empty.")
	}
	if ownerId == "" {
		return ApiLeaderboardRecordList{}, errors.New("'ownerId' is a required parameter but is null or empty.")
	}

	// Define the URL path
	urlPath := fmt.Sprintf(
		"/v2/leaderboard/%s/owner/%s",
		url.QueryEscape(leaderboardId),
		url.QueryEscape(ownerId),
	)
	queryParams := url.Values{}

	// Add optional parameters to the query
	if limit != nil {
		queryParams.Set("limit", fmt.Sprintf("%d", *limit))
	}
	if expiry != nil {
		queryParams.Set("expiry", *expiry)
	}
	if cursor != nil {
		queryParams.Set("cursor", *cursor)
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return ApiLeaderboardRecordList{}, err
	}

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiLeaderboardRecordList{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiLeaderboardRecordList{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiLeaderboardRecordList{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result ApiLeaderboardRecordList
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiLeaderboardRecordList{}, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiLeaderboardRecordList{}, err
			}
			return result, nil
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ApiLeaderboardRecordList{}, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

func (api *NakamaApi) ListMatches(
	bearerToken string,
	limit *int,
	authoritative *bool,
	label *string,
	minSize *int,
	maxSize *int,
	query *string,
	options map[string]string,
) (ApiMatchList, error) {

	// Define the URL path
	urlPath := "/v2/match"
	queryParams := url.Values{}

	// Add optional parameters to the query
	if limit != nil {
		queryParams.Set("limit", fmt.Sprintf("%d", *limit))
	}
	if authoritative != nil {
		queryParams.Set("authoritative", fmt.Sprintf("%t", *authoritative))
	}
	if label != nil {
		queryParams.Set("label", *label)
	}
	if minSize != nil {
		queryParams.Set("min_size", fmt.Sprintf("%d", *minSize))
	}
	if maxSize != nil {
		queryParams.Set("max_size", fmt.Sprintf("%d", *maxSize))
	}
	if query != nil {
		queryParams.Set("query", *query)
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return ApiMatchList{}, err
	}

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiMatchList{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiMatchList{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiMatchList{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result ApiMatchList
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiMatchList{}, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiMatchList{}, err
			}
			return result, nil
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ApiMatchList{}, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

func (api *NakamaApi) DeleteNotifications(
	bearerToken string,
	ids []string,
	options map[string]string,
) (any, error) {

	// Define the URL path
	urlPath := "/v2/notification"
	queryParams := url.Values{}

	// Add ids to the query parameters
	for _, id := range ids {
		queryParams.Add("ids", id)
	}

	// Convert the body to JSON (if necessary this can be removed if API doesn't require a body for DELETE)
	bodyJson := ""

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("DELETE", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result any
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

func (api *NakamaApi) ListNotifications(
	bearerToken string,
	limit *int,
	cacheableCursor *string,
	options map[string]string,
) (ApiNotificationList, error) {

	// Define the URL path
	urlPath := "/v2/notification"
	queryParams := url.Values{}

	// Add query parameters
	if limit != nil {
		queryParams.Set("limit", fmt.Sprintf("%d", *limit))
	}
	if cacheableCursor != nil {
		queryParams.Set("cacheable_cursor", *cacheableCursor)
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return ApiNotificationList{}, err
	}

	// Set Bearer Token authorization header
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiNotificationList{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiNotificationList{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiNotificationList{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result ApiNotificationList
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiNotificationList{}, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiNotificationList{}, err
			}
			return result, nil
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ApiNotificationList{}, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

func (api *NakamaApi) RpcFunc2(
	bearerToken string,
	id string,
	payload *string,
	httpKey *string,
	options map[string]string,
) (ApiRpc, error) {

	// Validate the required parameter 'id'
	if id == "" {
		return ApiRpc{}, errors.New("'id' is a required parameter but is empty")
	}

	// Define the URL path
	urlPath := fmt.Sprintf("/v2/rpc/%s", url.QueryEscape(id))

	// Add query parameters
	queryParams := url.Values{}
	if payload != nil {
		queryParams.Set("payload", *payload)
	}
	if httpKey != nil {
		queryParams.Set("http_key", *httpKey)
	}

	// Convert the body to JSON (if necessary, can be modified)
	bodyJson := ""

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return ApiRpc{}, err
	}

	// Set Bearer Token authorization header if provided
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if provided
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiRpc{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiRpc{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiRpc{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result ApiRpc
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiRpc{}, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiRpc{}, err
			}
			return result, nil
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ApiRpc{}, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

func (api *NakamaApi) RpcFunc(
	bearerToken string,
	id string,
	body string,
	httpKey *string,
	options map[string]string,
) (ApiRpc, error) {
	// Validate the required parameters 'id' and 'body'
	if id == "" {
		return ApiRpc{}, errors.New("'id' is a required parameter but is empty")
	}
	if body == "" {
		return ApiRpc{}, errors.New("'body' is a required parameter but is empty")
	}

	// Define the URL path
	urlPath := fmt.Sprintf("/v2/rpc/%s", url.QueryEscape(id))

	// Add query parameters
	queryParams := url.Values{}
	if httpKey != nil {
		queryParams.Set("http_key", *httpKey)
	}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return ApiRpc{}, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, strings.NewReader(string(bodyJson)))
	if err != nil {
		return ApiRpc{}, err
	}

	// Set Bearer Token authorization header if provided
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if provided
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiRpc{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiRpc{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiRpc{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result ApiRpc
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiRpc{}, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiRpc{}, err
			}
			return result, nil
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ApiRpc{}, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

func (api *NakamaApi) SessionLogout(
	bearerToken string,
	body ApiSessionLogoutRequest,
	options map[string]string,
) (interface{}, error) {
	// Validate the required parameter 'body'
	if body == (ApiSessionLogoutRequest{}) {
		return nil, errors.New("'body' is a required parameter but is null or undefined")
	}

	// Define the URL path
	urlPath := "/v2/session/logout"

	// Add query parameters (empty for this request)
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, strings.NewReader(string(bodyJson)))
	if err != nil {
		return nil, err
	}

	// Set Bearer Token authorization header if provided
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if provided
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			var result interface{}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

func (api *NakamaApi) ReadStorageObjects(
	bearerToken string,
	body ApiReadStorageObjectsRequest,
	options map[string]string,
) (ApiStorageObjects, error) {
	// Define the URL path
	urlPath := "/v2/storage"

	// Add query parameters (empty for this request)
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return ApiStorageObjects{}, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, strings.NewReader(string(bodyJson)))
	if err != nil {
		return ApiStorageObjects{}, err
	}

	// Set Bearer Token authorization header if provided
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if provided
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiStorageObjects{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiStorageObjects{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiStorageObjects{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiStorageObjects{}, err
			}
			var result ApiStorageObjects
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiStorageObjects{}, err
			}
			return result, nil
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ApiStorageObjects{}, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

func (api *NakamaApi) WriteStorageObjects(
	bearerToken string,
	body ApiWriteStorageObjectsRequest,
	options map[string]string,
) (ApiStorageObjectAcks, error) {
	// Define the URL path
	urlPath := "/v2/storage"

	// Add query parameters (empty for this request)
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return ApiStorageObjectAcks{}, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("PUT", fullUrl, strings.NewReader(string(bodyJson)))
	if err != nil {
		return ApiStorageObjectAcks{}, err
	}

	// Set Bearer Token authorization header if provided
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if provided
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiStorageObjectAcks{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiStorageObjectAcks{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiStorageObjectAcks{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiStorageObjectAcks{}, err
			}
			var result ApiStorageObjectAcks
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiStorageObjectAcks{}, err
			}
			return result, nil
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ApiStorageObjectAcks{}, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

func (api *NakamaApi) DeleteStorageObjects(
	bearerToken string,
	body ApiDeleteStorageObjectsRequest,
	options map[string]string,
) (any, error) {
	// Define the URL path
	urlPath := "/v2/storage/delete"

	// Add query parameters (empty for this request)
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("PUT", fullUrl, strings.NewReader(string(bodyJson)))
	if err != nil {
		return nil, err
	}

	// Set Bearer Token authorization header if provided
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if provided
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			return bodyBytes, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			return bodyBytes, nil
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

func (api *NakamaApi) ListStorageObjects(
	bearerToken string,
	collection string,
	userId *string,
	limit *int,
	cursor *string,
	options map[string]string,
) (ApiStorageObjectList, error) {
	// Validate the 'collection' parameter
	if collection == "" {
		return ApiStorageObjectList{}, errors.New("'collection' is a required parameter but is empty.")
	}

	// Define the URL path and replace the placeholder
	urlPath := strings.Replace("/v2/storage/{collection}", "{collection}", url.QueryEscape(collection), 1)

	// Add query parameters
	queryParams := url.Values{}
	if userId != nil {
		queryParams.Set("user_id", *userId)
	}
	if limit != nil {
		queryParams.Set("limit", fmt.Sprintf("%d", *limit))
	}
	if cursor != nil {
		queryParams.Set("cursor", *cursor)
	}

	// No request body for this function
	bodyJson := ""

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return ApiStorageObjectList{}, err
	}

	// Set Bearer Token authorization header if provided
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if provided
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiStorageObjectList{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiStorageObjectList{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiStorageObjectList{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiStorageObjectList{}, err
			}
			var result ApiStorageObjectList
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiStorageObjectList{}, err
			}
			return result, nil
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ApiStorageObjectList{}, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

func (api *NakamaApi) ListStorageObjects2(
	bearerToken string,
	collection string,
	userId string,
	limit *int,
	cursor *string,
	options map[string]string,
) (ApiStorageObjectList, error) {

	// Validate 'collection' and 'userId' parameters
	if collection == "" {
		return ApiStorageObjectList{}, errors.New("'collection' is a required parameter but is empty.")
	}
	if userId == "" {
		return ApiStorageObjectList{}, errors.New("'userId' is a required parameter but is empty.")
	}

	// Define the URL path and replace placeholders
	urlPath := strings.Replace(
		strings.Replace("/v2/storage/{collection}/{userId}", "{collection}", url.QueryEscape(collection), 1),
		"{userId}", url.QueryEscape(userId), 1,
	)

	// Add query parameters
	queryParams := url.Values{}
	if limit != nil {
		queryParams.Set("limit", fmt.Sprintf("%d", *limit))
	}
	if cursor != nil {
		queryParams.Set("cursor", *cursor)
	}

	// No request body for this function
	bodyJson := ""

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return ApiStorageObjectList{}, err
	}

	// Set Bearer Token authorization header if provided
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if provided
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiStorageObjectList{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiStorageObjectList{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiStorageObjectList{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiStorageObjectList{}, err
			}
			var result ApiStorageObjectList
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiStorageObjectList{}, err
			}
			return result, nil
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ApiStorageObjectList{}, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

// ListTournaments lists current or upcoming tournaments.
func (api *NakamaApi) ListTournaments(
	bearerToken string,
	categoryStart *int,
	categoryEnd *int,
	startTime *int64,
	endTime *int64,
	limit *int,
	cursor *string,
	options map[string]string,
) (ApiTournamentList, error) {
	// Define the URL path
	urlPath := "/v2/tournament"

	// Add query parameters
	queryParams := url.Values{}
	if categoryStart != nil {
		queryParams.Set("category_start", fmt.Sprintf("%d", *categoryStart))
	}
	if categoryEnd != nil {
		queryParams.Set("category_end", fmt.Sprintf("%d", *categoryEnd))
	}
	if startTime != nil {
		queryParams.Set("start_time", fmt.Sprintf("%d", *startTime))
	}
	if endTime != nil {
		queryParams.Set("end_time", fmt.Sprintf("%d", *endTime))
	}
	if limit != nil {
		queryParams.Set("limit", fmt.Sprintf("%d", *limit))
	}
	if cursor != nil {
		queryParams.Set("cursor", *cursor)
	}

	// No request body for this function
	bodyJson := ""

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return ApiTournamentList{}, err
	}

	// Set Bearer Token authorization header if provided
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if provided
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiTournamentList{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiTournamentList{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiTournamentList{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiTournamentList{}, err
			}
			var result ApiTournamentList
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiTournamentList{}, err
			}
			return result, nil
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ApiTournamentList{}, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

// DeleteTournamentRecord deletes a tournament record.
func (api *NakamaApi) DeleteTournamentRecord(
	bearerToken string,
	tournamentId string,
	options map[string]string,
) (any, error) {
	// Validate the tournamentId
	if tournamentId == "" {
		return nil, errors.New("'tournamentId' is a required parameter but is empty.")
	}

	// Define the URL path
	urlPath := "/v2/tournament/" + url.QueryEscape(tournamentId)

	// No query parameters for this function
	queryParams := url.Values{}

	// No request body for DELETE request
	bodyJson := ""

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("DELETE", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}

	// Set Bearer Token authorization header if provided
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if provided
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check the HTTP status code
	if resp.StatusCode == http.StatusNoContent {
		// Success with no content
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return bodyBytes, nil
	} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Success with content, consume body (optional)
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return bodyBytes, nil
	} else {
		// Handle error response
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected response: %s", string(bodyBytes))
	}
}

// ListTournamentRecords lists tournament records.
func (api *NakamaApi) ListTournamentRecords(
	bearerToken string,
	tournamentId string,
	ownerIds []string,
	limit *int,
	cursor *string,
	expiry *string,
	options map[string]string,
) (ApiTournamentRecordList, error) {

	// Validate the tournamentId
	if tournamentId == "" {
		return ApiTournamentRecordList{}, errors.New("'tournamentId' is a required parameter but is empty.")
	}

	// Define the URL path
	urlPath := "/v2/tournament/" + url.QueryEscape(tournamentId)

	// Add query parameters
	queryParams := url.Values{}
	if len(ownerIds) > 0 {
		for _, id := range ownerIds {
			queryParams.Add("owner_ids", id)
		}
	}
	if limit != nil {
		queryParams.Set("limit", fmt.Sprintf("%d", *limit))
	}
	if cursor != nil {
		queryParams.Set("cursor", *cursor)
	}
	if expiry != nil {
		queryParams.Set("expiry", *expiry)
	}

	// No request body for this function
	bodyJson := ""

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return ApiTournamentRecordList{}, err
	}

	// Set Bearer Token authorization header if provided
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if provided
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiTournamentRecordList{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiTournamentRecordList{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return ApiTournamentRecordList{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiTournamentRecordList{}, err
			}
			var result ApiTournamentRecordList
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiTournamentRecordList{}, err
			}
			return result, nil
		} else {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ApiTournamentRecordList{}, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

// WriteTournamentRecord2 writes a record to a tournament.
func (api *NakamaApi) WriteTournamentRecord2(
	bearerToken string,
	tournamentId string,
	record WriteTournamentRecordRequestTournamentRecordWrite,
	options map[string]string,
) (ApiLeaderboardRecord, error) {

	// Validate the tournamentId and record
	if tournamentId == "" {
		return ApiLeaderboardRecord{}, errors.New("'tournamentId' is a required parameter but is empty.")
	}
	if record == (WriteTournamentRecordRequestTournamentRecordWrite{}) {
		return ApiLeaderboardRecord{}, errors.New("'record' is a required parameter but is empty.")
	}

	// Define the URL path
	urlPath := "/v2/tournament/" + url.QueryEscape(tournamentId)

	// Prepare the request body
	bodyJson, err := json.Marshal(record)
	if err != nil {
		return ApiLeaderboardRecord{}, fmt.Errorf("failed to marshal record: %w", err)
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, nil)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, strings.NewReader(string(bodyJson)))
	if err != nil {
		return ApiLeaderboardRecord{}, err
	}

	// Set Bearer Token authorization header if provided
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if provided
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return ApiLeaderboardRecord{}, err
	}
	defer resp.Body.Close()

	// Handle the HTTP response
	if resp.StatusCode == http.StatusNoContent {
		// Success with no content
		return ApiLeaderboardRecord{}, nil
	} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Success with content, parse response body
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return ApiLeaderboardRecord{}, err
		}
		var result ApiLeaderboardRecord
		err = json.Unmarshal(bodyBytes, &result)
		if err != nil {
			return ApiLeaderboardRecord{}, err
		}
		return result, nil
	} else {
		// Handle error response
		bodyBytes, _ := io.ReadAll(resp.Body)
		return ApiLeaderboardRecord{}, fmt.Errorf("unexpected response: %s", string(bodyBytes))
	}
}

// WriteTournamentRecord writes a record to a tournament.
func (api *NakamaApi) WriteTournamentRecord(
	bearerToken string,
	tournamentId string,
	record WriteTournamentRecordRequestTournamentRecordWrite,
	options map[string]string,
) (ApiLeaderboardRecord, error) {

	// Validate the tournamentId and record
	if tournamentId == "" {
		return ApiLeaderboardRecord{}, errors.New("'tournamentId' is a required parameter but is empty.")
	}
	if record == (WriteTournamentRecordRequestTournamentRecordWrite{}) {
		return ApiLeaderboardRecord{}, errors.New("'record' is a required parameter but is empty.")
	}

	// Define the URL path
	urlPath := "/v2/tournament/" + url.QueryEscape(tournamentId)

	// Prepare the request body
	bodyJson, err := json.Marshal(record)
	if err != nil {
		return ApiLeaderboardRecord{}, fmt.Errorf("failed to marshal record: %w", err)
	}

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, nil)

	// Prepare the HTTP request
	req, err := http.NewRequest("PUT", fullUrl, strings.NewReader(string(bodyJson)))
	if err != nil {
		return ApiLeaderboardRecord{}, err
	}

	// Set Bearer Token authorization header if provided
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if provided
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiLeaderboardRecord{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiLeaderboardRecord{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNoContent {
			return ApiLeaderboardRecord{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Success with content, parse response body
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiLeaderboardRecord{}, err
			}
			var result ApiLeaderboardRecord
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiLeaderboardRecord{}, err
			}
			return result, nil
		} else {
			// Handle error response
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ApiLeaderboardRecord{}, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

// JoinTournament attempts to join an open and running tournament.
func (api *NakamaApi) JoinTournament(
	bearerToken string,
	tournamentId string,
	options map[string]string,
) (interface{}, error) {

	// Validate the tournamentId
	if tournamentId == "" {
		return nil, errors.New("'tournamentId' is a required parameter but is empty.")
	}

	// Define the URL path
	urlPath := "/v2/tournament/" + url.QueryEscape(tournamentId) + "/join"

	// Prepare the query params (if any, currently empty map)
	queryParams := url.Values{}

	// Prepare the request body
	bodyJson := ""

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}

	// Set Bearer Token authorization header if provided
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if provided
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNoContent {
			// Success with no content
			return nil, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Success with content, parse response body
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			var result interface{}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return result, nil
		} else {
			// Handle error response
			bodyBytes, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

// ListTournamentRecordsAroundOwner lists tournament records for a given owner.
func (api *NakamaApi) ListTournamentRecordsAroundOwner(
	bearerToken string,
	tournamentId string,
	ownerId string,
	limit *int,
	expiry *string,
	cursor *string,
	options map[string]string,
) (ApiTournamentRecordList, error) {

	// Validate the tournamentId and ownerId
	if tournamentId == "" {
		return ApiTournamentRecordList{}, errors.New("'tournamentId' is a required parameter but is empty.")
	}
	if ownerId == "" {
		return ApiTournamentRecordList{}, errors.New("'ownerId' is a required parameter but is empty.")
	}

	// Define the URL path
	urlPath := "/v2/tournament/" + url.QueryEscape(tournamentId) + "/owner/" + url.QueryEscape(ownerId)

	// Prepare the query params
	queryParams := url.Values{}
	if limit != nil {
		queryParams.Set("limit", fmt.Sprintf("%d", *limit))
	}
	if expiry != nil {
		queryParams.Set("expiry", *expiry)
	}
	if cursor != nil {
		queryParams.Set("cursor", *cursor)
	}

	// Prepare the request body (empty)
	bodyJson := ""

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return ApiTournamentRecordList{}, err
	}

	// Set Bearer Token authorization header if provided
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if provided
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiTournamentRecordList{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiTournamentRecordList{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNoContent {
			// Success with no content
			return ApiTournamentRecordList{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Success with content, parse response body
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiTournamentRecordList{}, err
			}
			var result ApiTournamentRecordList
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiTournamentRecordList{}, err
			}
			return result, nil
		} else {
			// Handle error response
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ApiTournamentRecordList{}, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

// GetUsers fetches zero or more users by ID and/or username.
func (api *NakamaApi) GetUsers(
	bearerToken string,
	ids []string,
	usernames []string,
	facebookIds []string,
	options map[string]string,
) (ApiUsers, error) {

	// Define the URL path
	urlPath := "/v2/user"

	// Prepare the query params
	queryParams := url.Values{}
	if len(ids) > 0 {
		queryParams["ids"] = ids
	}
	if len(usernames) > 0 {
		queryParams["usernames"] = usernames
	}
	if len(facebookIds) > 0 {
		queryParams["facebook_ids"] = facebookIds
	}

	// Prepare the request body (empty)
	bodyJson := ""

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return ApiUsers{}, err
	}

	// Set Bearer Token authorization header if provided
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if provided
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiUsers{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiUsers{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNoContent {
			// Success with no content
			return ApiUsers{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Success with content, parse response body
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiUsers{}, err
			}
			var result ApiUsers
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiUsers{}, err
			}
			return result, nil
		} else {
			// Handle error response
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ApiUsers{}, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

// ListUserGroups lists the groups the current user belongs to.
func (api *NakamaApi) ListUserGroups(
	bearerToken string,
	userId string,
	limit *int,
	state *int,
	cursor *string,
	options map[string]string,
) (ApiUserGroupList, error) {

	// Validate required parameters
	if userId == "" {
		return ApiUserGroupList{}, errors.New("'userId' is a required parameter but is empty.")
	}

	// Define the URL path and replace placeholder
	urlPath := "/v2/user/{userId}/group"
	urlPath = strings.Replace(urlPath, "{userId}", url.QueryEscape(userId), 1)

	// Prepare the query params
	queryParams := url.Values{}
	if limit != nil {
		queryParams.Set("limit", strconv.Itoa(*limit))
	}
	if state != nil {
		queryParams.Set("state", strconv.Itoa(*state))
	}
	if cursor != nil {
		queryParams.Set("cursor", *cursor)
	}

	// Prepare the request body (empty)
	bodyJson := ""

	// Construct the full URL
	fullUrl := api.buildFullUrl(api.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return ApiUserGroupList{}, err
	}

	// Set Bearer Token authorization header if provided
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	// Apply additional custom headers or options if provided
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(api.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- err
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return ApiUserGroupList{}, errors.New("request timed out")
	case err := <-errorChan:
		return ApiUserGroupList{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNoContent {
			// Success with no content
			return ApiUserGroupList{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			// Success with content, parse response body
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return ApiUserGroupList{}, err
			}
			var result ApiUserGroupList
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return ApiUserGroupList{}, err
			}
			return result, nil
		} else {
			// Handle error response
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ApiUserGroupList{}, fmt.Errorf("unexpected response: %s", string(bodyBytes))
		}
	}
}

func (api *NakamaApi) buildFullUrl(basePath string, fragment string, queryParams url.Values) string {
	fullPath := basePath + fragment + "?"

	for k, values := range queryParams {
		for _, v := range values {
			fullPath += fmt.Sprintf("%s=%s&", url.QueryEscape(k), url.QueryEscape(v))
		}
	}

	// Remove the trailing "&" if present
	if strings.HasSuffix(fullPath, "&") {
		fullPath = fullPath[:len(fullPath)-1]
	}

	return fullPath
}
