// Copyright 2019 Victor Tran
// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package ethash

import (
	"runtime"

	"github.com/sencha-dev/go-pow/internal/crypto"
	"github.com/sencha-dev/go-pow/internal/dag"
)

type Ethash struct {
	dag *dag.LightDAG
	cfg *dag.Config
}

func New(cfg *dag.Config) *Ethash {
	client := &Ethash{
		dag: dag.NewLightDAG(cfg),
		cfg: cfg,
	}

	return client
}

func NewEthereum() *Ethash {
	return New(dag.EthereumCfg)
}

func NewEthereumClassic() *Ethash {
	return New(dag.EthereumClassicCfg)
}

func (e *Ethash) Compute(height, nonce uint64, hash []byte) ([]byte, []byte) {
	epoch := dag.CalcEpoch(e.cfg, height)
	size := dag.DatasetSize(e.cfg, epoch)
	cache := e.dag.GetCache(epoch)

	keccak512Hasher := crypto.NewKeccak512Hasher()
	lookup := func(index uint32) []uint32 {
		return dag.GenerateDatasetItem512(e.cfg, cache.Cache(), index, keccak512Hasher)
	}

	mix, digest := hashimoto(hash, nonce, size, lookup)
	runtime.KeepAlive(cache)

	return mix, digest
}
