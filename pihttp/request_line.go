package pihttp

import "strings"

type requestLine struct {
	Method      string
	Uri         string
	HttpVersion string
}

func (rl requestLine) RawQuery() string {
	uriSplited := strings.Split(rl.Uri, "?")
	if len(uriSplited) < 1 {
		return ""
	}

	return uriSplited[1]
}
