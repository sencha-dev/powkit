package heavyhash

import (
	"fmt"
)

type Client struct{}

func New() *Client {
	client := &Client{}

	return client
}

func NewKaspa() *Client {
	return New()
}

func (c *Client) Compute(hash []byte, timestamp int64, nonce uint64) ([]byte, error) {
	if len(hash) != 32 {
		return nil, fmt.Errorf("hash must be 32 bytes")
	}

	digest := heavyHash(hash, timestamp, nonce)

	return digest, nil
}
