package stargate

import "context"

type Context struct {
	context.Context

	headers map[string]string
}

func (c *Context) AddHeader(header, value string) {
	if c.headers == nil {
		c.headers = map[string]string{}
	}

	c.headers[header] = value
}

func (c *Context) GetHeaders() map[string]string {
	return c.headers
}
