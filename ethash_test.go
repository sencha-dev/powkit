package pow

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestEthVerify(t *testing.T) {
	hasher, err := NewLightDag("ETH")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	difficulty := big.NewInt(8_726_282_760)
	if difficulty.Cmp(common.Big0) == 0 {
		t.Errorf("invalid block difficulty")
		return
	}

	nonce := uint64(5819316201154249538)
	height := uint64(12738427)
	hash, err := hex.DecodeString("28dcbf10a1cb49eb61f2e8b1b66636b46ea122dc6176de423f89ee3afd1467f4")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	expectedMix, err := hex.DecodeString("31bfa4e5a088d322c4d4704326fa77024414dc2f44db4d418b245cf53ce1819f")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	actualMix, actualDigest := hasher.Compute(hash, height, nonce)

	if bytes.Compare(actualMix, expectedMix) != 0 {
		t.Errorf("verify: mismatch on mixhash verify for ETH")
	}

	target := new(big.Int).Div(maxUint256, difficulty)
	if common.BytesToHash(actualDigest).Big().Cmp(target) > 0 {
		t.Errorf("verify: digest does not meet target difficulty")
	}
}
