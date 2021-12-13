package verthash

import (
	"bytes"
	"testing"

	"github.com/sencha-dev/powkit/internal/common/testutil"
)

func TestGenerate(t *testing.T) {
	New(false)
}

func TestVerify(t *testing.T) {
	tests := []struct {
		input  []byte
		output []byte
	}{
		{
			input:  testutil.MustDecodeHex("000000203a297b4b7685170d7644b43e5a6056234cc2414edde454a87580e1967d14c1078c13ea916117b0608732f3f65c2e03b81322efc0a62bcee77d8a9371261970a58a5a715da80e031b02560ad8"),
			output: testutil.MustDecodeHex("E0F6C10B4A38F35A6CDCC26D32A7ED8C3BFC5D827A9BC72647AFA324B70D0463"),
		},
	}

	verthash, err := New(true)
	if err != nil {
		t.Errorf("failed on new verthash: %v", err)
		return
	}
	defer verthash.Close()

	for i, tt := range tests {
		h := verthash.Hash(tt.input)

		if !bytes.Equal(h[:], tt.output) {
			t.Errorf("mismatch on verthash verify test %d: have %x want %x", i, h, tt.output)
		}
	}
}
