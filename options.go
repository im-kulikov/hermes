package hermes

import "time"

var defaultSocketURL = "wss://ws.pusherapp.com/app/de504dc5763aeef9ff52?protocol=7&client=js&version=2.1.6&flash=false"

type Option = func(*Options)

type Options struct {
	url      string
	deadline time.Duration
}

func URL(url string) Option {
	return func(o *Options) {
		o.url = url
	}
}

func Deadline(t time.Duration) Option {
	return func(o *Options) {
		o.deadline = t
	}
}

func newOptions(opts ...Option) Options {
	var opt = Options{
		url: defaultSocketURL,
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}
