package nakama

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gwaylib/errors"
	api "github.com/heroiclabs/nakama-common/api"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	ErrNoContent = errors.New("No content by 204")
)

type NakamaApi struct {
	ServerKey string
	BasePath  string
	TimeoutMs int // need set a validate value
}

func (napi NakamaApi) SetBasicAuth(req *http.Request, username, passwd string) {
	if checkStr(&username) {
		auth := username + ":"
		if checkStr(&passwd) {
			auth += passwd
		}
		encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
		req.Header.Set("Authorization", "Basic "+encodedAuth)
	}
}

func (napi *NakamaApi) doReq(bearerToken string, req *http.Request, options map[string]string, rsp proto.Message) error {
	if checkStr(&bearerToken) {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}
	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}

	// Run the HTTP request in a goroutine
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return errors.As(err)
	}
	defer resp.Body.Close()

	// Handle HTTP response
	if resp.StatusCode == http.StatusNoContent {
		if rsp != nil {
			return ErrNoContent.As(resp.StatusCode)
		}
		return nil
	} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return errors.As(err, string(bodyBytes))
		}
		if rsp == nil {
			return nil
		}

		if err := protojson.Unmarshal(bodyBytes, rsp); err != nil {
			return errors.As(err)
		}
		return nil
	}
	return errors.New(resp.Status).As(resp.StatusCode)
}

// Healthcheck is a healthcheck function that load balancers can use to check the service.
func (napi *NakamaApi) Healthcheck(bearerToken string, options map[string]string) error {
	// Define the URL path and query parameters
	urlPath := "/healthcheck"
	queryParams := url.Values{}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return errors.As(err)
	}
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

