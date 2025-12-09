package nakama

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// ISession represents a session authenticated for a user with the Nakama server.
type ISession interface {
	IsExpired(currentTime int64) bool
	IsRefreshExpired(currentTime int64) bool
}

// Session implements the ISession interface.
type Session struct {
	Token            string
	Created          bool
	CreatedAt        int64
	ExpiresAt        int64
	RefreshExpiresAt int64
	RefreshToken     string
	Username         string
	UserID           string
	Vars             map[string]interface{}
}

func (s *Session) ToJson() string {
	result, err := json.Marshal(s)
	if err != nil {
		panic(err)
	}
	return string(result)
}

// NewSession creates a new Session.
func NewSession(token, refreshToken string, created bool) *Session {
	unixTime := time.Now().Unix()
	session := &Session{
		Token:        token,
		RefreshToken: refreshToken,
		Created:      created,
		CreatedAt:    unixTime,
	}
	session.Update(token, refreshToken)
	return session
}

// IsExpired checks if the session token has expired.
func (s *Session) IsExpired(currentTime int64) bool {
	return (s.ExpiresAt - currentTime) < 0
}

// IsRefreshExpired checks if the refresh token has expired.
func (s *Session) IsRefreshExpired(currentTime int64) bool {
	return (s.RefreshExpiresAt - currentTime) < 0
}

// Update updates the session with a new token and refresh token.
func (s *Session) Update(token, refreshToken string) error {
	tokenDecoded, err := s.decodeJWT(token)
	if err != nil {
		return err
	}

	exp, err := parseInt64FromMap(tokenDecoded, "exp")
	if err != nil {
		return err
	}
	s.ExpiresAt = exp

	s.Token = token
	if username, ok := tokenDecoded["usn"].(string); ok {
		s.Username = username
	}
	if userID, ok := tokenDecoded["uid"].(string); ok {
		s.UserID = userID
	}
	if vars, ok := tokenDecoded["vrs"].(map[string]interface{}); ok {
		s.Vars = vars
	}

	// Handle refresh token
	if refreshToken != "" {
		refreshTokenDecoded, err := s.decodeJWT(refreshToken)
		if err != nil {
			return err
		}

		refreshExp, err := parseInt64FromMap(refreshTokenDecoded, "exp")
		if err != nil {
			return err
		}
		s.RefreshExpiresAt = refreshExp
		s.RefreshToken = refreshToken
	}

	return nil
}

// decodeJWT decodes a JWT token and returns its payload as a map.
func (s *Session) decodeJWT(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) < 2 {
		return nil, errors.New("invalid token format")
	}

	base64Raw := parts[1]
	base64Str := strings.ReplaceAll(strings.ReplaceAll(base64Raw, "-", "+"), "_", "/")
	decoded, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return nil, err
	}

	var payload map[string]interface{}
	if err := json.Unmarshal(decoded, &payload); err != nil {
		return nil, err
	}

	return payload, nil
}

// parseInt64FromMap parses an int64 value from a map by key.
func parseInt64FromMap(data map[string]interface{}, key string) (int64, error) {
	value, ok := data[key]
	if !ok {
		return 0, errors.New("key not found in map")
	}

	floatValue, ok := value.(float64)
	if !ok {
		return 0, errors.New("value is not a number")
	}

	return int64(floatValue), nil
}

// Restore creates a Session from an existing token and refresh token.
// It assumes that the session is not newly created.
func Restore(token, refreshToken string) *Session {
	return NewSession(token, refreshToken, false)
}
