# Go PoW

[![Go Test](https://github.com/sencha-dev/go-pow/actions/workflows/go.yml/badge.svg)](https://github.com/sencha-dev/go-pow/actions/workflows/go.yml)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/sencha-dev/go-pow)](https://pkg.go.dev/github.com/sencha-dev/go-pow?tab=doc)

*Note: This library is still in active development and is
subject to breaking changes*

This is a library designed for Proof of Work validation for chains that require a DAG - generally 
this refers to Ethash or ProgPOW chains. 

This was created due to the lack of a Kawpow verification library
in native Go. Since the DAG generation component of Kawpow varies only 
slightly from Ethash, `go-etchash` was generalized to handle the different
requirements of DAG generation and updated to implement Kawpow. Some more significant
changes have been made, such as memory mapping the L1 cache for the Kawpow DAG.
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
