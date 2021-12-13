// ethash: C/C++ implementation of Ethash, the Ethereum Proof of Work algorithm.
// Copyright 2018-2019 Pawel Bylica.
// Licensed under the Apache License, Version 2.0.

package progpow

import (
	"testing"
)

func TestRandomMath(t *testing.T) {
	tests := []struct {
		a        uint32
		b        uint32
		selector uint32
		value    uint32
	}{
		{
			a:        0x8626BB1F,
			b:        0xBBDFBC4E,
			selector: 0x883E5B49,
			value:    0x4206776D,
		},
		{
			a:        0x3F4BDFAC,
			b:        0xD79E414F,
			selector: 0x36B71236,
			value:    0x4C5CB214,
		},
		{
			a:        0x6D175B7E,
			b:        0xC4E89D4C,
			selector: 0x944ECABB,
			value:    0x53E9023F,
		},
		{
			a:        0x2EDDD94C,
			b:        0x7E70CB54,
			selector: 0x3F472A85,
			value:    0x2EDDD94C,
		},
		{
			a:        0x61AE0E62,
			b:        0xe0596b32,
			selector: 0x3F472A85,
			value:    0x61AE0E62,
		},
		{
			a:        0x8A81E396,
			b:        0x3F4BDFAC,
			selector: 0xCEC46E67,
			value:    0x1E3968A8,
		},
		{
			a:        0x8A81E396,
			b:        0x7E70CB54,
			selector: 0xDBE71FF7,
			value:    0x1E3968A8,
		},
		{
			a:        0xA7352F36,
			b:        0xA0EB7045,
			selector: 0x59E7B9D8,
			value:    0xA0212004,
		},
		{
			a:        0xC89805AF,
			b:        0x64291E2F,
			selector: 0x1BDC84A9,
			value:    0xECB91FAF,
		},
		{
			a:        0x760726D3,
			b:        0x79FC6A48,
			selector: 0xC675CAC5,
			value:    0x0FFB4C9B,
		},
		{
			a:        0x75551D43,
			b:        0x3383BA34,
			selector: 0x2863AD31,
			value:    0x00000003,
		},
		{
			a:        0xEA260841,
			b:        0xE92C44B7,
			selector: 0xF83FFE7D,
			value:    0x0000001B,
		},
	}

	for i, tt := range tests {
		value := randomMath(tt.a, tt.b, tt.selector)
		if value != tt.value {
			t.Errorf("failed on %d: have %d, want %d", i, value, tt.value)
		}
	}
}

func TestRandomMerge(t *testing.T) {
	tests := []struct {
		a        uint32
		b        uint32
		selector uint32
		value    uint32
	}{
		{
			a:        0x3B0BB37D,
			b:        0xA0212004,
			selector: 0x9BD26AB0,
			value:    0x3CA34321,
		},
		{
			a:        0x10C02F0D,
			b:        0x870FA227,
			selector: 0xD4F45515,
			value:    0x91C1326A,
		},
		{
			a:        0x24D2BAE4,
			b:        0x0FFB4C9B,
			selector: 0x7FDBC2F2,
			value:    0x2EDDD94C,
		},
		{
			a:        0xDA39E821,
			b:        0x089C4008,
			selector: 0x8B6CD8C3,
			value:    0x8A81E396,
		},
	}

	for i, tt := range tests {
		value := randomMerge(tt.a, tt.b, tt.selector)
		if value != tt.value {
			t.Errorf("failed on %d: have %d, want %d", i, value, tt.value)
		}
	}
}
