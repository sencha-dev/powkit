package equihash

import (
	"testing"

	"github.com/sencha-dev/go-pow/internal/common/testutil"
)

func TestSimple(t *testing.T) {
	tests := []struct {
		n     uint64
		k     uint64
		seed  []byte
		input []byte
		nonce uint32
		valid bool
	}{
		{
			n:     90,
			k:     5,
			seed:  testutil.MustDecodeHex("b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"),
			input: testutil.MustDecodeHex("f903000030070000e0140000a0fd0000a829000018410000890d0000b8120000242b00004d770000154f0000de730000902d000034a40000050e00007e2300000f0700001dcc0000725600000b790000329f000004e600006b1500000d8b0000285d00009b8b0000c01b0000cb970000e4080000721b00007ac40000a0e70000"),
			nonce: 4,
			valid: true,
		},
	}

	for i, tt := range tests {
		equihash, err := New(tt.n, tt.k, "")
		if err != nil {
			t.Errorf("failed on simple test %d: %v", i, err)
		}

		valid := equihash.Verify(tt.seed, tt.input, tt.nonce)
		if valid != tt.valid {
			t.Errorf("failed on simple test %d: have %t want %t", i, valid, tt.valid)
		}
	}
}
