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

The Nakama Client connects to a Nakama Server and is the entry point to access Nakama features. It is recommended to have one client per server per game.

To create a client for Sagi-shi pass in your server connection details:

```go
client := NewClient("defaultKey", "127.0.0.1", "7350", false, nil, nil)
```

### Configuring the Request Timeout Length

Each request to Nakama from the client must complete in a certain period of time before it is considered to have timed out. You can configure how long this period is (in milliseconds) by setting the timeout value on the client:

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