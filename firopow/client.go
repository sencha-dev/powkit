//go:generate ../.bin/gen-lookup -package firopow -cacheInit 16777216 -cacheGrowth 131072 -datasetInit 1610612736 -datasetGrowth 8388608

package firopow

import (
	"runtime"

	"github.com/sencha-dev/powkit/internal/common"
	"github.com/sencha-dev/powkit/internal/crypto"
	"github.com/sencha-dev/powkit/internal/dag"
)

type Firopow struct {
	dag *dag.LightDAG
	cfg *dag.Config
}

func New(cfg *dag.Config) *Firopow {
	client := &Firopow{
		dag: dag.NewLightDAG(cfg),
		cfg: cfg,
	}

	return client
}

func NewFiro() *Firopow {
	var cfg = &dag.Config{
		Name:       "FIRO",
		Revision:   23,
		StorageDir: common.DefaultDir(".powcache"),

		DatasetInitBytes:   (1 << 30) + (1 << 29),
		DatasetGrowthBytes: 1 << 23,
		CacheInitBytes:     1 << 24,
		CacheGrowthBytes:   1 << 17,

		CacheSizes:   dag.NewLookupTable(cacheSizes, 2048),
		DatasetSizes: dag.NewLookupTable(datasetSizes, 2048),

		DatasetParents:  512,
		EpochLength:     1300,
		SeedEpochLength: 1300,

		CacheRounds:    3,
		CachesCount:    3,
		CachesLockMmap: false,

		L1Enabled:       true,
		L1CacheSize:     4096 * 4,
		L1CacheNumItems: 4096,
	}

	return New(cfg)
}

func (e *Firopow) Compute(height, nonce uint64, hash []byte) ([]byte, []byte) {
	epoch := dag.CalcEpoch(e.cfg, height)
	datasetSize := dag.DatasetSize(e.cfg, epoch)
	cache := e.dag.GetCache(epoch)

	keccak512Hasher := crypto.NewKeccak512Hasher()
	lookup := func(index uint32) []uint32 {
		return dag.GenerateDatasetItem2048(e.cfg, cache.Cache(), index, keccak512Hasher)
	}

	mix, digest := compute(hash, height, nonce, datasetSize, lookup, cache.L1())
	runtime.KeepAlive(cache)

	return mix, digest
}
