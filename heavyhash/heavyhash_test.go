package heavyhash

import (
	"bytes"
	"testing"

	"github.com/sencha-dev/powkit/internal/common/testutil"
)

func TestHeavyHash(t *testing.T) {
	tests := []struct {
		digest    []byte
		nonce     uint64
		timestamp int64
		hash      []byte
	}{
		{
			hash:      testutil.MustDecodeHex("81553a695a0588998c413792e74ce8b8f8a096d64b3ee47387372434485c0b6f"),
			nonce:     0x2f8400000eba167c,
			timestamp: 0x000001848ca87c49,
			digest:    testutil.MustDecodeHex("000000001726686e851f02c584d7cc8a8fbe5938ecdb3ffa2ba16c260ee1fc40"),
		},
		{
			hash:      testutil.MustDecodeHex("9785c4d0e244b3564115fd110e8e608a688b8803baab9fa6948e9f7ba0540f4c"),
			nonce:     0x2f8400000ffdd00f,
			timestamp: 0x000001848cac94a5,
			digest:    testutil.MustDecodeHex("000000000943f66d7c28611e552e5523bbeb5e61c52d38b86924ca6268269331"),
		},
		{
			hash:      testutil.MustDecodeHex("2f6cea927f6dca4357c9aab4ac8b7957e3a349cc317d4631f26593b90f327256"),
			nonce:     0x2f84000005d52b98,
			timestamp: 0x000001848cae8db5,
			digest:    testutil.MustDecodeHex("0000000018c4ec84524ee859f660392858fa50b12c1419f36f480fedcd6ee3d8"),
		},
	}

	for i, tt := range tests {
		digest := heavyHash(tt.hash, tt.timestamp, tt.nonce)
		if bytes.Compare(digest, tt.digest) != 0 {
			t.Errorf("failed on %d: have %x, want %x", i, digest, tt.digest)
		}
	}
}
