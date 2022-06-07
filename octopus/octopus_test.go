package octopus

import (
	"bytes"
	"testing"

	"github.com/sencha-dev/powkit/internal/common/testutil"
)

func TestComputeConflux(t *testing.T) {
	tests := []struct {
		height uint64
		nonce  uint64
		hash   []byte
		digest []byte
	}{
		{
			height: 2,
			nonce:  151182848800,
			hash:   testutil.MustDecodeHex("0x4d99d0b41c7eb0dd1a801c35aae2df28ae6b53bc7743f0818a34b6ec97f5b4ae"),
			digest: testutil.MustDecodeHex("0xd45c965d3707e27a42995132637854234385cbf5626897259f1ee980554ddd5c"),
		},
		{
			height: 45749306,
			nonce:  0xd653b35d4689284f,
			hash:   testutil.MustDecodeHex("0x2115dd73ee8e3e15e65d218eedd6846514ac782b636bc5943fbe4f980e2d395b"),
			digest: testutil.MustDecodeHex("0x00000001b739e70d59c01c97652fa4c9540b1028cee284a72cb7c13bcab7536f"),
		},
		{
			height: 45749306,
			nonce:  0xe437cce8dbe11c39,
			hash:   testutil.MustDecodeHex("0x2115dd73ee8e3e15e65d218eedd6846514ac782b636bc5943fbe4f980e2d395b"),
			digest: testutil.MustDecodeHex("0x000000005b182d14a581cd3388cf2621fc787b416b5e0a3638ff8adc466e9686"),
		},
	}

	client := NewConflux()
	for i, tt := range tests {
		digest := client.Compute(tt.height, tt.nonce, tt.hash)
		if bytes.Compare(digest, tt.digest) != 0 {
			t.Errorf("failed on %d: digest mismatch: have %x, want %x", i, digest, tt.digest)
		}
	}
}
