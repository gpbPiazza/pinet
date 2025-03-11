package pihttp

type Response struct {
	Header     Header
	StatusCode int
	Body       any
	// statusLine     string
	// generalHeader  string
	// responseHeader string
	// entityHeader   string
	// entityBody     string
}
