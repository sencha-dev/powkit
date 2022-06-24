package equihash

type Client struct {
	n        uint32
	k        uint32
	personal []byte
	twist    bool
}

func New(n, k uint32, personal string, twist bool) *Client {
	cfg := &Client{
		n:        n,
		k:        k,
		personal: []byte(personal),
		twist:    twist,
	}

	return cfg
}

func NewFlux() *Client {
	return New(125, 4, "ZelProof", true)
}

func NewBitcoinGold() *Client {
	return New(144, 5, "BgoldPoW", false)
}

func NewZClassic() *Client {
	return New(192, 7, "ZcashPoW", false)
}

func NewZCash() *Client {
	return New(200, 9, "ZcashPoW", false)
}

func NewAion() *Client {
	return New(210, 9, "AION0PoW", false)
}

func (c *Client) Verify(header, soln []byte) (bool, error) {
	return verify(c.n, c.k, c.personal, header, soln, c.twist)
}
