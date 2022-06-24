// Copyright 2018 Nervos Foundation

package eaglesong

func permute(rounds int, state []uint32) {
	for i := 0; i < rounds; i++ {
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

func eaglesong(rounds, capacity, rate int, delim byte, input []byte) []byte {
	// absorbing
	state := make([]uint32, 16)
	for i := 0; i < ((len(input)+1)*8+rate-1)/rate; i++ {
		for j := 0; j < rate/32; j++ {
			var integer uint32
			for k := 0; k < 4; k++ {
				if i*rate/8+j*4+k < len(input) {
					integer = (integer << 8) ^ uint32(input[i*rate/8+j*4+k])
				} else if i*rate/8+j*4+k == len(input) {
					integer = (integer << 8) ^ uint32(delim)
				}
			}

			state[j] = state[j] ^ integer
		}

		permute(rounds, state)
	}

	// squeezing
	output := make([]byte, capacity)
	for i := 0; i < capacity/(rate/8); i++ {
		for j := 0; j < rate/32; j++ {
			for k := 0; k < 4; k++ {
				output[i*rate/8+j*4+k] = byte((state[j] >> (8 * k)) & 0xff)
			}
		}

		permute(rounds, state)
	}

	return output
}
