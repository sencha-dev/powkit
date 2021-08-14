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

type PowFunc func(size, height, nonce uint64, cache []uint32, hash []byte) ([]byte, []byte)

type cache struct {
	epoch       uint64
	dump        *os.File
	mmap        mmap.MMap
	cache       []uint32
	once        sync.Once
	used        time.Time
	epochLength uint64
	powFunc     PowFunc
}

// generate ensures that the cache content is generated before use.
func (c *cache) generate(dir string, limit int, lock bool) {
	c.once.Do(func() {
		size := cacheSize(c.epoch)
		seed := seedHash(c.epoch*c.epochLength+1, c.epochLength)

		// If we don't store anything on disk, generate and return.
		if dir == "" {
			c.cache = make([]uint32, size/4)
			generateCache(c.cache, c.epoch, c.epochLength, seed)
			return
		}
		// Disk storage is needed, this will get fancy
		var endian string
		if !isLittleEndian() {
			endian = ".be"
		}
		path := filepath.Join(dir, fmt.Sprintf("cache-R%d-%x%s", algorithmRevision, seed[:8], endian))

		// We're about to mmap the file, ensure that the mapping is cleaned up when the
		// cache becomes unused.
		runtime.SetFinalizer(c, (*cache).finalizer)

		// Try to load the file from disk and memory map it
		var err error
		c.dump, c.mmap, c.cache, err = memoryMap(path, lock)
		if err == nil {
			return
		}

		// No usable previous cache available, create a new cache file to fill
		c.dump, c.mmap, c.cache, err = memoryMapAndGenerate(path, size, lock, func(buffer []uint32) { generateCache(buffer, c.epoch, c.epochLength, seed) })
		if err != nil {
			c.cache = make([]uint32, size/4)
			generateCache(c.cache, c.epoch, c.epochLength, seed)
		}

		// Iterate over all previous instances and delete old ones
		for ep := int(c.epoch) - limit; ep >= 0; ep-- {
			seed := seedHash(uint64(ep)*c.epochLength+1, c.epochLength)
			path := filepath.Join(dir, fmt.Sprintf("cache-R%d-%x%s", algorithmRevision, seed[:8], endian))
			os.Remove(path)
		}
	})
}

// finalizer unmaps the memory and closes the file.
func (c *cache) finalizer() {
	if c.mmap != nil {
		c.mmap.Unmap()
		c.dump.Close()
		c.mmap, c.dump = nil, nil
	}
}

func (c *cache) compute(dagSize, height, nonce uint64, hash []byte) ([]byte, []byte) {
	digest, result := c.powFunc(dagSize, height, nonce, c.cache, hash)
	// Caches are unmapped in a finalizer. Ensure that the cache stays alive
	// until after the call to hashimotoLight so it's not unmapped while being used.
	runtime.KeepAlive(c)
	return digest, result
}

type LightDag struct {
	mu     sync.Mutex        // Protects the per-epoch map of verification caches
	caches map[uint64]*cache // Currently maintained verification caches
	future *cache            // Pre-generated cache for the estimated future DAG

	PowFunc        PowFunc
	Chain          string
	Algorithm      string
	NumCaches      int // Maximum number of caches to keep before eviction (only init, don't modify)
	DatasetParents int
	EpochLength    uint64
	MinimumHeight  uint64
}

func NewLightDag(chain string) (*LightDag, error) {
	chain = strings.ToUpper(chain)
	var dag *LightDag

	switch chain {
	case "ETH":
		dag = &LightDag{
			Chain:          "ETH",
			Algorithm:      "ethash",
			PowFunc:        hashimotoLight,
			EpochLength:    30000,
			DatasetParents: 256,
			NumCaches:      cachesOnDisk,
			MinimumHeight:  0,
		}
	case "ETC":
		dag = &LightDag{
			Chain:          "ETC",
			Algorithm:      "etchash",
			PowFunc:        hashimotoLight,
			EpochLength:    60000,
			DatasetParents: 256,
			NumCaches:      cachesOnDisk,
			MinimumHeight:  11700000,
		}
	case "RVN":
		dag = &LightDag{
			Chain:          "RVN",
			Algorithm:      "kawpow",
			PowFunc:        kawpowLight,
			EpochLength:    7500,
			DatasetParents: 512,
			NumCaches:      cachesOnDisk,
			MinimumHeight:  1219736,
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
			// DEBUG: fmt.Sprintf("evicting dag for epoch %d in favor of %d", evict.epoch, epoch)
			delete(dag.caches, evict.epoch)
		}

		// use the pre generated dag if exists
		if dag.future != nil && dag.future.epoch == epoch {
			// DEBUG: fmt.Sprintf("using pre-generated dag for epoch %d", epoch)
			c, dag.future = dag.future, nil
		} else {
			// DEBUG: fmt.Sprintf("creating new dag for epoch %d", epoch)
			c = &cache{epoch: epoch, epochLength: dag.EpochLength, powFunc: dag.PowFunc}
		}

		dag.caches[epoch] = c
		nextEpoch := epoch + 1
		if dag.future == nil || dag.future.epoch <= epoch {
			// DEBUG: fmt.Sprintf("pre-generating dag for epoch %d", nextEpoch)
			dag.future = &cache{epoch: nextEpoch, epochLength: dag.EpochLength, powFunc: dag.PowFunc}
			go dag.future.generate(defaultDir(), dag.NumCaches, cachesLockMmap)
		}
	}

	c.used = time.Now()
	dag.mu.Unlock()

	c.generate(defaultDir(), dag.NumCaches, cachesLockMmap)

	return c
}

func (dag *LightDag) Compute(hash []byte, height, nonce uint64) ([]byte, []byte, error) {
	if height < dag.MinimumHeight {
		return nil, nil, fmt.Errorf("%d is below the minimum height %d for %s", height, dag.MinimumHeight, dag.Chain)
	}

	epoch := calcEpoch(height, dag.EpochLength)
	dagSize := datasetSize(epoch)

	cache := dag.getCache(epoch)
	mix, digest := cache.compute(uint64(dagSize), height, nonce, hash)

	return mix, digest, nil
}
