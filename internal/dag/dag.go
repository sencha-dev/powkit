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
	"sync"
	"time"
)

type LightDAG struct {
	mu     sync.Mutex        // Protects the per-epoch map of verification caches
	caches map[uint64]*cache // Currently maintained verification caches
	future *cache            // Pre-generated cache for the estimated future DAG
	cfg    *Config
}

func NewLightDAG(cfg *Config) *LightDAG {
	dag := &LightDAG{
		caches: make(map[uint64]*cache),
		cfg:    cfg,
	}

	return dag
}

func (dag *LightDAG) GetCache(epoch uint64) *cache {
	var c *cache

	dag.mu.Lock()
	if dag.caches == nil {
		dag.caches = make(map[uint64]*cache)
	}

	c = dag.caches[epoch]
	if c == nil {
		// if cache limit is reached, evict the oldest cache entry
		if len(dag.caches) >= dag.cfg.CachesCount {
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
			go dag.future.generate(dag.cfg)
		}
	}

	c.used = time.Now()
	dag.mu.Unlock()

	c.generate(dag.cfg)

	return c
}