// DeleteAccount deletes the current user's account.
func (napi *NakamaApi) DeleteAccount(bearerToken string, options map[string]string) error {
	// Define the URL path and query parameters
	urlPath := "/v2/account"
	queryParams := url.Values{}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("DELETE", fullUrl, nil)
	if err != nil {
		return err
	}
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

// GetAccount fetches the current user's account.
func (napi *NakamaApi) GetAccount(bearerToken string, options map[string]string) (*api.Account, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account"
	queryParams := url.Values{}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	result := &api.Account{}
	if err := napi.doReq(bearerToken, req, options, result); err != nil {
		return nil, errors.As(err)
	}
	return result, nil
}

// UpdateAccount updates fields in the current user's account.
func (napi *NakamaApi) UpdateAccount(bearerToken string, body *api.UpdateAccountRequest, options map[string]string) error {
	// Check if the body is nil
	if body == nil {
		return errors.New("'body' is a required parameter but is null or undefined")
	}

	// Define the URL path and query parameters
	urlPath := "/v2/account"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("PUT", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}

	return nil
}

// AuthenticateApple authenticates a user with an Apple ID against the server.
func (napi *NakamaApi) AuthenticateApple(basicAuthUsername string, basicAuthPassword string, account *api.AccountApple, create *bool, username string, options map[string]string) (*api.Session, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/authenticate/apple"
	queryParams := url.Values{}
	if create != nil {
		queryParams.Set("create", fmt.Sprintf("%v", *create))
	}
	if username != "" {
		queryParams.Set("username", username)
	}

	// Convert the account to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Authorization header
	napi.SetBasicAuth(req, basicAuthUsername, basicAuthPassword)

	var result = &api.Session{}
	if err := napi.doReq("", req, options, result); err != nil {
		return nil, errors.As(err)
	}

	return result, nil
}

// AuthenticateCustom authenticates a user with a custom ID against the server.
func (napi *NakamaApi) AuthenticateCustom(
	basicAuthUsername string,
	basicAuthPassword string,
	account *api.AccountCustom,
	create *bool,
	username *string,
	options map[string]string,
) (*api.Session, error) {
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
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Authorization header
	napi.SetBasicAuth(req, basicAuthUsername, basicAuthPassword)
	var result api.Session
	if err := napi.doReq("", req, options, &result); err != nil {
		return nil, errors.As(err)
	}

	return &result, nil
}

// AuthenticateDevice authenticates a user with a device ID against the server.
func (napi *NakamaApi) AuthenticateDevice(
	basicAuthUsername string,
	basicAuthPassword string,
	account *api.AccountDevice,
	create *bool,
	username string,
	options map[string]string,
) (*api.Session, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/authenticate/device"
	queryParams := url.Values{}
	if create != nil {
		queryParams.Set("create", fmt.Sprintf("%v", *create))
	}
	if username != "" {
		queryParams.Set("username", username)
	}

	// Convert the account to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Authorization header
	napi.SetBasicAuth(req, basicAuthUsername, basicAuthPassword)

	var result api.Session
	if err := napi.doReq("", req, options, &result); err != nil {
		return nil, errors.As(err)
	}

	return &result, nil
}

// AuthenticateEmail authenticates a user with an email and password against the server.
func (napi *NakamaApi) AuthenticateEmail(
	basicAuthUsername string,
	basicAuthPassword string,
	account *api.AccountEmail,
	create *bool,
	username *string,
	options map[string]string,
) (*api.Session, error) {
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
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Authorization header
	napi.SetBasicAuth(req, basicAuthUsername, basicAuthPassword)

	var result = &api.Session{}
	if err := napi.doReq("", req, options, result); err != nil {
		return nil, errors.As(err)
	}

	return result, nil
}

// AuthenticateFacebook authenticates a user with a Facebook OAuth token against the server.
func (napi *NakamaApi) AuthenticateFacebook(
	basicAuthUsername string,
	basicAuthPassword string,
	account *api.AccountFacebook,
	create *bool,
	username string,
	sync *bool,
	options map[string]string,
) (*api.Session, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/authenticate/facebook"
	queryParams := url.Values{}
	if create != nil {
		queryParams.Set("create", fmt.Sprintf("%v", create))
	}
	if username != "" {
		queryParams.Set("username", username)
	}
	if sync != nil {
		queryParams.Set("sync", fmt.Sprintf("%v", sync))
	}

	// Convert the account to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Authorization header
	napi.SetBasicAuth(req, basicAuthUsername, basicAuthPassword)
	var result api.Session
	if err := napi.doReq("", req, options, &result); err != nil {
		return nil, errors.As(err)
	}

	return &result, nil
}

// AuthenticateFacebookInstantGame authenticates a user with a Facebook Instant Game token against the server.
func (napi *NakamaApi) AuthenticateFacebookInstantGame(
	basicAuthUsername string,
	basicAuthPassword string,
	account *api.AccountFacebookInstantGame,
	create *bool,
	username string,
	options map[string]string,
) (*api.Session, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/authenticate/facebookinstantgame"
	queryParams := url.Values{}
	if create != nil {
		queryParams.Set("create", fmt.Sprintf("%v", *create))
	}
	if username != "" {
		queryParams.Set("username", username)
	}

	// Convert the account to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Authorization header
	napi.SetBasicAuth(req, basicAuthUsername, basicAuthPassword)

	var result api.Session
	if err := napi.doReq("", req, options, &result); err != nil {
		return nil, errors.As(err)
	}

	return &result, nil
}

// AuthenticateGameCenter authenticates a user with Apple's GameCenter against the server.
func (napi *NakamaApi) AuthenticateGameCenter(
	basicAuthUsername string,
	basicAuthPassword string,
	account *api.AccountGameCenter,
	create *bool,
	username string,
	options map[string]string,
) (*api.Session, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/authenticate/gamecenter"
	queryParams := url.Values{}
	if create != nil {
		queryParams.Set("create", fmt.Sprintf("%v", *create))
	}
	if username != "" {
		queryParams.Set("username", username)
	}

	// Convert the account to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Authorization header
	napi.SetBasicAuth(req, basicAuthUsername, basicAuthPassword)

	var result api.Session
	if err := napi.doReq("", req, options, &result); err != nil {
		return nil, errors.As(err)
	}

	return &result, nil
}

// AuthenticateGoogle authenticates a user with Google against the server.
func (napi *NakamaApi) AuthenticateGoogle(
	basicAuthUsername string,
	basicAuthPassword string,
	account *api.AccountGoogle,
	create *bool,
	username string,
	options map[string]string,
) (*api.Session, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/authenticate/google"
	queryParams := url.Values{}
	queryParams.Set("create", fmt.Sprintf("%v", create))
	if username != "" {
		queryParams.Set("username", username)
	}

	// Convert the account to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Authorization header
	napi.SetBasicAuth(req, basicAuthUsername, basicAuthPassword)

	var result api.Session
	if err := napi.doReq("", req, options, &result); err != nil {
		return nil, errors.As(err)
	}

	return &result, nil

}

// AuthenticateSteam authenticates a user with Steam against the server.
func (napi *NakamaApi) AuthenticateSteam(
	basicAuthUsername string,
	basicAuthPassword string,
	account *api.AccountSteam,
	create *bool,
	username string,
	sync *bool,
	options map[string]string,
) (*api.Session, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/authenticate/steam"
	queryParams := url.Values{}
	if create != nil {
		queryParams.Set("create", fmt.Sprintf("%v", create))
	}
	if username != "" {
		queryParams.Set("username", username)
	}
	if sync != nil {
		queryParams.Set("sync", fmt.Sprintf("%v", sync))
	}

	// Convert the account to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return nil, errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Authorization header
	napi.SetBasicAuth(req, basicAuthUsername, basicAuthPassword)

	var result api.Session
	if err := napi.doReq("", req, options, &result); err != nil {
		return nil, errors.As(err)
	}

	return &result, nil
}

// LinkApple adds an Apple ID to the social profiles on the current user's account.
func (napi *NakamaApi) LinkApple(
	bearerToken string,
	body *api.AccountApple,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/account/link/apple"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

// LinkCustom adds a custom ID to the social profiles on the current user's account.
func (napi *NakamaApi) LinkCustom(
	bearerToken string,
	body *api.AccountCustom,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/account/link/custom"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

// LinkDevice adds a device ID to the social profiles on the current user's account.
func (napi *NakamaApi) LinkDevice(
	bearerToken string,
	body *api.AccountDevice,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/account/link/device"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

// LinkEmail adds an email and password to the social profiles on the current user's account.
func (napi *NakamaApi) LinkEmail(
	bearerToken string,
	body *api.AccountEmail,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/account/link/email"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

// LinkFacebook adds a Facebook account to the social profiles on the current user's account.
func (napi *NakamaApi) LinkFacebook(
	bearerToken string,
	account *api.AccountFacebook,
	sync *bool,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/account/link/facebook"
	queryParams := url.Values{}
	if sync != nil {
		queryParams.Set("sync", fmt.Sprintf("%t", *sync))
	}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return err
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil

}

// LinkFacebookInstantGame adds a Facebook Instant Game account to the social profiles on the current user's account.
func (napi *NakamaApi) LinkFacebookInstantGame(
	bearerToken string,
	body *api.AccountFacebookInstantGame,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/account/link/facebookinstantgame"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}

	return nil
}

// LinkGameCenter adds Apple's GameCenter to the social profiles on the current user's account.
func (napi *NakamaApi) LinkGameCenter(
	bearerToken string,
	body *api.AccountGameCenter,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/account/link/gamecenter"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil

}

// LinkGoogle adds a Google account to the social profiles on the current user's account.
func (napi *NakamaApi) LinkGoogle(
	bearerToken string,
	body *api.AccountGoogle,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/account/link/google"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

// LinkSteam adds a Steam account to the social profiles on the current user's account.
func (napi *NakamaApi) LinkSteam(
	bearerToken string,
	body *api.LinkSteamRequest,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/account/link/steam"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

// SessionRefresh refreshes a user's session using a refresh token retrieved from a previous authentication request.
func (napi *NakamaApi) SessionRefresh(
	basicAuthUsername string,
	basicAuthPassword string,
	body *api.SessionRefreshRequest,
	options map[string]string,
) (*api.Session, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/account/session/refresh"
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Basic Auth header
	napi.SetBasicAuth(req, basicAuthUsername, basicAuthPassword)

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
			var result api.Session
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return &result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// UnlinkApple removes the Apple ID from the social profiles on the current user's account.
func (napi *NakamaApi) UnlinkApple(
	bearerToken string,
	body *api.AccountApple,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/account/unlink/apple"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

// UnlinkCustom removes the custom ID from the social profiles on the current user's account.
func (napi *NakamaApi) UnlinkCustom(
	bearerToken string,
	body *api.AccountCustom,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/account/unlink/custom"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")

	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil

}

// UnlinkDevice removes the device ID from the social profiles on the current user's account.
func (napi *NakamaApi) UnlinkDevice(
	bearerToken string,
	body *api.AccountDevice,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/account/unlink/device"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

// UnlinkEmail removes the email+password from the social profiles on the current user's account.
func (napi *NakamaApi) UnlinkEmail(
	bearerToken string,
	body *api.AccountEmail,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/account/unlink/email"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

// UnlinkFacebook removes the Facebook profile from the social profiles on the current user's account.
func (napi *NakamaApi) UnlinkFacebook(
	bearerToken string,
	body *api.AccountFacebook,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/account/unlink/facebook"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

// UnlinkFacebookInstantGame removes the Facebook Instant Game profile from the social profiles on the current user's account.
func (napi *NakamaApi) UnlinkFacebookInstantGame(
	bearerToken string,
	body *api.AccountFacebookInstantGame,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/account/unlink/facebookinstantgame"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

// UnlinkGameCenter removes the GameCenter profile from the social profiles on the current user's account.
func (napi *NakamaApi) UnlinkGameCenter(
	bearerToken string,
	body *api.AccountGameCenter,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/account/unlink/gamecenter"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

// UnlinkGoogle removes the Google profile from the social profiles on the current user's account.
func (napi *NakamaApi) UnlinkGoogle(
	bearerToken string,
	body *api.AccountGoogle,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/account/unlink/google"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

// UnlinkSteam removes the Steam profile from the social profiles on the current user's account.
func (napi *NakamaApi) UnlinkSteam(
	bearerToken string,
	body *api.AccountSteam,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/account/unlink/steam"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

// ListChannelMessages lists a channel's message history.
func (napi *NakamaApi) ListChannelMessages(
	bearerToken *string,
	channelId *string,
	limit *int,
	forward *bool,
	cursor *string,
	options map[string]string,
) (*api.ChannelMessageList, error) {
	if !checkStr(channelId) {
		return nil, errors.New("'channelId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/channel/{channelId}", "{channelId}", url.PathEscape(*channelId), 1)
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
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
			return &api.ChannelMessageList{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result api.ChannelMessageList
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return &result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// Event submits an event for processing in the server's registered runtime custom events handler.
func (napi *NakamaApi) Event(
	bearerToken *string,
	body *api.Event,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/event"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != nil && *bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
		return errors.As(err)
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil
		} else {
			return errors.New(resp.Status).As(resp.StatusCode)
		}
	}
}

func (napi *NakamaApi) DeleteFriends(
	bearerToken *string,
	ids []string,
	usernames []string,
	options map[string]string,
) error {
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
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("DELETE", fullUrl, nil)
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != nil && *bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
		return errors.As(err)
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

// ListFriends fetches the list of all friends for the current user.
func (napi *NakamaApi) ListFriends(
	bearerToken *string,
	limit *int,
	state *int,
	cursor *string,
	options map[string]string,
) (*api.FriendList, error) {
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
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return &api.FriendList{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
		return &api.FriendList{}, errors.New("request timed out")
	case err := <-errorChan:
		return &api.FriendList{}, err
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return &api.FriendList{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result api.FriendList
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return &api.FriendList{}, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return &api.FriendList{}, err
			}
			return &result, nil
		} else {
			return &api.FriendList{}, errors.New(resp.Status)
		}
	}
}

func (napi *NakamaApi) AddFriends(
	bearerToken *string,
	ids []string,
	usernames []string,
	options map[string]string,
) error {
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
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken == nil && *bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
		return errors.As(err)
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

func (napi *NakamaApi) BlockFriends(
	bearerToken *string,
	ids []string,
	usernames []string,
	options map[string]string,
) error {
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
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != nil && *bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
		return errors.As(err)
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

func (napi *NakamaApi) ImportFacebookFriends(
	bearerToken *string,
	account *api.AccountFacebook,
	reset *bool,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/friend/facebook"
	queryParams := url.Values{}
	if checkBool(reset) {
		queryParams.Set("reset", strconv.FormatBool(*reset))
	}

	// Serialize the account object to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
		return errors.As(err)
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

func (napi *NakamaApi) ListFriendsOfFriends(
	bearerToken *string,
	limit *int,
	cursor *string,
	options map[string]string,
) (*api.FriendsOfFriendsList, error) {
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
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
			var result api.FriendsOfFriendsList
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return &result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

func (napi *NakamaApi) ImportSteamFriends(
	bearerToken *string,
	account *api.AccountSteam,
	reset *bool,
	options map[string]string,
) error {
	// Define the URL path and query parameters
	urlPath := "/v2/friend/steam"
	queryParams := url.Values{}
	if checkBool(reset) {
		queryParams.Set("reset", strconv.FormatBool(*reset))
	}

	// Serialize the account object to JSON
	bodyJson, err := json.Marshal(account)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
		return errors.As(err)
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

func (napi *NakamaApi) ListGroups(
	bearerToken *string,
	name *string,
	cursor *string,
	limit *int,
	langTag *string,
	members *int,
	open *bool,
	options map[string]string,
) (*api.GroupList, error) {
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
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
			var result *api.GroupList
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
func (napi *NakamaApi) CreateGroup(
	bearerToken *string,
	body *api.CreateGroupRequest,
	options map[string]string,
) (*api.Group, error) {
	// Define the URL path and query parameters
	urlPath := "/v2/group"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != nil && *bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
	defer cancel()

	// Make the HTTP request
	client := &http.Client{}
	responseChan := make(chan *http.Response, 1)
	errorChan := make(chan error, 1)

	// Run the HTTP request in a goroutine
	go func() {
		resp, err := client.Do(req.WithContext(ctx))
		if err != nil {
			errorChan <- errors.As(err)
			return
		}
		responseChan <- resp
	}()

	// Wait for the response or the timeout
	select {
	case <-ctx.Done():
		return nil, errors.New("request timed out")
	case err := <-errorChan:
		return nil, errors.As(err)
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, ErrNoContent.As(resp.StatusCode)
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, errors.As(err)
			}
			var result api.Group
			if err = json.Unmarshal(bodyBytes, &result); err != nil {
				return nil, errors.As(err)
			}
			return &result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// DeleteGroup deletes a group by ID.
func (napi *NakamaApi) DeleteGroup(
	bearerToken *string,
	groupId *string,
	options map[string]string,
) error {
	if !checkStr(groupId) {
		return errors.New("'groupId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}", "{groupId}", url.QueryEscape(*groupId), 1)
	queryParams := url.Values{}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("DELETE", fullUrl, nil)
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
		return errors.As(err)
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

// UpdateGroup updates fields in a given group.
func (napi *NakamaApi) UpdateGroup(
	bearerToken string,
	groupId *string,
	body *api.UpdateGroupRequest,
	options map[string]string,
) error {
	// Validate required parameters
	if checkStr(groupId) {
		return errors.New("'groupId' is a required parameter but is empty")
	}
	if body == nil {
		return errors.New("'body' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}", "{groupId}", url.QueryEscape(*groupId), 1)
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("PUT", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")

	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}

	return nil

}

// AddGroupUsers adds users to a group.
func (napi *NakamaApi) AddGroupUsers(
	bearerToken *string,
	groupId *string,
	userIds []string,
	options map[string]string,
) error {

	// Check required parameters
	if !checkStr(groupId) {
		return errors.New("'groupId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}/add", "{groupId}", url.QueryEscape(*groupId), 1)
	queryParams := url.Values{}
	for _, userId := range userIds {
		queryParams.Add("user_ids", userId)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != nil && *bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
		return errors.As(err)
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

// BanGroupUsers bans a set of users from a group.
func (napi *NakamaApi) BanGroupUsers(
	bearerToken *string,
	groupId *string,
	userIds []string,
	options map[string]string,
) error {
	// Check required parameters
	if !checkStr(groupId) {
		return errors.New("'groupId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}/ban", "{groupId}", url.QueryEscape(*groupId), 1)
	queryParams := url.Values{}
	for _, userId := range userIds {
		queryParams.Add("user_ids", userId)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if bearerToken != nil && *bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
		return errors.As(err)
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

// DemoteGroupUsers demotes a set of users in a group to the next role down.
func (napi *NakamaApi) DemoteGroupUsers(
	bearerToken *string,
	groupId *string,
	userIds []string,
	options map[string]string,
) error {
	// Check required parameters
	if groupId == nil || *groupId == "" {
		return errors.New("'groupId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}/demote", "{groupId}", url.QueryEscape(*groupId), 1)
	queryParams := url.Values{}
	for _, userId := range userIds {
		queryParams.Add("user_ids", userId)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
		return errors.As(err)
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

// JoinGroup immediately joins an open group, or requests to join a closed one.
func (napi *NakamaApi) JoinGroup(
	bearerToken *string,
	groupId *string,
	options map[string]string,
) error {
	if !checkStr(groupId) {
		return errors.New("'groupId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}/join", "{groupId}", url.QueryEscape(*groupId), 1)
	queryParams := url.Values{}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
		return errors.As(err)
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

// KickGroupUsers kicks a set of users from a group.
func (napi *NakamaApi) KickGroupUsers(
	bearerToken *string,
	groupId *string,
	userIds []string,
	options map[string]string,
) error {

	// Validate required parameter
	if !checkStr(groupId) {
		return errors.New("'groupId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}/kick", "{groupId}", url.QueryEscape(*groupId), 1)
	queryParams := url.Values{}
	for _, userId := range userIds {
		queryParams.Add("user_ids", userId)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
		return errors.As(err)
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

// LeaveGroup allows a user to leave a group they are a member of.
func (napi *NakamaApi) LeaveGroup(
	bearerToken *string,
	groupId *string,
	options map[string]string,
) error {
	// Validate the required parameter
	if checkStr(groupId) {
		return errors.New("'groupId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}/leave", "{groupId}", url.QueryEscape(*groupId), 1)
	queryParams := url.Values{}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
		return errors.As(err)
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

// PromoteGroupUsers promotes a set of users in a group to the next role up.
func (napi *NakamaApi) PromoteGroupUsers(
	bearerToken string,
	groupId string,
	userIds []string,
	options map[string]string,
) error {
	// Validate required parameter
	if checkStr(&groupId) {
		return errors.New("'groupId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}/promote", "{groupId}", url.QueryEscape(groupId), 1)
	queryParams := url.Values{}
	for _, userId := range userIds {
		queryParams.Add("user_ids", userId)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, nil)
	if err != nil {
		return errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")

	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}

	return nil

}

// ListGroupUsers lists all users that are part of a group.
func (napi *NakamaApi) ListGroupUsers(
	bearerToken *string,
	groupId *string,
	limit *int,
	state *int,
	cursor *string,
	options map[string]string,
) (*api.GroupUserList, error) {
	// Validate the required parameter
	if checkStr(groupId) {
		return nil, errors.New("'groupId' is a required parameter but is empty")
	}

	// Define the URL path and query parameters
	urlPath := strings.Replace("/v2/group/{groupId}/user", "{groupId}", url.QueryEscape(*groupId), 1)
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
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
			var result api.GroupUserList
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return &result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

func (napi *NakamaApi) ValidatePurchaseApple(
	bearerToken *string,
	body *api.ValidatePurchaseAppleRequest,
	options map[string]string,
) (*api.ValidatePurchaseResponse, error) {
	// Define the URL path
	urlPath := "/v2/iap/purchase/apple"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
			var result api.ValidatePurchaseResponse
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return &result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// ValidatePurchaseFacebookInstant validates an Instant IAP receipt from Facebook.
func (napi *NakamaApi) ValidatePurchaseFacebookInstant(
	bearerToken *string,
	body *api.ValidatePurchaseFacebookInstantRequest,
	options map[string]string,
) (*api.ValidatePurchaseResponse, error) {
	// Validate the required parameter
	if body == nil {
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
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
			var result api.ValidatePurchaseResponse
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return &result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// ValidatePurchaseGoogle validates an IAP receipt from Google.
func (napi *NakamaApi) ValidatePurchaseGoogle(
	bearerToken *string,
	body *api.ValidatePurchaseGoogleRequest,
	options map[string]string,
) (*api.ValidatePurchaseResponse, error) {
	// Validate the required parameter
	if body == nil {
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
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
			var result api.ValidatePurchaseResponse
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return &result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// ValidatePurchaseHuawei validates an IAP receipt from Huawei.
func (napi *NakamaApi) ValidatePurchaseHuawei(
	bearerToken *string,
	body *api.ValidatePurchaseHuaweiRequest,
	options map[string]string,
) (*api.ValidatePurchaseResponse, error) {
	// Validate the required parameter
	if body == nil {
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
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
			var result api.ValidatePurchaseResponse
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return &result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// ListSubscriptions lists user's subscriptions.
func (napi *NakamaApi) ListSubscriptions(
	bearerToken *string,
	body *api.ListSubscriptionsRequest,
	options map[string]string,
) (*api.SubscriptionList, error) {

	// Validate the required parameter
	if body == nil {
		return nil, errors.New("'body' is a required parameter but is null or undefined.")
	}

	// Define the URL path
	urlPath := "/v2/iap/subscription"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(bodyJson))
	if err != nil {
		return nil, errors.As(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
		return nil, errors.As(err)
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, ErrNoContent.As(resp.StatusCode)
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result api.SubscriptionList
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, errors.As(err)
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, errors.As(err)
			}
			return &result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// ValidateSubscriptionApple validates an Apple subscription receipt.
func (napi *NakamaApi) ValidateSubscriptionApple(
	bearerToken *string,
	body *api.ValidateSubscriptionAppleRequest,
	options map[string]string,
) (*api.ValidateSubscriptionResponse, error) {

	// Validate the required parameter
	if body == nil {
		return nil, errors.New("'body' is a required parameter but is null or undefined.")
	}

	// Define the URL path
	urlPath := "/v2/iap/subscription/apple"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
			return &api.ValidateSubscriptionResponse{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result api.ValidateSubscriptionResponse
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return &result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// ValidateSubscriptionGoogle validates a Google subscription receipt.
func (napi *NakamaApi) ValidateSubscriptionGoogle(
	bearerToken *string,
	body *api.ValidateSubscriptionGoogleRequest,
	options map[string]string,
) (*api.ValidateSubscriptionResponse, error) {

	// Validate the required parameter
	if body == nil {
		return nil, errors.New("'body' is a required parameter but is null or undefined.")
	}

	// Define the URL path
	urlPath := "/v2/iap/subscription/google"
	queryParams := url.Values{}

	// Serialize the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
			return &api.ValidateSubscriptionResponse{}, nil
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result api.ValidateSubscriptionResponse
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			err = json.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, err
			}
			return &result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// GetSubscription retrieves a subscription by product ID.
func (napi *NakamaApi) GetSubscription(
	bearerToken *string,
	productId *string,
	options map[string]string,
) (*api.ValidatedSubscription, error) {

	// Validate the required parameter
	if productId == nil || *productId == "" {
		return nil, errors.New("'productId' is a required parameter but is null or empty.")
	}

	// Define the URL path
	urlPath := fmt.Sprintf("/v2/iap/subscription/%s", url.QueryEscape(*productId))
	queryParams := url.Values{}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
		return nil, errors.As(err)
	case resp := <-responseChan:
		defer resp.Body.Close()

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, ErrNoContent.As(resp.StatusCode)
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			var result api.ValidatedSubscription
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, errors.As(err)
			}
			err = protojson.Unmarshal(bodyBytes, &result)
			if err != nil {
				return nil, errors.As(err)
			}
			return &result, nil
		} else {
			return nil, errors.New(resp.Status)
		}
	}
}

// DeleteLeaderboardRecord deletes a leaderboard record.
func (napi *NakamaApi) DeleteLeaderboardRecord(
	bearerToken *string,
	leaderboardId *string,
	options map[string]string,
) error {

	// Validate the required parameter
	if !checkStr(leaderboardId) {
		return errors.New("'leaderboardId' is a required parameter but is null or empty.")
	}

	// Define the URL path
	urlPath := fmt.Sprintf("/v2/leaderboard/%s", url.QueryEscape(*leaderboardId))
	queryParams := url.Values{}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("DELETE", fullUrl, nil)
	if err != nil {
		return err
	}

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
func (napi *NakamaApi) ListLeaderboardRecords(
	bearerToken *string,
	leaderboardId *string,
	ownerIds []string,
	limit *int,
	cursor *string,
	expiry *string,
	options map[string]string,
) (*api.LeaderboardRecordList, error) {

	// Validate the required parameter
	if !checkStr(leaderboardId) {
		return nil, errors.New("'leaderboardId' is a required parameter but is null or empty.")
	}

	// Define the URL path
	urlPath := fmt.Sprintf("/v2/leaderboard/%s", url.QueryEscape(*leaderboardId))
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
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, errors.As(err)
	}

	// Set Bearer Token authorization header
	if checkStr(bearerToken) {
		req.Header.Set("Authorization", "Bearer "+*bearerToken)
	}

	// Apply additional custom headers or options if needed
	for key, value := range options {
		req.Header.Set(key, value)
	}

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(napi.TimeoutMs)*time.Millisecond)
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
		return nil, errors.As(err)
	case resp := <-responseChan:
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, errors.As(err)
		}

		// Handle HTTP response
		if resp.StatusCode == http.StatusNoContent {
			return nil, ErrNoContent.As(resp.StatusCode)
		} else if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			result := &api.LeaderboardRecordList{}
			fmt.Println(string(bodyBytes))
			err = protojson.Unmarshal(bodyBytes, result)
			if err != nil {
				return nil, errors.As(err)
			}
			return result, nil
		} else {
			return nil, errors.New(resp.Status).As(string(bodyBytes))
		}
	}
}

// WriteLeaderboardRecord writes a record to a leaderboard.
func (napi *NakamaApi) WriteLeaderboardRecord(
	bearerToken string,
	leaderboardId string,
	record *api.WriteLeaderboardRecordRequest_LeaderboardRecordWrite,
	options map[string]string,
) (*api.LeaderboardRecord, error) {

	// Validate the required parameters
	if !checkStr(&leaderboardId) {
		return nil, errors.New("'leaderboardId' is a required parameter but is null or empty.")
	}
	if record == nil {
		return nil, errors.New("'record' is a required parameter but is null or empty.")
	}

	// Define the URL path
	urlPath := fmt.Sprintf("/v2/leaderboard/%s", url.QueryEscape(leaderboardId))
	queryParams := url.Values{}

	// Convert the record to JSON
	bodyJson, err := json.Marshal(record)
	if err != nil {
		return nil, errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, bytes.NewBuffer(bodyJson))
	if err != nil {
		return nil, errors.As(err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "application/json")

	result := &api.LeaderboardRecord{}
	if err := napi.doReq(bearerToken, req, options, result); err != nil {
		return nil, errors.As(err)
	}
	return result, nil
}

// ListLeaderboardRecordsAroundOwner lists leaderboard records that belong to a user.
func (napi *NakamaApi) ListLeaderboardRecordsAroundOwner(
	bearerToken string,
	leaderboardId string,
	ownerId string,
	limit int,
	expiry string,
	cursor string,
	options map[string]string,
) (*api.LeaderboardRecordList, error) {

	// Validate the required parameters
	if !checkStr(&leaderboardId) {
		return nil, errors.New("'leaderboardId' is a required parameter but is null or empty.")
	}
	if !checkStr(&ownerId) {
		return nil, errors.New("'ownerId' is a required parameter but is null or empty.")
	}

	// Define the URL path
	urlPath := fmt.Sprintf(
		"/v2/leaderboard/%s/owner/%s",
		url.QueryEscape(leaderboardId),
		url.QueryEscape(ownerId),
	)
	queryParams := url.Values{}

	// Add optional parameters to the query
	if limit > 0 {
		queryParams.Set("limit", fmt.Sprintf("%d", limit))
	}
	if expiry != "" {
		queryParams.Set("expiry", expiry)
	}
	if cursor != "" {
		queryParams.Set("cursor", cursor)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, errors.As(err)
	}

	result := &api.LeaderboardRecordList{}
	if err := napi.doReq(bearerToken, req, options, result); err != nil {
		return nil, errors.As(err)
	}

	return result, nil
}

func (napi *NakamaApi) ListMatches(
	bearerToken string,
	limit int,
	authoritative *bool,
	label string,
	minSize int,
	maxSize int,
	query string,
	options map[string]string,
) (*api.MatchList, error) {

	// Define the URL path
	urlPath := "/v2/match"
	queryParams := url.Values{}

	// Add optional parameters to the query
	if limit > 0 {
		queryParams.Set("limit", fmt.Sprintf("%d", limit))
	}
	if authoritative != nil {
		queryParams.Set("authoritative", fmt.Sprintf("%t", *authoritative))
	}
	if label != "" {
		queryParams.Set("label", label)
	}
	if minSize > 0 {
		queryParams.Set("min_size", fmt.Sprintf("%d", minSize))
	}
	if maxSize > 0 {
		queryParams.Set("max_size", fmt.Sprintf("%d", maxSize))
	}
	if query != "" {
		queryParams.Set("query", query)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}

	var result api.MatchList
	if err := napi.doReq(bearerToken, req, options, &result); err != nil {
		return nil, errors.As(err)
	}

	return &result, nil
}

func (napi *NakamaApi) DeleteNotifications(
	bearerToken string,
	ids []string,
	options map[string]string,
) error {

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
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("DELETE", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

func (napi *NakamaApi) ListNotifications(
	bearerToken string,
	limit int,
	cacheableCursor string,
	options map[string]string,
) (*api.NotificationList, error) {

	// Define the URL path
	urlPath := "/v2/notification"
	queryParams := url.Values{}

	// Add query parameters
	if limit > 0 {
		queryParams.Set("limit", fmt.Sprintf("%d", limit))
	}
	if cacheableCursor != "" {
		queryParams.Set("cacheable_cursor", cacheableCursor)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, errors.As(err)
	}

	result := &api.NotificationList{}
	if err := napi.doReq(bearerToken, req, options, result); err != nil {
		return nil, errors.As(err)
	}

	return result, nil
}

func (napi *NakamaApi) RpcFunc2(
	bearerToken string,
	id string,
	payload string,
	httpKey string,
	options map[string]string,
) (*api.Rpc, error) {

	// Validate the required parameter 'id'
	if !checkStr(&id) {
		return nil, errors.New("'id' is a required parameter but is empty")
	}

	// Define the URL path
	urlPath := fmt.Sprintf("/v2/rpc/%s", url.QueryEscape(id))

	// Add query parameters
	queryParams := url.Values{}
	if payload != "" {
		queryParams.Set("payload", payload)
	}
	if httpKey != "" {
		queryParams.Set("http_key", httpKey)
	}

	// Convert the body to JSON (if necessary, can be modified)
	bodyJson := ""

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return nil, errors.As(err)
	}

	result := &api.Rpc{}
	if err := napi.doReq(bearerToken, req, options, result); err != nil {
		return nil, errors.As(err)
	}
	return result, nil
}

func (napi *NakamaApi) RpcFunc(
	bearerToken string,
	id string,
	body string,
	httpKey string,
	options map[string]string,
) (*api.Rpc, error) {
	// Validate the required parameters 'id' and 'body'
	if !checkStr(&id) {
		return nil, errors.New("'id' is a required parameter but is empty")
	}
	if !checkStr(&body) {
		return nil, errors.New("'body' is a required parameter but is empty")
	}

	// Define the URL path
	urlPath := fmt.Sprintf("/v2/rpc/%s", url.QueryEscape(id))

	// Add query parameters
	queryParams := url.Values{}
	if httpKey != "" {
		queryParams.Set("http_key", httpKey)
	}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, strings.NewReader(string(bodyJson)))
	if err != nil {
		return nil, errors.As(err)
	}

	result := &api.Rpc{}
	if err := napi.doReq(bearerToken, req, options, result); err != nil {
		return nil, errors.As(err)
	}

	return result, nil
}

func (napi *NakamaApi) SessionLogout(
	bearerToken string,
	body *api.SessionLogoutRequest,
	options map[string]string,
) error {
	// Validate the required parameter 'body'
	if body == nil {
		return errors.New("'body' is a required parameter but is null or undefined")
	}

	// Define the URL path
	urlPath := "/v2/session/logout"

	// Add query parameters (empty for this request)
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, strings.NewReader(string(bodyJson)))
	if err != nil {
		return errors.As(err)
	}

	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

func (napi *NakamaApi) ReadStorageObjects(
	bearerToken string,
	body *api.ReadStorageObjectsRequest,
	options map[string]string,
) (*api.StorageObjects, error) {
	// Define the URL path
	urlPath := "/v2/storage"

	// Add query parameters (empty for this request)
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, strings.NewReader(string(bodyJson)))
	if err != nil {
		return nil, errors.As(err)
	}

	result := &api.StorageObjects{}
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return nil, errors.As(err)
	}

	return result, nil
}

func (napi *NakamaApi) WriteStorageObjects(
	bearerToken string,
	body *api.WriteStorageObjectsRequest,
	options map[string]string,
) (*api.StorageObjectAcks, error) {
	// Define the URL path
	urlPath := "/v2/storage"

	// Add query parameters (empty for this request)
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return nil, errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("PUT", fullUrl, strings.NewReader(string(bodyJson)))
	if err != nil {
		return nil, errors.As(err)
	}

	var result api.StorageObjectAcks
	if err := napi.doReq(bearerToken, req, options, &result); err != nil {
		return nil, errors.As(err)
	}

	return &result, nil
}

func (napi *NakamaApi) DeleteStorageObjects(
	bearerToken string,
	body *api.DeleteStorageObjectsRequest,
	options map[string]string,
) error {
	// Define the URL path
	urlPath := "/v2/storage/delete"

	// Add query parameters (empty for this request)
	queryParams := url.Values{}

	// Convert the body to JSON
	bodyJson, err := json.Marshal(body)
	if err != nil {
		return errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("PUT", fullUrl, strings.NewReader(string(bodyJson)))
	if err != nil {
		return errors.As(err)
	}

	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

func (napi *NakamaApi) ListStorageObjects(
	bearerToken string,
	collection string,
	userId string,
	limit int,
	cursor string,
	options map[string]string,
) (*api.StorageObjectList, error) {
	// Validate the 'collection' parameter
	if !checkStr(&collection) {
		return nil, errors.New("'collection' is a required parameter but is empty.")
	}

	// Define the URL path and replace the placeholder
	urlPath := strings.Replace("/v2/storage/{collection}", "{collection}", url.QueryEscape(collection), 1)

	// Add query parameters
	queryParams := url.Values{}
	if userId != "" {
		queryParams.Set("user_id", userId)
	}
	if limit > 0 {
		queryParams.Set("limit", fmt.Sprintf("%d", limit))
	}
	if cursor != "" {
		queryParams.Set("cursor", cursor)
	}

	// No request body for this function
	bodyJson := ""

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return nil, errors.As(err)
	}

	result := &api.StorageObjectList{}
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return nil, errors.As(err)
	}

	return result, nil
}

func (napi *NakamaApi) ListStorageObjects2(
	bearerToken string,
	collection string,
	userId string,
	limit int,
	cursor string,
	options map[string]string,
) (*api.StorageObjectList, error) {

	// Validate 'collection' and 'userId' parameters
	if !checkStr(&collection) {
		return nil, errors.New("'collection' is a required parameter but is empty.")
	}
	if checkStr(&userId) {
		return nil, errors.New("'userId' is a required parameter but is empty.")
	}

	// Define the URL path and replace placeholders
	urlPath := strings.Replace(
		strings.Replace("/v2/storage/{collection}/{userId}", "{collection}", url.QueryEscape(collection), 1),
		"{userId}", url.QueryEscape(userId), 1,
	)

	// Add query parameters
	queryParams := url.Values{}
	if limit > 0 {
		queryParams.Set("limit", fmt.Sprintf("%d", limit))
	}
	if cursor != "" {
		queryParams.Set("cursor", cursor)
	}

	// No request body for this function
	bodyJson := ""

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return nil, errors.As(err)
	}

	result := &api.StorageObjectList{}
	if err := napi.doReq(bearerToken, req, options, result); err != nil {
		return nil, errors.As(err)
	}

	return result, nil
}

// ListTournaments lists current or upcoming tournaments.
func (napi *NakamaApi) ListTournaments(
	bearerToken string,
	categoryStart *int,
	categoryEnd *int,
	startTime *int64,
	endTime *int64,
	limit int,
	cursor string,
	options map[string]string,
) (*api.TournamentList, error) {
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
	if limit > 0 {
		queryParams.Set("limit", fmt.Sprintf("%d", limit))
	}
	if cursor != "" {
		queryParams.Set("cursor", cursor)
	}

	// No request body for this function
	bodyJson := ""

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	var result api.TournamentList
	if err := napi.doReq(bearerToken, req, options, &result); err != nil {
		return nil, errors.As(err)
	}

	return &result, nil
}

// DeleteTournamentRecord deletes a tournament record.
func (napi *NakamaApi) DeleteTournamentRecord(
	bearerToken string,
	tournamentId string,
	options map[string]string,
) error {
	// Validate the tournamentId
	if tournamentId == "" {
		return errors.New("'tournamentId' is a required parameter but is empty.")
	}

	// Define the URL path
	urlPath := "/v2/tournament/" + url.QueryEscape(tournamentId)

	// No query parameters for this function
	queryParams := url.Values{}

	// No request body for DELETE request
	bodyJson := ""

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("DELETE", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

// ListTournamentRecords lists tournament records.
func (napi *NakamaApi) ListTournamentRecords(
	bearerToken string,
	tournamentId string,
	ownerIds []string,
	limit int,
	cursor string,
	expiry string,
	options map[string]string,
) (*api.TournamentRecordList, error) {

	// Validate the tournamentId
	if !checkStr(&tournamentId) {
		return nil, errors.New("'tournamentId' is a required parameter but is empty.")
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
	if limit > 0 {
		queryParams.Set("limit", fmt.Sprintf("%d", limit))
	}
	if cursor != "" {
		queryParams.Set("cursor", cursor)
	}
	if expiry != "" {
		queryParams.Set("expiry", expiry)
	}

	// No request body for this function
	bodyJson := ""

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return nil, errors.As(err)
	}

	result := &api.TournamentRecordList{}
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return nil, errors.As(err)
	}

	return result, nil
}

// WriteTournamentRecord2 writes a record to a tournament.
func (napi *NakamaApi) WriteTournamentRecord2(
	bearerToken string,
	tournamentId string,
	record *api.WriteTournamentRecordRequest,
	options map[string]string,
) (*api.LeaderboardRecord, error) {

	// Validate the tournamentId and record
	if checkStr(&tournamentId) {
		return nil, errors.New("'tournamentId' is a required parameter but is empty.")
	}
	if record == nil {
		return nil, errors.New("'record' is a required parameter but is empty.")
	}

	// Define the URL path
	urlPath := "/v2/tournament/" + url.QueryEscape(tournamentId)

	// Prepare the request body
	bodyJson, err := json.Marshal(record)
	if err != nil {
		return nil, errors.As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, nil)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, strings.NewReader(string(bodyJson)))
	if err != nil {
		return nil, errors.As(err)
	}

	result := &api.LeaderboardRecord{}
	if err := napi.doReq(bearerToken, req, options, result); err != nil {
		return nil, errors.As(err)
	}

	return result, nil
}

// WriteTournamentRecord writes a record to a tournament.
func (napi *NakamaApi) WriteTournamentRecord(
	bearerToken string,
	tournamentId string,
	record *api.WriteTournamentRecordRequest_TournamentRecordWrite,
	options map[string]string,
) (*api.LeaderboardRecord, error) {

	// Validate the tournamentId and record
	if !checkStr(&tournamentId) {
		return nil, errors.New("'tournamentId' is a required parameter but is empty.")
	}
	if record == nil {
		return nil, errors.New("'record' is a required parameter but is empty.")
	}

	// Define the URL path
	urlPath := "/v2/tournament/" + url.QueryEscape(tournamentId)

	// Prepare the request body
	bodyJson, err := json.Marshal(record)
	if err != nil {
		return nil, errors.New("failed to marshal record").As(err)
	}

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, nil)

	// Prepare the HTTP request
	req, err := http.NewRequest("PUT", fullUrl, strings.NewReader(string(bodyJson)))
	if err != nil {
		return nil, errors.As(err)
	}

	result := &api.LeaderboardRecord{}
	if err := napi.doReq(bearerToken, req, options, result); err != nil {
		return nil, errors.As(err)
	}

	return result, nil
}

// JoinTournament attempts to join an open and running tournament.
func (napi *NakamaApi) JoinTournament(
	bearerToken string,
	tournamentId string,
	options map[string]string,
) error {

	// Validate the tournamentId
	if !checkStr(&tournamentId) {
		return errors.New("'tournamentId' is a required parameter but is empty.")
	}

	// Define the URL path
	urlPath := "/v2/tournament/" + url.QueryEscape(tournamentId) + "/join"

	// Prepare the query params (if any, currently empty map)
	queryParams := url.Values{}

	// Prepare the request body
	bodyJson := ""

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("POST", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return errors.As(err)
	}
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return errors.As(err)
	}
	return nil
}

// ListTournamentRecordsAroundOwner lists tournament records for a given owner.
func (napi *NakamaApi) ListTournamentRecordsAroundOwner(
	bearerToken string,
	tournamentId string,
	ownerId string,
	limit int,
	expiry string,
	cursor string,
	options map[string]string,
) (*api.TournamentRecordList, error) {

	// Validate the tournamentId and ownerId
	if !checkStr(&tournamentId) {
		return nil, errors.New("'tournamentId' is a required parameter but is empty.")
	}
	if !checkStr(&ownerId) {
		return nil, errors.New("'ownerId' is a required parameter but is empty.")
	}

	// Define the URL path
	urlPath := "/v2/tournament/" + url.QueryEscape(tournamentId) + "/owner/" + url.QueryEscape(ownerId)

	// Prepare the query params
	queryParams := url.Values{}
	if limit > 0 {
		queryParams.Set("limit", fmt.Sprintf("%d", limit))
	}
	if expiry != "" {
		queryParams.Set("expiry", expiry)
	}
	if cursor != "" {
		queryParams.Set("cursor", cursor)
	}

	// Prepare the request body (empty)
	bodyJson := ""

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return nil, errors.As(err)
	}

	var result api.TournamentRecordList
	if err := napi.doReq(bearerToken, req, options, nil); err != nil {
		return nil, errors.As(err)
	}
	return &result, nil
}

// GetUsers fetches zero or more users by ID and/or username.
func (napi *NakamaApi) GetUsers(
	bearerToken *string,
	ids []string,
	usernames []string,
	facebookIds []string,
	options map[string]string,
) (*api.Users, error) {

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
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return nil, err
	}
	var result api.Users
	if err := napi.doReq("", req, options, &result); err != nil {
		return nil, errors.As(err)
	}

	return &result, nil

}

// ListUserGroups lists the groups the current user belongs to.
func (napi *NakamaApi) ListUserGroups(
	bearerToken string,
	userId string,
	state *int,
	limit int,
	cursor string,
	options map[string]string,
) (*api.UserGroupList, error) {

	// Validate required parameters
	if !checkStr(&userId) {
		return nil, errors.New("'userId' is a required parameter but is empty.")
	}

	// Define the URL path and replace placeholder
	urlPath := "/v2/user/{userId}/group"
	urlPath = strings.Replace(urlPath, "{userId}", url.QueryEscape(userId), 1)

	// Prepare the query params
	queryParams := url.Values{}
	if limit > 0 {
		queryParams.Set("limit", strconv.Itoa(limit))
	}
	if state != nil {
		queryParams.Set("state", strconv.Itoa(*state))
	}
	if cursor != "" {
		queryParams.Set("cursor", cursor)
	}

	// Prepare the request body (empty)
	bodyJson := ""

	// Construct the full URL
	fullUrl := napi.buildFullUrl(napi.BasePath, urlPath, queryParams)

	// Prepare the HTTP request
	req, err := http.NewRequest("GET", fullUrl, strings.NewReader(bodyJson))
	if err != nil {
		return nil, errors.As(err)
	}

	var result api.UserGroupList
	if err := napi.doReq(bearerToken, req, options, &result); err != nil {
		return nil, errors.As(err)
	}

	return &result, nil
}

func (napi *NakamaApi) buildFullUrl(basePath string, fragment string, queryParams url.Values) string {
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
