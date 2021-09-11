# Proof of Work Algorithms

[![Go Test](https://github.com/sencha-dev/go-pow/actions/workflows/go.yml/badge.svg)](https://github.com/sencha-dev/go-pow/actions/workflows/go.yml)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/sencha-dev/go-pow)](https://pkg.go.dev/github.com/sencha-dev/go-pow?tab=doc)

*Note: This library is still in active development and is
subject to breaking changes*

This libary initially focused on Proof of Work algorithms that required DAG generation,
but has expanded beyond that scope. Currently it supports `ethash` (Ethereum), `etchash`
(Ethereum Classic), `kawpow` (Ravencoin), `equihash` (Zcash), and `verthash` (Vertcoin).
`cuckoo/cuckaroo` (Grin, Aeternity), `autolykos2` (Ergo), and `BeamHashIII` (Beam) will
be added eventually.

All DAG-based algorithms only implement a light DAG, which is sufficient for verification
but not mining. For the DAG-based algorithms and verthash, data is cached in `~/.powcache`.
Ethash will generally be between 40-80Mb per epoch (and we cache 3 epochs), but verthash
requires a strict 1.2Gb, so be careful if you're using verthash in memory. 

# References

*Note: I have done my best to keep the copyright notices in each given file, though
some libraries like `cpp-kawpow` were forked three or four times without license headers,
making it fairly difficult to keep track. I'm always happy to update them if I've made
a mistake.*

  - [Victor Tran: ethashproof](https://github.com/tranvictor/ethashproof)
  - [Ethereum Classic Labs: go-etchash](https://github.com/etclabscore/go-etchash)
  - [RavencoinCommunity: cpp-kawpow](https://github.com/RavenCommunity/cpp-kawpow/)
  - [Zcash: librustzcash (equihash)](https://github.com/zcash/librustzcash/tree/master/components/equihash)
  - [Gert-Jaap Glasbergen: verthash-go](https://github.com/gertjaap/verthash-go/)