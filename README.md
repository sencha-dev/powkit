# Go PoW

This is a library designed for Proof of Work validation
for chains that require a DAG - generally this refers
to Ethash or ProgPOW chains. Currently `ethereum`, 
`ethereum classic`, and `ravencoin`. 

The library is built for mining pools, so it only
returns the mix and digest, leaving the difficulty
validation for the client. 

The DAG code is heavily borrowing from [go-etchash](https://github.com/etclabscore/go-etchash)
(which was actually inspired by [ethashproof](https://github.com/tranvictor/ethashproof)).
The Kawpow code is mostly ported over from [cpp-kawpow](https://github.com/RavenCommunity/cpp-kawpow/).
There are other small bits taken from elsewhere with sources in the code.

This library was created due to the lack of a Kawpow verification library
in native Go. Since the DAG generation component of Kawpow varies only 
slightly from Ethash, `go-etchash` was generalized to handle the different
requirements of DAG generation and implemented the native algorithm for Kawpow.

For `ETC` and `RVN`, this will only work above the blocks for ECIP-1099 (`11700000` ) and the 
Kawpow hard fork (`1219736`), respectively.

Finally, this library only implements a light DAG - the full DAG is very large and there
isn't really a use case for mining pools to use it. It wouldn't be too difficult to 
implement though, it exists in `go-etchash` so it would just be a matter of adding
the ProgPOW full DAG implementation.

# Todos

- [X] Fix the Ravencoin algorithm to function properly
- [X] Implement minimum block heights
- [ ] Handle the L1 cache in ProgPOW more elegantly (memory map if needed)
- [ ] Clean up Ravencoin generally to remove the hacky solutions currently implemented 
- [ ] Remove ethereum dependency
- [ ] Utility functions for difficulty validation
- [ ] More extensive testing for `ETH`, `ETC`, and `RVN`
- [ ] Implement testnet support

# Usage

To instantiate the DAG (*the only error `NewLightDag` 
will return is if the chain symbol is not supported*).

```go
hasher, err := NewLightDag("ETH")
if err != nil {
	panic(err)
}

hasher, err := NewLightDag("ETC")
if err != nil {
	panic(err)
}

hasher, err := NewLightDag("RVN")
if err != nil {
	panic(err)
}
```

To compute a mix and digest (*the only error `Compute`
will return is if the height is below the chain's 
minimum height*).

```go
mix, digest, err := hasher.Compute(hash, height, nonce)
```


Example (ETH):

```go
package main

import (
	"encoding/hex"
	"fmt"

	"github.com/sencha-dev/go-pow"
)

func main() {
	nonce := uint64(5819316201154249538)
	height := uint64(12738427)
	hash, err := hex.DecodeString("28dcbf10a1cb49eb61f2e8b1b66636b46ea122dc6176de423f89ee3afd1467f4")
	if err != nil {
		panic(err)
	}

	hasher, err := pow.NewLightDag("ETH")
	if err != nil {
		panic(err)
	}

	mix, digest, err := hasher.Compute(hash, height, nonce)
	if err != nil {
		panic(err)
	}

	fmt.Println(hex.EncodeToString(mix), hex.EncodeToString(digest))
}
```

Example (RVN):

```go
package main

import (
	"encoding/hex"
	"fmt"

	"github.com/sencha-dev/go-pow"
)

func main() {
	height := uint64(1881757)
	nonce := uint64(96461238819045)
	headerHash, err := hex.DecodeString("cf63e993ca10d7b6667cc6de7c896a6f32ffe49a3916ece271744030805489a3")
	if err != nil {
		panic(err)
	}

	hasher, err := pow.NewLightDag("RVN")
	if err != nil {
		panic(err)
	}

	mix, digest, err := hasher.Compute(headerHash, height, nonce)
	if err != nil {
		panic(err)
	}

	fmt.Println(hex.EncodeToString(mix), hex.EncodeToString(digest))
}
```
