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
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

func cachePaths(cfg *Config, seed []byte) (string, string) {
	cacheName := fmt.Sprintf("cache-%s-R%d-%x", cfg.Name, cfg.Revision, seed)
	l1Name := fmt.Sprintf("l1-%s-R%d-%x", cfg.Name, cfg.Revision, seed)

	cachePath := filepath.Join(cfg.StorageDir, cacheName)
	l1Path := filepath.Join(cfg.StorageDir, l1Name)

	return cachePath, l1Path
}

type cache struct {
	epoch uint64
	once  sync.Once
	used  time.Time
	cache dataFile
	l1    dataFile
}

func (c *cache) Cache() []uint32 {
	return c.cache.data
}

func (c *cache) L1() []uint32 {
	return c.l1.data
}

// generate ensures that the cache content is generated before use.
func (c *cache) generate(cfg *Config) {
	c.once.Do(func() {
		size := CacheSize(cfg, c.epoch)
		seed := SeedHash(cfg, c.epoch*cfg.EpochLength+1)

		// If we don't store anything on disk, generate and return.
		if cfg.StorageDir == "" {
			c.cache.data = make([]uint32, size/4)
			generateCache(cfg, c.cache.data, c.epoch, seed)

			if cfg.L1Enabled {
				c.l1.data = make([]uint32, cfg.L1CacheNumItems)
				generateL1Cache(cfg, c.l1.data, c.cache.data)
			}

			return
		}

		cachePath, l1Path := cachePaths(cfg, seed[:8])

		// We're about to mmap the file, ensure that the mapping is cleaned up when the
		// cache becomes unused.
		runtime.SetFinalizer(c, (*cache).finalizer)

		// Try to load the file from disk and memory map it
		var err error
		c.cache, err = memoryMap(cachePath, cfg.CachesLockMmap)
		needsCache := err != nil

		needsL1 := cfg.L1Enabled
		if cfg.L1Enabled {
			c.l1, err = memoryMap(l1Path, cfg.CachesLockMmap)
			needsL1 = err != nil
		}

		if !needsL1 && !needsCache {
			return
		}

		// No usable previous cache available, create a new cache file to fill
		if needsCache {
			cacheGenerator := func(buffer []uint32) { generateCache(cfg, buffer, c.epoch, seed) }
			c.cache, err = memoryMapAndGenerate(cachePath, size, cfg.CachesLockMmap, cacheGenerator)
			if err != nil {
				c.cache.data = make([]uint32, size/4)
				generateCache(cfg, c.cache.data, c.epoch, seed)
			}
		}

		if needsL1 {
			l1Generator := func(buffer []uint32) { generateL1Cache(cfg, buffer, c.cache.data) }
			c.l1, err = memoryMapAndGenerate(l1Path, cfg.L1CacheSize, cfg.CachesLockMmap, l1Generator)
			if err != nil {
				c.l1.data = make([]uint32, cfg.L1CacheNumItems)
				generateL1Cache(cfg, c.l1.data, c.cache.data)
			}
		}

		// Iterate over all previous instances and delete old ones
		for ep := int(c.epoch) - cfg.CachesCount; ep >= 0; ep-- {
			seed := SeedHash(cfg, uint64(ep)*cfg.EpochLength+1)
			cachePath, l1Path := cachePaths(cfg, seed[:8])
			os.Remove(cachePath)
			os.Remove(l1Path)
		}
	})
}

// finalizer unmaps the memory and closes the file.
func (c *cache) finalizer() {
	if c.cache.mmap != nil {
		c.cache.mmap.Unmap()
		c.cache.dump.Close()
		c.cache.mmap, c.cache.dump = nil, nil
	}

	if c.l1.mmap != nil {
		c.l1.mmap.Unmap()
		c.l1.dump.Close()
		c.l1.mmap, c.l1.dump = nil, nil
	}
}
