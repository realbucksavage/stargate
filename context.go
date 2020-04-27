package stargate

import "context"

type Context struct {
	context.Context

	headers map[string]string
}

func (c *Context) AddHeader(header, value string) {
	c.headers[header] = value
}

func NewContext() *Context {
	return WithContext(context.Background())
}

func WithContext(ctx context.Context) *Context {
	return &Context{
		Context: ctx,
		headers: map[string]string{},
	}
}
