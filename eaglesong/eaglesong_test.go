package eaglesong

import (
	"bytes"
	"testing"

	"github.com/sencha-dev/powkit/internal/common/testutil"
)

func TestEaglesong(t *testing.T) {
	tests := []struct {
		input  []byte
		digest []byte
	}{
		{
			input:  []byte("Hello, world!\n"),
			digest: testutil.MustDecodeHex("0x64867e2441d162615dc2430b6bcb4d3f4b95e4d0db529fca1eece73c077d72d6"),
		},
	}

	client := NewNervos()
	for i, tt := range tests {
		digest := client.Compute(tt.input)
		if bytes.Compare(digest, tt.digest) != 0 {
			t.Errorf("failed on %d: digest mismatch: have %x, want %x", i, digest, tt.digest)
		}
	}
}
