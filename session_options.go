package lrmr

import "time"

type SessionOptions struct {
	Name         string
	Timeout      time.Duration
	NodeSelector map[string]string
}

type SessionOption func(o *SessionOptions)

func WithName(n string) SessionOption {
	return func(o *SessionOptions) {
		o.Name = n
	}
}

func WithTimeout(d time.Duration) SessionOption {
	return func(o *SessionOptions) {
		o.Timeout = d
	}
}

func WithNodeSelector(selector map[string]string) SessionOption {
	return func(o *SessionOptions) {
		o.NodeSelector = selector
	}
}

func buildSessionOptions(opts []SessionOption) (o SessionOptions) {
	for _, optFn := range opts {
		optFn(&o)
	}
	return o
}
