package equihash

type Config struct {
	n        uint32
	k        uint32
	personal []byte
}

func New(n, k uint32, personal string) *Config {
	cfg := &Config{
		n:        n,
		k:        k,
		personal: []byte(personal),
	}

	return cfg
}

func (cfg *Config) TraditionalVerify(seed, input []byte, nonce uint32) bool {
	return TraditionalVerify(cfg.n, cfg.k, cfg.personal, seed, input, nonce)
}

func (cfg *Config) ZCashVerify(header, soln []byte) (bool, error) {
	return ZCashVerify(cfg.n, cfg.k, cfg.personal, header, soln)
}
