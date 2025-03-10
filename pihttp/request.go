package pihttp

type Request struct {
	requestLine
	Headers    header
	RawQuery   string
	EntityBody []byte
}
