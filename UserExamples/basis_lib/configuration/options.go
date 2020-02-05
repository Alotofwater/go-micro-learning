package configuration

import "github.com/micro/go-micro/config/source"

type Options struct {
	Sources []source.Source
	PathPrefix string
}

type Option func(o *Options)


func WithSource(src source.Source) Option {
	return func(o *Options) {
		o.Sources = append(o.Sources, src)
	}
}


func WithPathPrefix(val string) Option {
	return func(o *Options) {
		o.PathPrefix = val
	}
}