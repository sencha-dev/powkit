package pow

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestVerifyETH(t *testing.T) {
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
	hash := MustDecodeHex("28dcbf10a1cb49eb61f2e8b1b66636b46ea122dc6176de423f89ee3afd1467f4")
	expectedMix := MustDecodeHex("31bfa4e5a088d322c4d4704326fa77024414dc2f44db4d418b245cf53ce1819f")

	actualMix, actualDigest, err := hasher.Compute(hash, height, nonce)
	if err != nil {
		panic(err)
	}

	if bytes.Compare(actualMix, expectedMix) != 0 {
		t.Errorf("verify: mismatch on mixhash verify for ETH")
	}

	target := new(big.Int).Div(maxUint256, difficulty)
	if common.BytesToHash(actualDigest).Big().Cmp(target) > 0 {
		t.Errorf("verify: digest does not meet target difficulty")
	}
}

/*
func TestVerifyRVN(t *testing.T) {
	hasher, err := NewLightDag("RVN")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	const height uint64 = 1881757
	const nonce uint64 = 96461238819045
	headerHash := MustDecodeHex("cf63e993ca10d7b6667cc6de7c896a6f32ffe49a3916ece271744030805489a3")
	expectedMix := MustDecodeHex("76772038fdd6ed503752e29933d346f05e83dcbdf1939e59a45477dc3d520770")
	expectedDigest := MustDecodeHex("000000000000a6861f21601535f488a96c3c88f0219ad2385771787d0564b679")

	actualMix, actualDigest, err := hasher.Compute(headerHash, height, nonce)
	if err != nil {
		panic(err)
	}

	if bytes.Compare(expectedMix, actualMix) != 0 {
		t.Errorf("mixhash: rvn does not match")
	}

	if bytes.Compare(expectedDigest, actualDigest) != 0 {
		t.Errorf("digest: rvn does not match")
	}
}*/
