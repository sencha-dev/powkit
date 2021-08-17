package pow

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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

	Chain           string
	Algorithm       string
	NumCaches       int // Maximum number of caches to keep before eviction (only init, don't modify)
	DatasetParents  uint32
	EpochLength     uint64
	SeedEpochLength uint64 // ETC uses 30000 for the seed epoch length but 60000 for the rest
	MinimumHeight   uint64
	NeedsL1         bool
}

func NewLightDag(chain string) (*LightDag, error) {
	chain = strings.ToUpper(chain)
	var dag *LightDag

	switch chain {
	case "ETH":
		dag = &LightDag{
			Chain:           "ETH",
			Algorithm:       "ethash",
			EpochLength:     30000,
			SeedEpochLength: 30000,
			DatasetParents:  256,
			NumCaches:       cachesOnDisk,
			MinimumHeight:   0,
			NeedsL1:         false,
		}
	case "ETC":
		dag = &LightDag{
			Chain:           "ETC",
			Algorithm:       "etchash",
			EpochLength:     60000,
			SeedEpochLength: 30000,
			DatasetParents:  256,
			NumCaches:       cachesOnDisk,
			MinimumHeight:   11700000,
			NeedsL1:         false,
		}
	case "RVN":
		dag = &LightDag{
			Chain:           "RVN",
			Algorithm:       "kawpow",
			EpochLength:     7500,
			SeedEpochLength: 7500,
			DatasetParents:  512,
			NumCaches:       cachesOnDisk,
			MinimumHeight:   1219736,
			NeedsL1:         true,
		}
	default:
		return nil, fmt.Errorf("%s is not supported", chain)
	}

	return dag, nil
}

func (dag *LightDag) getCache(epoch uint64) *cache {
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
			go dag.future.generate(dag.Chain, defaultDir(), dag.NumCaches, cachesLockMmap, dag.NeedsL1)
		}
	}

	c.used = time.Now()
	dag.mu.Unlock()

	c.generate(dag.Chain, defaultDir(), dag.NumCaches, cachesLockMmap, dag.NeedsL1)

	return c
}

func (dag *LightDag) Compute(hash []byte, height, nonce uint64) ([]byte, []byte, error) {
	if height < dag.MinimumHeight {
		return nil, nil, fmt.Errorf("%d is below the minimum height %d for %s", height, dag.MinimumHeight, dag.Chain)
	}

	var mix, digest []byte
	switch dag.Chain {
	case "ETH", "ETC":
		mix, digest = dag.hashimotoLight(height, nonce, hash)
	case "RVN":
		mix, digest = dag.kawpowLight(height, nonce, hash)
	default:
		return nil, nil, fmt.Errorf("unsupported chain %s", dag.Chain)
	}

	return mix, digest, nil
}
