package server

type options struct {
}

type Option interface {
	apply(*options)
}
