package autolykos2

import (
	"fmt"
)

type Client struct {
	k     uint32
	n     uint32
	nBase uint32
}

func New(k, n uint32) *Client {
	c := &Client{
		k:     k,
		n:     n,
		nBase: 1 << n,
	}

	return c
}

func NewErgo() *Client {
	return New(32, 26)
}

func (c *Client) Compute(msg []byte, height, nonce uint64) ([]byte, error) {
	if len(msg) != 32 {
		return nil, fmt.Errorf("msg must be 32 bytes")
	}

	return compute(c.k, c.nBase, msg, nonce, height), nil
}
