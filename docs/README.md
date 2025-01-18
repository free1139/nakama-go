# Nakama Go Client Guide

This client library guide will show you how to use the core Nakama features in Go by showing you how to develop the
Nakama specific parts (without full game logic or UI) of
an [Among Us (external)](https://www.innersloth.com/games/among-us/) inspired game called Sagi-shi (Japanese
for "Imposter").

## Prerequisites

Before proceeding ensure that you have:

- Installed Nakama server
- Installed the [Nakama Go SDK](https://pkg.go.dev/github.com/NorthNorthGames/nakama-go)

### Full API documentation

For the full API documentation please visit the [API docs](#nakama-client).

### Installation

Add the dependency to your `go.mod` file:

```shell
go get go get https://pkg.go.dev/github.com/NorthNorthGames/nakama-go
```

After installing the client import it into your project:

```go
import(
"github.com/NorthNorthGames/nakama-go"
)
```

In your main Go function create a [client object](#nakama-client)

### Updates

New versions of the Nakama Go Client and the corresponding improvements are documented in the [Release Notes]().

## Getting started

Learn how to get started using the Nakama Client and Socket objects to start building Sagi-shi and your own game.

### Nakama Client

The Nakama Client connects to a Nakama Server and is the entry point to access Nakama features. It is recommended to
have one client per server per game.

To create a client for Sagi-shi pass in your server connection details:

```go
client := NewClient("defaultKey", "127.0.0.1", "7350", false, nil, nil)
```

### Configuring the Request Timeout Length

Each request to Nakama from the client must complete in a certain period of time before it is considered to have timed
out. You can configure how long this period is (in milliseconds) by setting the timeout value on the client:

```go
client.timeout = 10000
```

### Nakama Socket

The Nakama Socket is used for gameplay and real-time latency-sensitive features such as chat, parties, matches and RPCs.

From the client create a socket:

```go
socket := client.CreateSocket(false, false, nil, nil)

appearOnline := true
socket.Connect(session, appearOnline)
```

## Authentication

Nakama has many [authentication methods](https://heroiclabs.com/docs/nakama/concepts/authentication/) and supports
creating [custom authentication](https://heroiclabs.com/docs/nakama/concepts/authentication/#custom) on the server.

Sagi-shi will use device and Facebook authentication, linked to the same user account so that players can play from
multiple devices.

### Device authentication

Nakama [Device Authentication](https://heroiclabs.com/docs/nakama/concepts/authentication/#device) uses the physical
device’s unique identifier to easily authenticate a user and create an account if one does not exist.

When using only device authentication, you don’t need a login UI as the player can automatically authenticate when the
game launches.

Authentication is an example of a Nakama feature accessed from a Nakama Client instance.

```go
// Simulated device ID retrieval and storage
var deviceID string

// Check if device ID exists in a file (simulating AsyncStorage equivalent)
file, err := os.Open(storageKey)
if err == nil {
// Read the device ID from the file
fmt.Fscanf(file, "%s", &deviceID)
file.Close()
} else {
// Generate a new unique ID for the device
deviceID = fmt.Sprintf("device-%d", memory.TotalMemory())

// Save the device ID in a file for future reference
file, err = os.Create(storageKey)
if err != nil {
log.Fatalf("Failed to save device ID: %v", err)
}
_, writeErr := file.WriteString(deviceID)
if writeErr != nil {
log.Fatalf("Failed to write device ID: %v", writeErr)
}
file.Close()
}

// Authenticate with the Nakama server using Device Authentication
session, authErr := client.AuthenticateDevice(deviceID, create, username)
if authErr != nil {
log.Fatalf("Failed to authenticate: %v", authErr)
}

log.Printf("Successfully authenticated: %+v\n", session)
```

### Facebook authentication

Nakama [Facebook Authentication](https://heroiclabs.com/docs/nakama/concepts/authentication/#facebook) is an easy-to-use
authentication method which lets you optionally import the player’s Facebook friends and add them to their Nakama
Friends list.

```go
oauthToken := "<token>"
importFriends := true

session, err := client.AuthenticateFacebook(oauthToken, true, "mycustomusername", importFriends)
if err != nil {
log.Fatalf("Error authenticating with Facebook: %v", err)
}
log.Printf("Successfully authenticated: %+v\n", session)
```

### Custom authentication

Nakama supports [Custom Authentication](https://heroiclabs.com/docs/nakama/concepts/authentication/#custom) methods to
integrate with additional identity services.

See
the [Itch.io custom authentication](https://heroiclabs.com/docs/nakama/client-libraries/snippets/custom-authentication/)
recipe for an example.
