// Copyright (c) 2021 Electric Coin Company

package equihash

import (
	"bytes"
	"reflect"
	"testing"
)

func TestExpandArray(t *testing.T) {
	tests := []struct {
		length   uint32
		pad      uint32
		compact  []byte
		expected []byte
	}{
		{
			length: 11,
			pad:    0,
			compact: []byte{
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff,
			},
			expected: []byte{
				0x07, 0xff, 0x07, 0xff, 0x07, 0xff, 0x07, 0xff,
				0x07, 0xff, 0x07, 0xff, 0x07, 0xff, 0x07, 0xff,
			},
		},
		{
			length: 21,
			pad:    0,
			compact: []byte{
				0xaa, 0xaa, 0xad, 0x55, 0x55, 0x6a, 0xaa, 0xab,
				0x55, 0x55, 0x5a, 0xaa, 0xaa, 0xd5, 0x55, 0x56,
				0xaa, 0xaa, 0xb5, 0x55, 0x55,
			},
			expected: []byte{
				0x15, 0x55, 0x55, 0x15, 0x55, 0x55, 0x15, 0x55,
				0x55, 0x15, 0x55, 0x55, 0x15, 0x55, 0x55, 0x15,
				0x55, 0x55, 0x15, 0x55, 0x55, 0x15, 0x55, 0x55,
			},
		},
		{
			length: 21,
			pad:    0,
			compact: []byte{
				0x00, 0x02, 0x20, 0x00, 0x0a, 0x7f, 0xff, 0xfe,
				0x00, 0x12, 0x30, 0x22, 0xb3, 0x82, 0x26, 0xac,
				0x19, 0xbd, 0xf2, 0x34, 0x56,
			},
			expected: []byte{
				0x00, 0x00, 0x44, 0x00, 0x00, 0x29, 0x1f, 0xff,
				0xff, 0x00, 0x01, 0x23, 0x00, 0x45, 0x67, 0x00,
				0x89, 0xab, 0x00, 0xcd, 0xef, 0x12, 0x34, 0x56,
			},
		},
		{
			length: 14,
			pad:    0,
			compact: []byte{
				0xcc, 0xcf, 0x33, 0x3c, 0xcc, 0xf3, 0x33, 0xcc,
				0xcf, 0x33, 0x3c, 0xcc, 0xf3, 0x33, 0xcc, 0xcf,
				0x33, 0x3c, 0xcc, 0xf3, 0x33, 0xcc, 0xcf, 0x33,
				0x3c, 0xcc, 0xf3, 0x33,
			},
			expected: []byte{
				0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33,
				0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33,
				0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33,
				0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33, 0x33,
			},
		},
		{
			length: 11,
			pad:    2,
			compact: []byte{
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff,
			},
			expected: []byte{
				0x00, 0x00, 0x07, 0xff, 0x00, 0x00, 0x07, 0xff,
				0x00, 0x00, 0x07, 0xff, 0x00, 0x00, 0x07, 0xff,
				0x00, 0x00, 0x07, 0xff, 0x00, 0x00, 0x07, 0xff,
				0x00, 0x00, 0x07, 0xff, 0x00, 0x00, 0x07, 0xff,
			},
		},
	}

	for i, tt := range tests {
		actual, err := expandArray(tt.compact, tt.length, tt.pad)
		if err != nil {
			t.Errorf("failed on test %d: %v", i, err)
			continue
		}

		if bytes.Compare(actual, tt.expected) != 0 {
			t.Errorf("failed on test %d: have %x want %x", i, actual, tt.expected)
			continue
		}
	}
}

