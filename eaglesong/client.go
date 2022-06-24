package eaglesong

type Client struct {
	rounds   int
	capacity int
	rate     int
	length   int
	delim    byte
}

func New(rounds, capacity, rate, length int, delim byte) *Client {
	cfg := &Client{
		rounds:   rounds,
		capacity: capacity,
		rate:     rate,
		length:   length,
		delim:    delim,
	}

	return cfg
}

func NewNervos() *Client {
	return New(43, 32, 256, 32, 0x06)
}

func (c *Client) Compute(input []byte) []byte {
	return eaglesong(c.rounds, c.capacity, c.rate, c.delim, input)
}
