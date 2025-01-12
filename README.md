# Nakama Go client

[Nakama](https://github.com/heroiclabs/nakama) is an open-source server designed to power modern games and apps.
Features include user accounts, chat, social, matchmaker, realtime multiplayer, and much [more](https://heroiclabs.com).

This client implements the full API and socket options with the server. It's written in Go with minimal dependencies.

## Getting Started

You'll need to set up the server and database before you can connect with the client. The simplest way is to use Docker
but have a look at the [server documentation](https://github.com/heroiclabs/nakama#getting-started) for other options.

1. Install and run the servers. Follow
   these [instructions](https://heroiclabs.com/docs/nakama/getting-started/install/docker/).

2. Import the client into your project.
   It's [available on pkg.go.dev](https://pkg.go.dev/github.com/NorthNorthGames/nakama-go)

   ```shell
   go get https://pkg.go.dev/github.com/NorthNorthGames/nakama-go
   ```

3. Use the connection credentials to build a client object.

   ```go
   package main

   import "github.com/NorthNorthGames/nakama-go"
   
   func main() {
	   useSsl := false // Enable if server is run with an SSL certificate
	   client := NewClient("defaultKey", "127.0.0.1", "7350", useSsl, nil, nil)
   }
   ```

## Usage

The client object has many methods to execute various features in the server or open realtime socket connections with
the server.

### Authenticate

```go
email := "super@heroes.com"
password := "batsignal"
session, error := client.AuthenticateEmail(email, password, nil, nil, nil)

if error != nil {
   log.Fatalf("Failed to authenticate email: %v", error)
}
log.Printf("Authenticated successfully. Session Token: %v", *session)
```

### Sessions

When authenticated the server responds with an auth token (JWT) which contains useful properties and gets deserialized
into a `Session` struct.

```go
log.Print(session.Token)
log.Print(session.RefreshToken)
log.Print(session.UserID)
log.Print(session.Username)

log.Printf("Session has expired? %t", session.IsExpired(time.Now().UnixMilli()/1000))

expiresAt := session.ExpiresAt
var expiresAtTime time.Time

if expiresAt != nil {
   expiresAtTime = time.UnixMilli(*expiresAt)
   log.Printf("Session will expire at: %s", expiresAtTime.Format(time.RFC3339))
} else {
    log.Print("Expiration time is not set.")
}
```

It is recommended to store the auth token from the session and check at startup if it has expired. If the token has
expired you must reauthenticate. The expiry time of the token can be changed as a setting in the server.

```go
session := Restore(authToken, refreshToken)

// Check whether a session is close to expiry

unixTimeInFuture := time.Now().Add(24 * time.Hour).UnixMilli()

if session.IsExpired(unixTimeInFuture / 1000) {
    session, error = client.SessionRefresh(session, make(map[string]string))

   if error != nil {
      log.Print("Session can no longer be refreshed. Must reauthenticate!")
   }
}
```

### Requests

The client includes lots of builtin APIs for various features of the game server. These can be accessed with the methods
which return Promise objects. It can also call custom logic as RPC functions on the server. These can also be executed
with a socket object.

All requests are sent with a session object which authorizes the client.

```go
account, error := client.GetAccount(session)

if account == nil {
    log.Fatalf("Failed to get account: %v", error)
}
if error != nil {
    log.Fatalf("Failed to get account: %v", error)
}

log.Print(account.User.ID)
log.Print(account.User.Username)
log.Print(account.Wallet)
```

### Socket

The client can create one or more sockets with the server. Each socket can have its own event listeners registered for
responses received from the server.

```go
secure := false
trace := false
socket := client.CreateSocket(secure, trace, nil, nil)

session, _ := socket.Connect(*session, nil, nil)
// Socket it open
```

There's many messages for chat, realtime, status events, notifications, etc. which can be sent or received from the socket.

```go
// Join a chat channel
roomName := "mychannel"
chatType := 1 // 1 = Room, 2 = Direct Message, 3 = Group
persistence := false
hidden := false

channel, err := socket.JoinChat(roomName, chatType, persistence, hidden)
if err != nil {
    log.Fatalf("Failed to join chat: %v", err)
}

// Send a message to the channel
message := map[string]interface{}{
    "hello": "world",
}

_, err = socket.WriteChatMessage(channel.ID, message)
if err != nil {
    log.Fatalf("Failed to send chat message: %v", err)
}
```

## Contribute

The development roadmap is managed as GitHub issues and pull requests are welcome. If you're interested in enhancing the code please open an issue to discuss the changes.

### License

This project is licensed under the [MIT License](https://github.com/NorthNorthGames/nakama-go/blob/main/LICENSE).
