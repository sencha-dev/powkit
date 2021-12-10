// ethash: C/C++ implementation of Ethash, the Ethereum Proof of Work algorithm.
// Copyright 2018-2019 Pawel Bylica.
// Licensed under the Apache License, Version 2.0.

package progpow

import (
	"testing"
)

func TestRandomMath(t *testing.T) {
	cases := [][4]uint32{
		[4]uint32{0x8626BB1F, 0xBBDFBC4E, 0x883E5B49, 0x4206776D},
		[4]uint32{0x3F4BDFAC, 0xD79E414F, 0x36B71236, 0x4C5CB214},
		[4]uint32{0x6D175B7E, 0xC4E89D4C, 0x944ECABB, 0x53E9023F},
		[4]uint32{0x2EDDD94C, 0x7E70CB54, 0x3F472A85, 0x2EDDD94C},
		[4]uint32{0x61AE0E62, 0xe0596b32, 0x3F472A85, 0x61AE0E62},
		[4]uint32{0x8A81E396, 0x3F4BDFAC, 0xCEC46E67, 0x1E3968A8},
		[4]uint32{0x8A81E396, 0x7E70CB54, 0xDBE71FF7, 0x1E3968A8},
		[4]uint32{0xA7352F36, 0xA0EB7045, 0x59E7B9D8, 0xA0212004},
		[4]uint32{0xC89805AF, 0x64291E2F, 0x1BDC84A9, 0xECB91FAF},
		[4]uint32{0x760726D3, 0x79FC6A48, 0xC675CAC5, 0x0FFB4C9B},
		[4]uint32{0x75551D43, 0x3383BA34, 0x2863AD31, 0x00000003},
		[4]uint32{0xEA260841, 0xE92C44B7, 0xF83FFE7D, 0x0000001B},
	}

	for i, c := range cases {
		actual := randomMath(c[0], c[1], c[2])
		if actual != c[3] {
			t.Errorf("failed randomMath test on case %d", i)
		}
	}
}

func TestRandomMerge(t *testing.T) {
	cases := [][4]uint32{
		[4]uint32{0x3B0BB37D, 0xA0212004, 0x9BD26AB0, 0x3CA34321},
		[4]uint32{0x10C02F0D, 0x870FA227, 0xD4F45515, 0x91C1326A},
		[4]uint32{0x24D2BAE4, 0x0FFB4C9B, 0x7FDBC2F2, 0x2EDDD94C},
		[4]uint32{0xDA39E821, 0x089C4008, 0x8B6CD8C3, 0x8A81E396},
	}

	for i, c := range cases {
		actual := randomMerge(c[0], c[1], c[2])
		if actual != c[3] {
			t.Errorf("failed randomMerge test on case %d", i)
		}
	}
}
