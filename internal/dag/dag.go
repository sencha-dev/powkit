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

package dag

import (
	"fmt"
	"math/big"
	"path/filepath"
	"sync"
	"time"

	"github.com/sencha-dev/powkit/internal/crypto"
)

type DAG struct {
	Config
	mu     sync.Mutex        // Protects the per-epoch map of verification caches
	caches map[uint64]*cache // Currently maintained verification caches
	future *cache            // Pre-generated cache for the estimated future DAG
}

func New(cfg Config) *DAG {
	dag := &DAG{
		Config: cfg,
		caches: make(map[uint64]*cache),
	}

	return dag
}

/* helpers */

func (d *DAG) cacheStorageLocation(seed []byte) string {
	name := fmt.Sprintf("cache-%s-R%d-%x", d.Name, d.Revision, seed)
	path := filepath.Join(d.StorageDir, name)

	return path
}

func (d *DAG) l1StorageLocation(seed []byte) string {
	name := fmt.Sprintf("l1-%s-R%d-%x", d.Name, d.Revision, seed)
	path := filepath.Join(d.StorageDir, name)

	return path
}

/* calculations */

func (d *DAG) CalcEpoch(height uint64) uint64 {
	epoch := height / d.EpochLength

	return epoch
}

func (d *DAG) SeedHash(height uint64) []byte {
	seed := make([]byte, 32)
	if height < d.SeedEpochLength {
		return seed
	}

	keccak256Hasher := crypto.NewKeccak256Hasher()
	for i := 0; i < int(height/d.SeedEpochLength); i++ {
		keccak256Hasher(seed, seed)
	}

	return seed
}

func (d *DAG) CacheSize(epoch uint64) uint64 {
	if d.CacheSizes != nil && epoch < d.CacheSizes.maxEpoch {
		return d.CacheSizes.table[epoch]
	}

	return d.calcCacheSize(epoch)
}

func (d *DAG) DatasetSize(epoch uint64) uint64 {
	if d.DatasetSizes != nil && epoch < d.DatasetSizes.maxEpoch {
		return d.DatasetSizes.table[epoch]
	}

	return d.calcDatasetSize(epoch)
}

func (d *DAG) calcCacheSize(epoch uint64) uint64 {
	size := d.CacheInitBytes + d.CacheGrowthBytes*epoch - hashBytes

	// Always accurate for n < 2^64
	for !new(big.Int).SetUint64(size / hashBytes).ProbablyPrime(1) {
		size -= 2 * hashBytes
	}

	return size
}

func (d *DAG) calcDatasetSize(epoch uint64) uint64 {
	size := d.DatasetInitBytes + d.DatasetGrowthBytes*epoch - d.MixBytes

	// Always accurate for n < 2^64
	for !new(big.Int).SetUint64(size / d.MixBytes).ProbablyPrime(1) {
		size -= 2 * d.MixBytes
	}

	return size
}

/* cache */

func (dag *DAG) GetCache(epoch uint64) *cache {
	var c *cache

	dag.mu.Lock()
	if dag.caches == nil {
		dag.caches = make(map[uint64]*cache)
	}

	c = dag.caches[epoch]
	if c == nil {
		// if cache limit is reached, evict the oldest cache entry
		if len(dag.caches) >= dag.CachesCount {
			var evict *cache
			for _, cache := range dag.caches {
				if evict == nil || evict.used.After(cache.used) {
					evict = cache
				}
			}
			delete(dag.caches, evict.epoch)
		}

		// use the pre generated dag if exists
		if dag.future != nil && dag.future.epoch == epoch {
			c, dag.future = dag.future, nil
		} else {
			c = &cache{epoch: epoch}
		}

		dag.caches[epoch] = c
		nextEpoch := epoch + 1
		if dag.future == nil || dag.future.epoch <= epoch {
			dag.future = &cache{epoch: nextEpoch}
			go dag.future.generate(dag)
		}
	}

	c.used = time.Now()
	dag.mu.Unlock()

	c.generate(dag)

	return c
}

/* lookups */

type LookupFunc func(index uint32) []uint32

func (dag *DAG) NewLookupFunc512(c *cache, epoch uint64) LookupFunc {
	keccak512Hasher := crypto.NewKeccak512Hasher()
	lookup := func(index uint32) []uint32 {
		return dag.generateDatasetItemUint(c.Cache(), index, 1, keccak512Hasher)
	}

	return lookup
}

func (dag *DAG) NewLookupFunc1024(c *cache, epoch uint64) LookupFunc {
	keccak512Hasher := crypto.NewKeccak512Hasher()
	lookup := func(index uint32) []uint32 {
		return dag.generateDatasetItemUint(c.Cache(), index, 2, keccak512Hasher)
	}

	return lookup
}

func (dag *DAG) NewLookupFunc2048(c *cache, epoch uint64) LookupFunc {
	keccak512Hasher := crypto.NewKeccak512Hasher()
	lookup := func(index uint32) []uint32 {
		return dag.generateDatasetItemUint(c.Cache(), index, 4, keccak512Hasher)
	}

	return lookup
}
