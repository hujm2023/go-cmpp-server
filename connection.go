package cmppserver

import (
	"github.com/cloudwego/netpoll"
)

type Connection struct {
	conn netpoll.Reader
}

func (c *Connection) Peek(n int) ([]byte, error) {
	return c.conn.Peek(n)
}

func (c *Connection) Discard(n int) (int, error) {
	return n, c.conn.Skip(n)
}

func (c *Connection) Size() int {
	return c.conn.Len()
}

// Read .
// Do not use.
func (c *Connection) Read(p []byte) (n int, err error) {
	// make a copy
	r, err := c.conn.Slice(n)
	if err != nil {
		return 0, err
	}
	data, err := r.Next(-1)
	if err != nil {
		return 0, err
	}
	copy(p, data)
	return len(data), nil
}
