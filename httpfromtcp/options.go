package httpfromtcp

type options struct {
}

type Option interface {
	apply(*options)
}
