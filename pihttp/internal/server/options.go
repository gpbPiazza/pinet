package server

type options struct {
	handler Handler
}

type Option interface {
	apply(*options)
}

func WithHandler(handler Handler) Option {
	return &optionWithHandler{
		handler: handler,
	}
}

type optionWithHandler struct {
	handler Handler
}

func (o *optionWithHandler) apply(opts *options) {
	opts.handler = o.handler
}
