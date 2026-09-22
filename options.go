package jev

import (
	"net/http"
)

type Option func(j *Jev)

func WithClient(c *http.Client) Option {
	return func(j *Jev) {
		if c != nil {
			j.transport.HTTP = c
		}
	}
}

// WithLogger sends the client's debug output to l. The default logger
// discards it. A nil logger is rejected by New rather than ignored.
func WithLogger(l Logger) Option {
	return func(j *Jev) {
		j.transport.Logger = l
	}
}

// WithEndpoint points the client at another URL: a proxy, a gateway or a
// compatible service. It is the full endpoint the request is posted to, path
// included, not just a host -- nothing is appended to it.
func WithEndpoint(endpoint string) Option {
	return func(j *Jev) {
		j.transport.Endpoint = endpoint
	}
}

func WithModel(model string) Option {
	return func(j *Jev) {
		j.model = model
	}
}
