package pihttp

type options struct {
	httpV                           string
	readBufferSize, writeBufferSize int
}

type Option interface {
	apply(*options)
}
