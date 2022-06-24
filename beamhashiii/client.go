package beamhashiii

import (
	"fmt"
)

type Client struct {
	n        uint32
	k        uint32
	personal []byte
}

func New(n, k uint32, personal string) *Client {
	client := &Client{
		n:        n,
		k:        k,
		personal: []byte(personal),
	}

	return client
}

func NewBeam() *Client {
	return New(150, 5, "Beam-PoW")
}

func (c *Client) Verify(header, soln []byte) (bool, error) {
	if len(header) != 40 {
		return false, fmt.Errorf("header must be 40 bytes")
	} else if len(soln) != 104 {
		return false, fmt.Errorf("soln must be 104 bytes")
	}

	return verify(c.n, c.k, c.personal, header, soln)
}
