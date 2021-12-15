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
`make test` will throw about 2.2Gb of data into `~/.powcache` due to the variety and breadth of tests.

# Supported Algorithms

The estimated hash times are meant to be relative since
it can vary across different hardware. All tests were run
with the respective DAG already calculated. In almost all cases hashing 
functions are the bottleneck (which is how PoW algorithms
are designed to behave).

| Algorithm   | DAG         | Estimated Hash Time       |
| ----------- | ----------- | --------------------------|
| Ethash      | yes         | 4ms                       |
| Etchash     | yes         | 4ms                       |
| Kawpow      | yes         | 17ms                      |
| Firopow     | yes         | 18ms                      |
| Verthash    | yes         | 8ms disk, 270μs in memory |
| Equihash    | no          | 30μs                      |
| Autolykos2  | no          | 580μs                     |



# Roadmap

Most profitable Proof of Work chains nowdays use some sort of DAG, and that is almost always the Ethash DAG.
There are then layers on top of the Ethash DAG like ProgPow, which can be considered both a Proof of Work 
algorithm and a class of algorithms (Kawpow, Firopow). There isn't really a good way to organize all of these
into a given structure but even if there was, a new algorithm could appear tomorrow and break that structure.
powkit takes the approach of allowing most parameters to be varied and will implement new algorithms on an
as-needed basis.

Currently, though powkit is used in production internally, it probably isn't a good idea to use yourself. The
API is still in flux and each minor version will probably be breaking. Once we do a v1.0.0 release, the structure
will probably be pretty set in stone. Hopefully that will happen in the next few months.

# References

  - [Ethereum: go-ethereum](https://github.com/ethereum/go-ethereum/blob/master/consensus/ethash/ethash.go)
  - [Ethereum Classic Labs: go-etchash](https://github.com/etclabscore/go-etchash)
  - [RavencoinCommunity: cpp-kawpow](https://github.com/RavenCommunity/cpp-kawpow/)
  - [Zcash: librustzcash (equihash)](https://github.com/zcash/librustzcash/tree/master/components/equihash)
  - [Gert-Jaap Glasbergen: verthash-go](https://github.com/gertjaap/verthash-go/)
  - [Firo: firo](https://github.com/firoorg/firo/tree/master/src/crypto/progpow)
  - [Ergo: ergo](https://github.com/ergoplatform/ergo/blob/0af9dd9d8846d672c1e2a77f8ab29963fa5acd1e/src/main/scala/org/ergoplatform/mining/AutolykosPowScheme.scala)
  - [leifjacky: erg-gominer-demo](https://github.com/leifjacky/erg-gominer-demo)