func TestIndicesFromMinimal(t *testing.T) {
	tests := []struct {
		n        uint32
		k        uint32
		minimal  []byte
		expected []uint32
	}{
		{
			n: 80,
			k: 3,
			minimal: []byte{
				0x00, 0x00, 0x08, 0x00, 0x00, 0x40, 0x00, 0x02,
				0x00, 0x00, 0x10, 0x00, 0x00, 0x80, 0x00, 0x04,
				0x00, 0x00, 0x20, 0x00, 0x01,
			},
			expected: []uint32{
				1, 1, 1, 1, 1, 1, 1, 1,
			},
		},
		{
			n: 80,
			k: 3,
			minimal: []byte{
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				0xff, 0xff, 0xff, 0xff, 0xff,
			},
			expected: []uint32{
				2097151, 2097151, 2097151, 2097151, 2097151, 2097151, 2097151, 2097151,
			},
		},
		{
			n: 80,
			k: 3,
			minimal: []byte{
				0x0f, 0xff, 0xf8, 0x00, 0x20, 0x03, 0xff, 0xfe,
				0x00, 0x08, 0x00, 0xff, 0xff, 0x80, 0x02, 0x00,
				0x3f, 0xff, 0xe0, 0x00, 0x80,
			},
			expected: []uint32{
				131071, 128, 131071, 128, 131071, 128, 131071, 128,
			},
		},
		{
			n: 80,
			k: 3,
			minimal: []byte{
				0x00, 0x02, 0x20, 0x00, 0x0a, 0x7f, 0xff, 0xfe,
				0x00, 0x4d, 0x10, 0x01, 0x4c, 0x80, 0x0f, 0xfc,
				0x00, 0x00, 0x2f, 0xff, 0xff,
			},
			expected: []uint32{
				68, 41, 2097151, 1233, 665, 1023, 1, 1048575,
			},
		},
	}

	for i, tt := range tests {
		actual, err := indicesFromMinimal(tt.n, tt.k, tt.minimal)
		if err != nil {
			t.Errorf("failed on test %d: %v", i, err)
			continue
		}

		if !reflect.DeepEqual(actual, tt.expected) {
			t.Errorf("failed on test %d: have %x want %x", i, actual, tt.expected)
			continue
		}
	}
}

