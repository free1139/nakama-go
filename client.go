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

type LeaderboardRecord struct {
	CreateTime    *string
	ExpiryTime    *string
	LeaderboardID *string
	Metadata      map[string]interface{}
	NumScore      *int
	OwnerID       *string
	Rank          *int
	Score         *int
	SubScore      *int
	UpdateTime    *string
	Username      *string
	MaxNumScore   *int
}

type LeaderboardRecordList struct {
	NextCursor   *string
	OwnerRecords []LeaderboardRecord
	PrevCursor   *string
	RankCount    *int
	Records      []LeaderboardRecord
}

type Tournament struct {
	Authoritative *bool
	ID            *string
	Title         *string
	Description   *string
	Duration      *int
	Category      *int
	SortOrder     *int
	Size          *int
	MaxSize       *int
	MaxNumScore   *int
	CanEnter      *bool
	EndActive     *int
	NextReset     *int
	Metadata      map[string]interface{}
	CreateTime    *string
	StartTime     *string
	EndTime       *string
	StartActive   *int
}

type TournamentList struct {
	Tournaments []Tournament
	Cursor      *string
}

type TournamentRecordList struct {
	NextCursor   *string
	OwnerRecords []LeaderboardRecord
	PrevCursor   *string
	Records      []LeaderboardRecord
}

type WriteTournamentRecord struct {
	Metadata map[string]interface{}
	Score    *string
	SubScore *string
}

type WriteLeaderboardRecord struct {
	Metadata map[string]interface{}
	Score    *string
	SubScore *string
}

type WriteStorageObject struct {
	Collection      *string
	Key             *string
	PermissionRead  *int
	PermissionWrite *int
	Value           map[string]interface{}
	Version         *string
}

type StorageObject struct {
	Collection      *string
	CreateTime      *string
	Key             *string
	PermissionRead  *int
	PermissionWrite *int
	UpdateTime      *string
	UserID          *string
	Value           map[string]interface{}
	Version         *string
}

type StorageObjectList struct {
	Cursor  *string
	Objects []StorageObject
}

type StorageObjects struct {
	Objects []StorageObject
}

type ChannelMessage struct {
	ChannelID   *string
	Code        *int
	Content     map[string]interface{}
	CreateTime  *string
	GroupID     *string
	MessageID   *string
	Persistent  *bool
	RoomName    *string
	ReferenceID *string
	SenderID    *string
	UpdateTime  *string
	UserIDOne   *string
	UserIDTwo   *string
	Username    *string
}

type ChannelMessageList struct {
	CacheableCursor *string
	Messages        []ChannelMessage
	NextCursor      *string
	PrevCursor      *string
}

type User struct {
	AvatarURL             *string
	CreateTime            *string
	DisplayName           *string
	EdgeCount             *int
	FacebookID            *string
	FacebookInstantGameID *string
	GamecenterID          *string
	GoogleID              *string
	ID                    *string
	LangTag               *string
	Location              *string
	Metadata              map[string]interface{}
	Online                *bool
	SteamID               *string
	Timezone              *string
	UpdateTime            *string
	Username              *string
}

type Users struct {
	Users []User
}

type Friend struct {
	State *int
	User  *User
}

type Friends struct {
	Friends []Friend
	Cursor  *string
}

type FriendOfFriend struct {
	Referrer *string
	User     *User
}

type FriendsOfFriends struct {
	Cursor           *string
	FriendsOfFriends []FriendOfFriend
}

type GroupUser struct {
	User  *User
	State *int
}

type GroupUserList struct {
	GroupUsers []GroupUser
	Cursor     *string
}

type Group struct {
	AvatarURL   *string
	CreateTime  *string
	CreatorID   *string
	Description *string
	EdgeCount   *int
	ID          *string
	LangTag     *string
	MaxCount    *int
	Metadata    map[string]interface{}
	Name        *string
	Open        *bool
	UpdateTime  *string
}

type GroupList struct {
	Cursor *string
	Groups []Group
}

type UserGroup struct {
	Group *Group
	State *int
}

type UserGroupList struct {
	UserGroups []UserGroup
	Cursor     *string
}

type Notification struct {
	Code       *int
	Content    map[string]interface{}
	CreateTime *string
	ID         *string
	Persistent *bool
	SenderID   *string
	Subject    *string
}

type NotificationList struct {
	CacheableCursor *string
	Notifications   []Notification
}

type ValidatedSubscription struct {
	Active                *bool
	CreateTime            *string
	Environment           *string
	ExpiryTime            *string
	OriginalTransactionID *string
	ProductID             *string
	ProviderNotification  *string
	ProviderResponse      *string
	PurchaseTime          *string
	RefundTime            *string
	Store                 *string
	UpdateTime            *string
	UserID                *string
}

type SubscriptionList struct {
	Cursor                 *string
	PrevCursor             *string
	ValidatedSubscriptions []ValidatedSubscription
}
