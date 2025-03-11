package pihttp

type Request struct {
	requestLine
	Headers    Header
	RawQuery   string
	EntityBody []byte
}