func TestValidSolutionFull(t *testing.T) {
	tests := []struct {
		n     uint32
		k     uint32
		input []byte
		nonce []byte
		soln  []byte
	}{
		{
			n:     96,
			k:     5,
			input: []byte("Equihash is an asymmetric PoW based on the Generalised Birthday problem."),
			nonce: []byte{
				1, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
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
		valid, err := IsValidZCashSolution(tt.n, tt.k, tt.input, tt.nonce, tt.soln)
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

func TestInvalidSolution(t *testing.T) {
	tests := []struct {
		n       uint32
		k       uint32
		input   []byte
		nonce   []byte
		indices []uint32
		error   string
	}{
		{
			n:     96,
			k:     5,
			input: []byte("Equihash is an asymmetric PoW based on the Generalised Birthday problem."),
			nonce: []byte{
				1, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
			// Change one index
			indices: []uint32{
				2262, 15185, 36112, 104243, 23779, 118390, 118332, 130041,
				32642, 69878, 76925, 80080, 45858, 116805, 92842, 111026,
				15972, 115059, 85191, 90330, 68190, 122819, 81830, 91132,
				23460, 49807, 52426, 80391, 69567, 114474, 104973, 122568,
			},
			error: "collision",
		},
		{
			n:     96,
			k:     5,
			input: []byte("Equihash is an asymmetric PoW based on the Generalised Birthday problem."),
			nonce: []byte{
				1, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
			// Swap two arbitrary indices
			indices: []uint32{
				45858, 15185, 36112, 104243, 23779, 118390, 118332, 130041,
				32642, 69878, 76925, 80080, 2261, 116805, 92842, 111026,
				15972, 115059, 85191, 90330, 68190, 122819, 81830, 91132,
				23460, 49807, 52426, 80391, 69567, 114474, 104973, 122568,
			},
			error: "collision",
		},
		{
			n:     96,
			k:     5,
			input: []byte("Equihash is an asymmetric PoW based on the Generalised Birthday problem."),
			nonce: []byte{
				1, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
			// Reverse the first pair of indices
			indices: []uint32{
				15185, 2261, 36112, 104243, 23779, 118390, 118332, 130041,
				32642, 69878, 76925, 80080, 45858, 116805, 92842, 111026,
				15972, 115059, 85191, 90330, 68190, 122819, 81830, 91132,
				23460, 49807, 52426, 80391, 69567, 114474, 104973, 122568,
			},
			error: "out of order",
		},
		{
			n:     96,
			k:     5,
			input: []byte("Equihash is an asymmetric PoW based on the Generalised Birthday problem."),
			nonce: []byte{
				1, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
			// Swap the first and second pairs of indices
			indices: []uint32{
				36112, 104243, 2261, 15185, 23779, 118390, 118332, 130041,
				32642, 69878, 76925, 80080, 45858, 116805, 92842, 111026,
				15972, 115059, 85191, 90330, 68190, 122819, 81830, 91132,
				23460, 49807, 52426, 80391, 69567, 114474, 104973, 122568,
			},
			error: "out of order",
		},
		{
			n:     96,
			k:     5,
			input: []byte("Equihash is an asymmetric PoW based on the Generalised Birthday problem."),
			nonce: []byte{
				1, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
			// Swap the second-to-last and last pairs of indices
			indices: []uint32{
				2261, 15185, 36112, 104243, 23779, 118390, 118332, 130041,
				32642, 69878, 76925, 80080, 45858, 116805, 92842, 111026,
				15972, 115059, 85191, 90330, 68190, 122819, 81830, 91132,
				23460, 49807, 52426, 80391, 104973, 122568, 69567, 114474,
			},
			error: "out of order",
		},
		{
			n:     96,
			k:     5,
			input: []byte("Equihash is an asymmetric PoW based on the Generalised Birthday problem."),
			nonce: []byte{
				1, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
			// Swap the first half and second half
			indices: []uint32{
				15972, 115059, 85191, 90330, 68190, 122819, 81830, 91132,
				23460, 49807, 52426, 80391, 69567, 114474, 104973, 122568,
				2261, 15185, 36112, 104243, 23779, 118390, 118332, 130041,
				32642, 69878, 76925, 80080, 45858, 116805, 92842, 111026,
			},
			error: "out of order",
		},
		{
			n:     96,
			k:     5,
			input: []byte("Equihash is an asymmetric PoW based on the Generalised Birthday problem."),
			nonce: []byte{
				1, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
			// Sort the indices
			indices: []uint32{
				2261, 15185, 15972, 23460, 23779, 32642, 36112, 45858,
				49807, 52426, 68190, 69567, 69878, 76925, 80080, 80391,
				81830, 85191, 90330, 91132, 92842, 104243, 104973, 111026,
				114474, 115059, 116805, 118332, 118390, 122568, 122819, 130041,
			},
			error: "collision",
		},
		{
			n:     96,
			k:     5,
			input: []byte("Equihash is an asymmetric PoW based on the Generalised Birthday problem."),
			nonce: []byte{
				1, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
			// Duplicate indices
			indices: []uint32{
				2261, 2261, 15185, 15185, 36112, 36112, 104243, 104243,
				23779, 23779, 118390, 118390, 118332, 118332, 130041, 130041,
				32642, 32642, 69878, 69878, 76925, 76925, 80080, 80080,
				45858, 45858, 116805, 116805, 92842, 92842, 111026, 111026,
			},
			error: "duplicate indices",
		},

		{
			n:     96,
			k:     5,
			input: []byte("Equihash is an asymmetric PoW based on the Generalised Birthday problem."),
			nonce: []byte{
				1, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
			// Duplicate first half
			indices: []uint32{
				2261, 15185, 36112, 104243, 23779, 118390, 118332, 130041,
				32642, 69878, 76925, 80080, 45858, 116805, 92842, 111026,
				2261, 15185, 36112, 104243, 23779, 118390, 118332, 130041,
				32642, 69878, 76925, 80080, 45858, 116805, 92842, 111026,
			},
			error: "duplicate indices",
		},
	}

	for i, tt := range tests {
		valid, err := isValidSolutionIterative(tt.n, tt.k, tt.input, tt.nonce, tt.indices)
		if err.Error() != tt.error {
			t.Errorf("failed on test %d: have %v want %s", i, err, tt.error)
			continue
		}

		if valid {
			t.Errorf("failed on test %d: have true want false", i)
			continue
		}
	}
}
