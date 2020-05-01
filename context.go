package stargate

type Context struct {
	headers map[string]string
}

func (c *Context) AddHeader(header, value string) {
	c.headers[header] = value
}
