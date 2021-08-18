# Go PoW

[![Go Test](https://github.com/sencha-dev/go-pow/actions/workflows/go.yml/badge.svg)](https://github.com/sencha-dev/go-pow/actions/workflows/go.yml)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/sencha-dev/go-pow)](https://pkg.go.dev/github.com/sencha-dev/go-pow?tab=doc)

This is a library designed for Proof of Work validation for chains that require a DAG - generally 
this refers to Ethash or ProgPOW chains. Currently `ethereum`, `ethereum classic`, 
and `ravencoin` are supported. Target difficulty validation is left to the client because mining
pools generally do two difficulty checks (share difficulty and block difficulty). Eventually
there will probably be some difficulty validation utilities in here though.

This was created due to the lack of a Kawpow verification library
in native Go. Since the DAG generation component of Kawpow varies only 
slightly from Ethash, `go-etchash` was generalized to handle the different
requirements of DAG generation and implemented the native algorithm for Kawpow.
Before `go-pow`, the only clean solution for Kawpow hashing was using the RPC
call `getkawpowhash` from the full node.

The DAG code is heavily borrowing from [go-etchash](https://github.com/etclabscore/go-etchash)
(which was actually inspired by [ethashproof](https://github.com/tranvictor/ethashproof)).
The Kawpow code is mostly ported over from [cpp-kawpow](https://github.com/RavenCommunity/cpp-kawpow/).
There are other small bits taken from elsewhere with sources in the code.

Finally, this library only implements a light DAG - the full DAG is very large and there
isn't really a use case for pools to use it during validation. It wouldn't be too difficult to 
implement though, it exists in `go-etchash` so it would just be a matter of adding
the ProgPOW full DAG implementation.

*note: for `ETC` and `RVN`, this will only work above the blocks for ECIP-1099 (`11700000` ) and the 
Kawpow hard fork (`1219736`)*

# Feature Todos

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
minimum height*).

```go
mix, digest, err := dag.Compute(hash, height, nonce)
```
