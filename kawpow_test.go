package pow

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestVerifyRVN(t *testing.T) {
	hasher, err := NewLightDag("RVN")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	const height uint64 = 1881757
	const nonce uint64 = 96461238819045
	headerHash, err := hex.DecodeString("cf63e993ca10d7b6667cc6de7c896a6f32ffe49a3916ece271744030805489a3")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	expectedMix, err := hex.DecodeString("76772038fdd6ed503752e29933d346f05e83dcbdf1939e59a45477dc3d520770")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	expectedDigest, err := hex.DecodeString("000000000000a6861f21601535f488a96c3c88f0219ad2385771787d0564b679")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	actualMix, actualDigest := hasher.Compute(headerHash, height, nonce)

	if bytes.Compare(expectedMix, actualMix) != 0 {
		t.Errorf("mixhash: rvn does not match")
	}

	if bytes.Compare(expectedDigest, actualDigest) != 0 {
		t.Errorf("digest: rvn does not match")
	}
}
