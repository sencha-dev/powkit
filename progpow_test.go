/*
test vectors taken from ProgPOW spec
github.com/ifdefelse/ProgPOW/blob/master/test-vectors
*/

package pow

/*
import (
	"testing"
)

func TestFnv1a(t *testing.T) {
	var h1, d1, h2, d2, h3, d3 uint32
	var expected1, expected2, expected3 uint32

	h1, d1, expected1 = 0x811C9DC5, 0xDDD0A47B, 0xD37EE61A
	h2, d2, expected2 = 0xD37EE61A, 0xEE304846, 0xDEDC7AD4
	h3, d3, expected3 = 0xDEDC7AD4, 0x00000000, 0xA9155BBC

	actual1 := fnv1a(h1, d1)
	if actual1 != expected1 {
		t.Errorf("failed fnv1a test 1")
	}

	actual2 := fnv1a(h2, d2)
	if actual2 != expected2 {
		t.Errorf("failed fnv1a test 2")
	}

	actual3 := fnv1a(h3, d3)
	if actual3 != expected3 {
		t.Errorf("failed fnv1a test 3")
	}
}

func TestKiss99(t *testing.T) {
	var z, w, jsr, jcong uint32
	z, w = 362436069, 521288629
	jsr, jcong = 123456789, 380116160

	results := map[int]uint32{
		1:      769445856,
		2:      742012328,
		3:      2121196314,
		4:      2805620942,
		100000: 941074834,
	}

	kiss := NewKiss99(z, w, jsr, jcong)

	for i := 1; i < 100001; i++ {
		actual := kiss.Next()

		if expected, ok := results[i]; ok {
			if actual != expected {
				t.Errorf("failed Kiss99 test on iteration %d", i)
			}
		}
	}
}

func TestInitMixRngState(t *testing.T) {
	const period uint64 = 50
	const height uint64 = 30000

	expectedSrc := []byte{
		0x1A, 0x1E, 0x01, 0x13, 0x0B, 0x15, 0x0F, 0x12,
		0x03, 0x11, 0x1F, 0x10, 0x1C, 0x04, 0x16, 0x17,
		0x02, 0x0D, 0x1D, 0x18, 0x0A, 0x0C, 0x05, 0x14,
		0x07, 0x08, 0x0E, 0x1B, 0x06, 0x19, 0x09, 0x00,
	}

	expectedDst := []byte{
		0x00, 0x04, 0x1B, 0x1A, 0x0D, 0x0F, 0x11, 0x07,
		0x0E, 0x08, 0x09, 0x0C, 0x03, 0x0A, 0x01, 0x0B,
		0x06, 0x10, 0x1C, 0x1F, 0x02, 0x13, 0x1E, 0x16,
		0x1D, 0x05, 0x18, 0x12, 0x19, 0x17, 0x15, 0x14,
	}

	const expectedZ uint32 = 0x6535921C
	const expectedW uint32 = 0x29345B16
	const expectedJsr uint32 = 0xC0DD7F78
	const expectedJcong uint32 = 0x1165D7EB

	number := height / period
	state := initMixRngState(number)

	for i := 0; i < len(expectedSrc); i++ {
		if state.SrcSeq[i] != uint32(expectedSrc[i]) {
			t.Errorf("failed initMixRngState test for SrcSeq value at index %d", i)
		}

		if state.DstSeq[i] != uint32(expectedDst[i]) {
			t.Errorf("failed initMixRngState test for DstSeq value at index %d", i)
		}
	}

	if state.Rng.z != expectedZ {
		t.Errorf("failed initMixRngState test for Rng value z")
	}

	if state.Rng.w != expectedW {
		t.Errorf("failed initMixRngState test for Rng value w")
	}

	if state.Rng.jsr != expectedJsr {
		t.Errorf("failed initMixRngState test for Rng value jsr")
	}

	if state.Rng.jcong != expectedJcong {
		t.Errorf("failed initMixRngState test for Rng value jcong")
	}

}

/*
func TestInitMix(t *testing.T) {
	seed := [2]uint32{0xEE304846, 0xDDD0A47B}

	lanesExpected := map[int][32]uint32{
		0: [32]uint32{
			0x10C02F0D, 0x99891C9E, 0xC59649A0, 0x43F0394D,
			0x24D2BAE4, 0xC4E89D4C, 0x398AD25C, 0xF5C0E467,
			0x7A3302D6, 0xE6245C6C, 0x760726D3, 0x1F322EE7,
			0x85405811, 0xC2F1E765, 0xA0EB7045, 0xDA39E821,
			0x79FC6A48, 0x089E401F, 0x8488779F, 0xD79E414F,
			0x041A826B, 0x313C0D79, 0x10125A3C, 0x3F4BDFAC,
			0xA7352F36, 0x7E70CB54, 0x3B0BB37D, 0x74A3E24A,
			0xCC37236A, 0xA442B311, 0x955AB27A, 0x6D175B7E,
		},
		13: [32]uint32{
			0x4E46D05D, 0x2E77E734, 0x2C479399, 0x70712177,
			0xA75D7FF5, 0xBEF18D17, 0x8D42252E, 0x35B4FA0E,
			0x462C850A, 0x2DD2B5D5, 0x5F32B5EC, 0xED5D9EED,
			0xF9E2685E, 0x1F29DC8E, 0xA78F098B, 0x86A8687B,
			0xEA7A10E7, 0xBE732B9D, 0x4EEBCB60, 0x94DD7D97,
			0x39A425E9, 0xC0E782BF, 0xBA7B870F, 0x4823FF60,
			0xF97A5A1C, 0xB00BCAF4, 0x02D0F8C4, 0x28399214,
			0xB4CCB32D, 0x83A09132, 0x27EA8279, 0x3837DDA3,
		},
	}

	mix := initMix(seed)

	for lane := range mix {
		if expected, ok := lanesExpected[lane]; ok {
			actual := mix[lane]

			for reg := range expected {
				if actual[reg] != expected[reg] {
					t.Errorf("failed Kiss99 test on iteration %d", lane)
					continue
				}
			}
		}
	}
}

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
}*/
