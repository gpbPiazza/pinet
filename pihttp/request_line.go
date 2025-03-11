package pihttp

import "strings"

const (
	queryDelimiter = "?"
)

var (
	queryDelimiterBytes = []byte(queryDelimiter)
)

type requestLine struct {
	Method      string
	Uri         string
	Path        string
	HttpVersion string
}

func (rl requestLine) RawQuery() string {
	uriSplited := strings.Split(rl.Uri, queryDelimiter)
	if len(uriSplited) < 1 {
		return ""
	}

	return uriSplited[1]
}
