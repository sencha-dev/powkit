# Proof of Work Algorithms

[![Go Test](https://github.com/sencha-dev/powkit/actions/workflows/go.yml/badge.svg)](https://github.com/sencha-dev/powkit/actions/workflows/go.yml)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/sencha-dev/powkit)](https://pkg.go.dev/github.com/sencha-dev/powkit?tab=doc)

# Overview

*Note: This library is still in active development and is
subject to breaking changes*

Even though there are a wide variety of Proof of Work algorithms, finding the technical details
for the implementations is quite a task. For example, both Kawpow and Firopow are variations off of ProgPow
but finding the exact differences is painful. This is meant to be a unified library to
make the specification of existing Proof of Work algorithms more standardized. 

All DAG-based algorithms only implement a light DAG, which is sufficient for validation
but not full nodes or miners. For the DAG-based algorithms, data is cached in `~/.powcache`.
Ethash will generally be between 40-80Mb per epoch (and generally 3 caches are stored). At the time of writing, running 
`make test` will throw about 800Mb of data into `~/.powcache` due to the variety and breadth of tests.

# Algorithms

| Algorithm     | DAG         | Supported |
| ------------- | ----------- | ----------
| Ethash        | yes         | yes
| Etchash       | yes         | yes
| Kawpow        | yes         | yes
| Firopow       | yes         | yes
| Octopus       | yes         | yes
| Equihash      | no          | yes
| HeavyHash     | no          | yes
| Autolykos2    | no          | yes
| Cuckoo Cycle  | no          | yes
| Eaglesong     | no          | yes
| BeamHashIII   | no          | yes
| ZelHash       | no          | yes
| Cortex        | no          | yes

# Things to Note

  - Most of these algorithms are partially optimized but I'm sure they could be improved. That being said, that will probably 
  never happen since these have never been intended to be used for miner clients. All of these algorithms far surpass 
  a reasonable threshold for performance and I have no intention of hypertuning them.
  - The base ProgPow implementation ("ProgPow094") exists in the `internal/progpow` package.
  - Since ZelHash is such a minor Equihash variant, it is treated as just "twisted Equihash" (in `equihash/`).
  - All testing is done on linux, windows support is hazy at best. 
  - The library assumes the host architecture is little-endian, I'm fairly confident big-endian architectures will not function properly.
  - As of now, the only other algorithms that are on the list of "maybes" are: [cryptonight](https://github.com/Equim-chan/cryptonight),
  [randomx](https://git.dero.io/DERO_Foundation/RandomX), X25X, and cuckatoo. 

# Roadmap

Currently, though powkit is used in production internally, it probably isn't a good idea to use yourself. The
API is still in flux and each minor version will probably be breaking. Once we do a v1.0.0 release, the structure
will probably be pretty set in stone. 

# References

  - [Ethereum: go-ethereum](https://github.com/ethereum/go-ethereum/blob/master/consensus/ethash/ethash.go)
  - [Ethereum Classic Labs: go-etchash](https://github.com/etclabscore/go-etchash)
  - [RavencoinCommunity: cpp-kawpow](https://github.com/RavenCommunity/cpp-kawpow/)
  - [Zcash: librustzcash (equihash)](https://github.com/zcash/librustzcash/tree/master/components/equihash)
  - [Firo: firo](https://github.com/firoorg/firo/tree/master/src/crypto/progpow)
  - [Ergo: ergo](https://github.com/ergoplatform/ergo/blob/0af9dd9d8846d672c1e2a77f8ab29963fa5acd1e/src/main/scala/org/ergoplatform/mining/AutolykosPowScheme.scala)
  - [leifjacky: erg-gominer-demo](https://github.com/leifjacky/erg-gominer-demo)
  - [tromp: cuckoo](https://github.com/tromp/cuckoo)
  - [Nervos Network: rfcs (eaglesong)](https://github.com/nervosnetwork/rfcs/tree/master/rfcs/0010-eaglesong)
  - [Conflux Chain: conflux-rust (Octopus)](https://github.com/Conflux-Chain/conflux-rust/tree/8fdc0773ccc447f5f6af142e84ae507284f0e411/core/src/pow)
