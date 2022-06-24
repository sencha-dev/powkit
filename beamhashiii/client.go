package beamhashiii

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
	return verify(c.n, c.k, c.personal, header, soln)
}
