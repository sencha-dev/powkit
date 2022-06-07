// Copyright (c) 2021 Electric Coin Company

package equihash

import (
	"bytes"
	"testing"

	"github.com/sencha-dev/powkit/internal/common/testutil"
)

func TestTraditionalVerify(t *testing.T) {
	tests := []struct {
		n        uint32
		k        uint32
		personal string
		seed     []byte
		input    []byte
		nonce    uint32
		valid    bool
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
		valid := New(tt.n, tt.k, tt.personal).TraditionalVerify(tt.seed, tt.input, tt.nonce)
		if valid != tt.valid {
			t.Errorf("failed on simple test %d: have %t want %t", i, valid, tt.valid)
		}
	}
}

func TestZCashVerify(t *testing.T) {
	tests := []struct {
		n        uint32
		k        uint32
		personal string
		header   []byte
		soln     []byte
	}{
		{
			n:        96,
			k:        5,
			personal: "ZcashPoW",
			header: bytes.Join([][]byte{
				[]byte("Equihash is an asymmetric PoW based on the Generalised Birthday problem."),
				[]byte{
					1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
					0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
				},
			}, nil),
			soln: []byte{
				0x04, 0x6a, 0x8e, 0xd4, 0x51, 0xa2, 0x19, 0x73,
				0x32, 0xe7, 0x1f, 0x39, 0xdb, 0x9c, 0x79, 0xfb,
				0xf9, 0x3f, 0xc1, 0x44, 0x3d, 0xa5, 0x8f, 0xb3,
				0x8d, 0x05, 0x99, 0x17, 0x21, 0x16, 0xd5, 0x55,
				0xb1, 0xb2, 0x1f, 0x32, 0x70, 0x5c, 0xe9, 0x98,
				0xf6, 0x0d, 0xa8, 0x52, 0xf7, 0x7f, 0x0e, 0x7f,
				0x4d, 0x63, 0xfc, 0x2d, 0xd2, 0x30, 0xa3, 0xd9,
				0x99, 0x53, 0xa0, 0x78, 0x7d, 0xfe, 0xfc, 0xab,
				0x34, 0x1b, 0xde, 0xc8,
			},
		},
	}

	for i, tt := range tests {
		valid, err := New(tt.n, tt.k, tt.personal).ZCashVerify(tt.header, tt.soln)
		if err != nil {
			t.Errorf("failed on test %d: %v", i, err)
			continue
		}

		if !valid {
			t.Errorf("failed on test %d: have false want true", i)
			continue
		}
	}
}
