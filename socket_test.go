package nakama

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupSocket(t *testing.T) (*Client, *Session) {
	client := NewClient("defaultkey", "127.0.0.1", "7350", false, nil, nil)

	deviceId := "376C007D-260F-579B-BD75-A3CBBFC2EF99"
	create := true
	session, err := client.AuthenticateDevice(deviceId, &create, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	return client, session
}

func TestCreateMatch_NoName(t *testing.T) {
	client, session := setupSocket(t)

	socket := createSocket(t, client, session)

	//matchName := "Test"
	match, err := socket.CreateMatch(nil)

	assert.NoError(t, err)
	assert.NotNil(t, match)
	assert.IsType(t, &Match{}, match)
}

func TestCreateMatch_WithName(t *testing.T) {
	client, session := setupSocket(t)

	socket := createSocket(t, client, session)

	matchName := "Test"
	match, err := socket.CreateMatch(&matchName)

	assert.NoError(t, err)
	assert.NotNil(t, match)
	assert.IsType(t, &Match{}, match)
}

func createSocket(t *testing.T, client *Client, session *Session) *DefaultSocket {
	timeout := 1000
	socket := client.CreateSocket(false, true, nil, &timeout)

	assert.IsType(t, DefaultSocket{}, socket)

	err := socket.Connect(session, nil, &timeout, nil)
	assert.NoError(t, err)

	return socket
}
