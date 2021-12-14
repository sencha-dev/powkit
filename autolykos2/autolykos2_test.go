package autolykos2

import (
	"bytes"
	"testing"

	"github.com/sencha-dev/powkit/internal/common/testutil"
)

func TestCalcN(t *testing.T) {
	tests := []struct {
		height uint64
		value  uint32
	}{
		{
			height: 500000,
			value:  67108864,
		},
		{
			height: 600000,
			value:  67108864,
		},
		{
			height: 614400,
			value:  70464240,
		},
		{
			height: 665600,
			value:  73987410,
		},
		{
			height: 700000,
			value:  73987410,
		},
		{
			height: 788400,
			value:  81571035,
		},
		{
			height: 1051200,
			value:  104107290,
		},
		{
			height: 4198400,
			value:  2143944600,
		},
		{
			height: 41984000,
			value:  2143944600,
		},
	}

	for i, tt := range tests {
		value := New(32, 26).calcN(tt.height)
		if value != tt.value {
			t.Errorf("failed on %d: have %d, want %d", i, value, tt.value)
		}
	}
}

func TestCompute(t *testing.T) {
	tests := []struct {
		msg    []byte
		nonce  uint64
		height uint64
		result []byte
	}{
		// https://www.ergoforum.org/t/test-vectors-for-increased-n-values/2887
		{
			msg:    testutil.MustDecodeHex("548c3e602a8f36f8f2738f5f643b02425038044d98543a51cabaa9785e7e864f"),
			nonce:  0x3105,
			height: 614400,
			result: testutil.MustDecodeHex("0002fcb113fe65e5754959872dfdbffea0489bf830beb4961ddc0e9e66a1412a"),
		},
	}

	for i, tt := range tests {
		result := New(32, 26).Compute(tt.msg, tt.nonce, tt.height)
		if bytes.Compare(result, tt.result) != 0 {
			t.Errorf("failed on %d: have %x, want %x", i, result, tt.result)
		}
	}
}
