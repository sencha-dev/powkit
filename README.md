# Proof of Work Algorithms

[![Go Test](https://github.com/sencha-dev/go-pow/actions/workflows/go.yml/badge.svg)](https://github.com/sencha-dev/go-pow/actions/workflows/go.yml)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/sencha-dev/go-pow)](https://pkg.go.dev/github.com/sencha-dev/go-pow?tab=doc)

# Overview

*Note: This library is still in active development and is
subject to breaking changes*

Though there is a wide variety of Proof of Work algorithms, finding the technical details
for the implementations is quite a task. Both Kawpow and Firopow are variations off of ProgPow,
though finding the exact differences is no easy task. This is meant to be a unified library to
make the implementation of existing Proof of Work algorithms easier. 

All DAG-based algorithms only implement a light DAG, which is sufficient for verification
but not full nodes or miners. For the DAG-based algorithms and verthash, data is cached in `~/.powcache`.
Ethash will generally be between 40-80Mb per epoch (and generally 3 caches are stored), but verthash
requires a strict 1.2Gb, so be careful if you're using verthash in memory. At the time of writing, running 
`make test` will throw about 3.3Gb of data into `~/.powcache` due to the variety and breadth of tests.

# Algorithms

| Algorithm     | DAG         | Supported |
| ------------- | ----------- | ----------
| Ethash        | yes         | yes
| Etchash       | yes         | yes
| Kawpow        | yes         | yes
| Firopow       | yes         | yes
| Octopus       | yes         | yes
| Verthash      | yes         | yes
| Equihash      | no          | yes
| Autolykos2    | no          | yes
| Cuckoo Cycle  | no          | yes
| Eaglesong     | no          | yes
| BeamHashIII   | no          | no

# Things to Note

  - Cuckoo Cycle is built specifically for Aeternity. There is a modification of the `sipnode` function in the current version
  of tromp's cuckoo algorithms that Aeternity does not use (a Nicehash dev gives more details [here](https://forum.aeternity.com/t/support-aeternity-stratum-implementation/3140/6)). It wouldn't be hard to implement other Cuckoo Cycle algorithms (cuckatoo, cuckaroo),
  there just isn't really a need at this point since Grin is fairly annoying. BlockCypher implements the other algorithms [here](https://github.com/blockcypher/libgrin/tree/master/core/pow).
  - Equihash is built around ZCash's variation of Equihash. The original implementation is left for compatibility reasons, hopefully one day
  I'll find a way to unify the two (though this may not be possible in a reasonable way). 
  - All non-DAG algorithms are less organized than I would like, they'll probably be overhauled at some point for a more coherent general standard.
  - All testing is done on linux, windows support is hazy at best. 
  - The library assumes the host architecture is little-endian, I'm fairly confident big-endian architectures will not function properly.
  - The base ProgPow implementation ("ProgPow094") exists in the `internal/progpow` package.
  - One day I'll implement BeamHashIII, it is just a slight modification of Equihash (I think 150_5?). Other than that, no other
  algorithms are planned - there is a [cryptonight](https://github.com/Equim-chan/cryptonight) and a 
  [randomx](https://git.dero.io/DERO_Foundation/RandomX) implementation in Go, though these aren't really of interest. 
  Of course, this can change if new algorithms become popular.

# Roadmap

Most profitable Proof of Work chains nowdays use some sort of DAG, and that is almost always the Ethash DAG.
There are then layers on top of the Ethash DAG like ProgPow, which can be considered both a Proof of Work 
algorithm and a class of algorithms (Kawpow, Firopow). There isn't really a good way to organize all of these
into a given structure but even if there was, a new algorithm could appear tomorrow and break that structure.
powkit takes the approach of allowing most parameters to be varied and will implement new algorithms on an
as-needed basis.

Currently, though powkit is used in production internally, it probably isn't a good idea to use yourself. The
API is still in flux and each minor version will probably be breaking. Once we do a v1.0.0 release, the structure
will probably be pretty set in stone. 

# References

  - [Ethereum: go-ethereum](https://github.com/ethereum/go-ethereum/blob/master/consensus/ethash/ethash.go)
  - [Ethereum Classic Labs: go-etchash](https://github.com/etclabscore/go-etchash)
  - [RavencoinCommunity: cpp-kawpow](https://github.com/RavenCommunity/cpp-kawpow/)
  - [Zcash: librustzcash (equihash)](https://github.com/zcash/librustzcash/tree/master/components/equihash)
  - [Gert-Jaap Glasbergen: verthash-go](https://github.com/gertjaap/verthash-go/)
  - [Firo: firo](https://github.com/firoorg/firo/tree/master/src/crypto/progpow)
  - [Ergo: ergo](https://github.com/ergoplatform/ergo/blob/0af9dd9d8846d672c1e2a77f8ab29963fa5acd1e/src/main/scala/org/ergoplatform/mining/AutolykosPowScheme.scala)
  - [leifjacky: erg-gominer-demo](https://github.com/leifjacky/erg-gominer-demo)
  - [tromp: cuckoo](https://github.com/tromp/cuckoo)
  - [Nervos Network: rfcs (eaglesong)](https://github.com/nervosnetwork/rfcs/tree/master/rfcs/0010-eaglesong)
  - [Conflux Chain: conflux-rust (Octopus)](https://github.com/Conflux-Chain/conflux-rust/tree/8fdc0773ccc447f5f6af142e84ae507284f0e411/core/src/pow)
