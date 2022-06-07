package eaglesong

type Config struct {
	rounds   int
	capacity int
	rate     int
	length   int
	delim    byte
}

func New(rounds, capacity, rate, length int, delim byte) *Config {
	cfg := &Config{
		rounds:   rounds,
		capacity: capacity,
		rate:     rate,
		length:   length,
		delim:    delim,
	}

	return cfg
}

func NewNervos() *Config {
	var cfg = &Config{
		rounds:   43,
		capacity: 32,
		rate:     256,
		length:   32,
		delim:    0x06,
	}

	return cfg
}

func (cfg *Config) Compute(input []byte) []byte {
	return cfg.eaglesong(input)
}
