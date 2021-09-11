// Copyright 2019 Victor Tran
// Copyright 2015 The go-ethereum Authors
// Copyright 2015 Lefteris Karapetsas <lefteris@refu.co>
// Copyright 2015 Matthew Wampler-Doty <matthew.wampler.doty@gmail.com>
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
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/edsrzf/mmap-go"
)

type cache struct {
	epoch           uint64
	epochLength     uint64
	seedEpochLength uint64
	once            sync.Once
	used            time.Time

	cacheDump *os.File
	cacheMmap mmap.MMap
	cache     []uint32
	l1Dump    *os.File
	l1Mmap    mmap.MMap
	l1        []uint32
}

func (c *cache) Cache() []uint32 {
	return c.cache
}

func (c *cache) L1() []uint32 {
	return c.l1
}

// generate ensures that the cache content is generated before use.
func (c *cache) generate(chain, dir string, limit int, lock bool, generateL1 bool) {
	c.once.Do(func() {
		size := cacheSize(c.epoch)
		seed := seedHash(c.epoch*c.epochLength+1, c.seedEpochLength)

		// If we don't store anything on disk, generate and return.
		if dir == "" {
			c.cache = make([]uint32, size/4)
			generateCache(c.cache, c.epoch, c.epochLength, seed)

			if generateL1 {
				c.l1 = make([]uint32, l1CacheNumItems)
				generateL1Cache(c.l1, c.cache)
			}

			return
		}
		// Disk storage is needed, this will get fancy
		var endian string
		if !isLittleEndian() {
			endian = ".be"
		}
		cachePath := filepath.Join(dir, fmt.Sprintf("cache-%s-R%d-%x%s", chain, algorithmRevision, seed[:8], endian))
		l1Path := filepath.Join(dir, fmt.Sprintf("l1-%s-R%d-%x%s", chain, algorithmRevision, seed[:8], endian))

		// We're about to mmap the file, ensure that the mapping is cleaned up when the
		// cache becomes unused.
		runtime.SetFinalizer(c, (*cache).finalizer)

		// Try to load the file from disk and memory map it
		var err error
		c.cacheDump, c.cacheMmap, c.cache, err = memoryMap(cachePath, lock)

		needsCache := err != nil
		needsL1 := generateL1

		if generateL1 {
			c.l1Dump, c.l1Mmap, c.l1, err = memoryMap(l1Path, lock)
			needsL1 = err != nil
		}

		if !needsL1 && !needsCache {
			return
		}

		// No usable previous cache available, create a new cache file to fill
		if needsCache {
			cacheGenerator := func(buffer []uint32) { generateCache(buffer, c.epoch, c.epochLength, seed) }
			c.cacheDump, c.cacheMmap, c.cache, err = memoryMapAndGenerate(cachePath, size, lock, cacheGenerator)
			if err != nil {
				c.cache = make([]uint32, size/4)
				generateCache(c.cache, c.epoch, c.epochLength, seed)
			}
		}

		if needsL1 {
			l1Generator := func(buffer []uint32) { generateL1Cache(buffer, c.cache) }
			c.l1Dump, c.l1Mmap, c.l1, err = memoryMapAndGenerate(l1Path, uint64(l1CacheSize), lock, l1Generator)
			if err != nil {
				c.l1 = make([]uint32, l1CacheNumItems)
				generateL1Cache(c.l1, c.cache)
			}

		}

		// Iterate over all previous instances and delete old ones
		for ep := int(c.epoch) - limit; ep >= 0; ep-- {
			seed := seedHash(uint64(ep)*c.epochLength+1, c.seedEpochLength)
			cachePath := filepath.Join(dir, fmt.Sprintf("cache-%s-R%d-%x%s", chain, algorithmRevision, seed[:8], endian))
			l1Path := filepath.Join(dir, fmt.Sprintf("l1-%s-R%d-%x%s", chain, algorithmRevision, seed[:8], endian))
			os.Remove(cachePath)
			os.Remove(l1Path)
		}
	})
}

// finalizer unmaps the memory and closes the file.
func (c *cache) finalizer() {
	if c.cacheMmap != nil {
		c.cacheMmap.Unmap()
		c.cacheDump.Close()
		c.cacheMmap, c.cacheDump = nil, nil
	}
}

type LightDag struct {
	mu     sync.Mutex        // Protects the per-epoch map of verification caches
	caches map[uint64]*cache // Currently maintained verification caches
	future *cache            // Pre-generated cache for the estimated future DAG

	Name            string
	NumCaches       int // Maximum number of caches to keep before eviction (only init, don't modify)
	DatasetParents  uint32
	EpochLength     uint64
	SeedEpochLength uint64 // ETC uses 30000 for the seed epoch length but 60000 for the rest
	NeedsL1         bool
}

func (dag *LightDag) GetCache(epoch uint64) *cache {
	var c *cache

	dag.mu.Lock()
	if dag.caches == nil {
		dag.caches = make(map[uint64]*cache)
	}

	c = dag.caches[epoch]
	if c == nil {
		// if cache limit is reached, evict the oldest cache entry
		if len(dag.caches) >= dag.NumCaches {
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
			c = &cache{epoch: epoch, epochLength: dag.EpochLength, seedEpochLength: dag.SeedEpochLength}
		}

		dag.caches[epoch] = c
		nextEpoch := epoch + 1
		if dag.future == nil || dag.future.epoch <= epoch {
			dag.future = &cache{epoch: nextEpoch, epochLength: dag.EpochLength, seedEpochLength: dag.SeedEpochLength}
			go dag.future.generate(dag.Name, defaultDir(), dag.NumCaches, cachesLockMmap, dag.NeedsL1)
		}
	}

	c.used = time.Now()
	dag.mu.Unlock()

	c.generate(dag.Name, defaultDir(), dag.NumCaches, cachesLockMmap, dag.NeedsL1)

	return c
}
