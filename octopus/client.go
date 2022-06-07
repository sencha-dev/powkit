//go:generate ../.bin/gen-lookup -package octopus -mixBytes 256 -cacheInit 16777216 -cacheGrowth 65536 -datasetInit 4294967296 -datasetGrowth 16777216

package octopus

import (
	"runtime"

	"github.com/sencha-dev/powkit/internal/common"
	"github.com/sencha-dev/powkit/internal/dag"
)

type Client struct {
	*dag.DAG
}

func New(cfg dag.Config) *Client {
	client := &Client{
		DAG: dag.New(cfg),
	}

	return client
}

func NewConflux() *Client {
	var cfg = dag.Config{
		Name:       "CFX",
		Revision:   1,
		StorageDir: common.DefaultDir(".powcache"),

		DatasetInitBytes:   2 * (1 << 31),
		DatasetGrowthBytes: 1 << 24,
		CacheInitBytes:     2 * (1 << 23),
		CacheGrowthBytes:   1 << 16,

		CacheSizes:   dag.NewLookupTable(cacheSizes, 2048),
		DatasetSizes: dag.NewLookupTable(datasetSizes, 2048),

		MixBytes:        256,
		DatasetParents:  256,
		EpochLength:     1 << 19,
		SeedEpochLength: 1 << 19,

		CacheRounds:    3,
		CachesCount:    3,
		CachesLockMmap: false,

		L1Enabled: false,
	}

	return New(cfg)
}

func (c *Client) Compute(height, nonce uint64, hash []byte) []byte {
	epoch := c.CalcEpoch(height)
	size := c.DatasetSize(epoch)
	cache := c.GetCache(epoch)
	lookup := c.NewLookupFunc512(cache, epoch)

	digest := octopus(hash, nonce, size, lookup)
	runtime.KeepAlive(cache)

	return digest
}
