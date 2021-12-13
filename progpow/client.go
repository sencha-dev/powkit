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

package progpow

import (
	"runtime"

	"github.com/sencha-dev/powkit/internal/crypto"
	"github.com/sencha-dev/powkit/internal/dag"
)

type Progpow struct {
	dag *dag.LightDAG
	cfg *dag.Config
}

func New(cfg *dag.Config) *Progpow {
	client := &Progpow{
		dag: dag.NewLightDAG(cfg),
		cfg: cfg,
	}

	return client
}

func NewProgpow094() *Progpow {
	return New(dag.Progpow094Cfg)
}

func (e *Progpow) Compute(height, nonce uint64, hash []byte) ([]byte, []byte) {
	epoch := dag.CalcEpoch(e.cfg, height)
	cache := e.dag.GetCache(epoch)

	keccak512Hasher := crypto.NewKeccak512Hasher()
	lookup := func(index uint32) []uint32 {
		return dag.GenerateDatasetItem2048(e.cfg, cache.Cache(), index, keccak512Hasher)
	}

	mix, digest := compute(hash, height, nonce, lookup, cache.L1())
	runtime.KeepAlive(cache)

	return mix, digest
}
