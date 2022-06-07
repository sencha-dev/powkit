// Copyright 2018 Nervos Foundation

package eaglesong

func (cfg *Config) permutation(state []uint32) {
	for i := 0; i < cfg.rounds; i++ {
		temp := make([]uint32, 16)
		for j := 0; j < 16; j++ {
			for k := 0; k < 16; k++ {
				temp[j] = temp[j] ^ (bitMatrix[k*16+j] * state[k])
			}
		}

		for j := 0; j < 16; j++ {
			state[j] = temp[j]
		}

		// circulant multiplication
		for j := 0; j < 16; j++ {
			value := state[j]
			state[j] ^= (value << coefficients[3*j+1])
			state[j] ^= (value >> (32 - coefficients[3*j+1]))
			state[j] ^= (value << coefficients[3*j+2])
			state[j] ^= (value >> (32 - coefficients[3*j+2]))
		}

		// constants injection
		for j := 0; j < 16; j++ {
			state[j] = state[j] ^ injectionConstants[i*16+j]
		}

		// addition / rotation / addition
		for j := 0; j < 16; j += 2 {
			state[j] = state[j] + state[j+1]
			state[j] = (state[j] << 8) ^ (state[j] >> 24)
			state[j+1] = (state[j+1] << 24) ^ (state[j+1] >> 8)
			state[j+1] = state[j] + state[j+1]
		}
	}
}

func (cfg *Config) eaglesong(input []byte) []byte {
	// absorbing
	state := make([]uint32, 16)
	for i := 0; i < ((len(input)+1)*8+cfg.rate-1)/cfg.rate; i++ {
		for j := 0; j < cfg.rate/32; j++ {
			var integer uint32
			for k := 0; k < 4; k++ {
				if i*cfg.rate/8+j*4+k < len(input) {
					integer = (integer << 8) ^ uint32(input[i*cfg.rate/8+j*4+k])
				} else if i*cfg.rate/8+j*4+k == len(input) {
					integer = (integer << 8) ^ uint32(cfg.delim)
				}
			}

			state[j] = state[j] ^ integer
		}
		cfg.permutation(state)
	}

	// squeezing
	output := make([]byte, cfg.capacity)
	for i := 0; i < cfg.capacity/(cfg.rate/8); i++ {
		for j := 0; j < cfg.rate/32; j++ {
			for k := 0; k < 4; k++ {
				output[i*cfg.rate/8+j*4+k] = byte((state[j] >> (8 * k)) & 0xff)
			}
		}
		cfg.permutation(state)
	}

	return output
}
