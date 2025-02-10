package pihttp

// -> https://pkg.go.dev/net#pkg-examples
// -> https://itsfuad.medium.com/understanding-http-at-a-low-level-a-developers-guide-with-c-213728d6c41d

// Create a struct that is capeable to stablish a connection using TCP/IP protocol from net package.
// This struct shoul be able to handle clients to connect in they open socket IP connection and response the request client

// Start small, just create the server, connect into a PORT and try connect in this port as client e return OK every time a
// client request comes.

type HTTP struct {
}

func NewServer() *HTTP {
	return &HTTP{}
}
