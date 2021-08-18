package pow_test

import (
	"encoding/hex"
	"fmt"

	"github.com/sencha-dev/go-pow"
)

// ETH example of generating a PoW hash and mix digest
func ExampleETH() {
	nonce := uint64(0x956e895d988798e)
	height := uint64(12965001)
	hash, err := hex.DecodeString("cf133ce0cccd4ad877d671b310c27f5ce19c28c14455dac45b90171bac5581c7")
	if err != nil {
		panic(err)
	}

	dag, err := pow.NewLightDag("ETH")
	if err != nil {
		panic(err)
	}

	mix, digest, err := dag.Compute(hash, height, nonce)
	if err != nil {
		panic(err)
	}

	fmt.Println(hex.EncodeToString(mix), hex.EncodeToString(digest))
	// Output: cb3166ebb1888430069b769145b20ba5e3a55f32fd2fa39f0ebdc08d60b4557e 00000000000000012923a9ab2605573e0158adeb21c86b22d8ebd33b8ee08856
}

// ETC example of generating a PoW hash and mix digest
func ExampleETC() {
	nonce := uint64(0x9827862e22a92ff1)
	height := uint64(13344137)
	hash, err := hex.DecodeString("27eaf677273c9147cd27b99c34b3783243255864a54b169af238750c39b3c167")
	if err != nil {
		panic(err)
	}

	dag, err := pow.NewLightDag("ETC")
	if err != nil {
		panic(err)
	}

	mix, digest, err := dag.Compute(hash, height, nonce)
	if err != nil {
		panic(err)
	}

	fmt.Println(hex.EncodeToString(mix), hex.EncodeToString(digest))
	// Output: 6dd0879bfe248c4ac73160a3d2554ce12431d2033b5f4464559368d855795df7 000000000000ca4bd2875398d73ab24e8467d5986bfe85ec7fca1b860a540d14
}

// RVN example of generating a block hash and mix digest
func ExampleRVN() {
	height := uint64(1888509)
	nonce := uint64(0xf09b0e1342275f3f)
	headerHash, err := hex.DecodeString("14f2c18d74d48abe637437458c10ff5283a9a5197e8b5e740a161f4411b97a43")
	if err != nil {
		panic(err)
	}

	dag, err := pow.NewLightDag("RVN")
	if err != nil {
		panic(err)
	}

	mix, digest, err := dag.Compute(headerHash, height, nonce)
	if err != nil {
		panic(err)
	}

	fmt.Println(hex.EncodeToString(mix), hex.EncodeToString(digest))
	// Output: 3dedcc6fb28c6bf3f5d203f29188bef3ff86688be34c93f28bd227eced9226e4 0000000000005e6585e5e6ab7e4d75a98810204341def05823ad3a5ca1fa0d83
}
