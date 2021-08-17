# Go PoW

This is a library designed for Proof of Work validation
for chains that require a DAG - generally this refers
to Ethash or ProgPOW chains. Currently `ethereum`, 
`ethereum classic`, and `ravencoin` are supported. 

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
- [X] Remove ethereum dependency
- [X] Memory map L1 cache for RVN to remove 20ms on `Compute`
- [X] Clean up Ravencoin generally to remove the incomplete solutions 
- [X] More extensive testing for `ETH`, `ETC`, and `RVN`
- [ ] Make file cleanup more consistent
- [ ] Utility functions for difficulty validation
- [ ] Implement testnet support

# Usage

To instantiate the DAG (*the only error `NewLightDag` 
will return is if the chain symbol is not supported*).

```go
dag, err := NewLightDag("ETH")
if err != nil {
	panic(err)
}

dag, err := NewLightDag("ETC")
if err != nil {
	panic(err)
}

dag, err := NewLightDag("RVN")
if err != nil {
	panic(err)
}
```

To compute a mix and digest (*the only error `Compute`
will return is if the height is below the chain's 
minimum height or the input hash is invalid*).

```go
mix, digest, err := dag.Compute(hash, height, nonce)
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
}
```

Example (ETC):

```go
package main

import (
	"encoding/hex"
	"fmt"

	"github.com/sencha-dev/go-pow"
)

func main() {
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
}
```
