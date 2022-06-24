//go:generate ../.bin/gen-lookup -package kawpow -mixBytes 128 -cacheInit 16777216 -cacheGrowth 131072 -datasetInit 1073741824 -datasetGrowth 8388608

package kawpow

import (
	"fmt"
	"runtime"

	"github.com/sencha-dev/powkit/internal/common"
	"github.com/sencha-dev/powkit/internal/dag"
)

type Client struct {
	data *dag.DAG
}

func New(cfg dag.Config) *Client {
	client := &Client{
		data: dag.New(cfg),
	}

	return client
}

func NewRavencoin() *Client {
	var cfg = dag.Config{
		Name:       "RVN",
		Revision:   23,
		StorageDir: common.DefaultDir(".powcache"),

		DatasetInitBytes:   1 << 30,
		DatasetGrowthBytes: 1 << 23,
		CacheInitBytes:     1 << 24,
		CacheGrowthBytes:   1 << 17,

		CacheSizes:   dag.NewLookupTable(cacheSizes, 2048),
		DatasetSizes: dag.NewLookupTable(datasetSizes, 2048),

		MixBytes:        128,
		DatasetParents:  512,
		EpochLength:     7500,
		SeedEpochLength: 7500,

		CacheRounds:    3,
		CachesCount:    3,
		CachesLockMmap: false,

		L1Enabled:       true,
		L1CacheSize:     4096 * 4,
		L1CacheNumItems: 4096,
	}

	return New(cfg)
}

func (c *Client) Compute(hash []byte, height, nonce uint64) ([]byte, []byte, error) {
	if len(hash) != 32 {
		return nil, nil, fmt.Errorf("hash must be 32 bytes")
	}

	epoch := c.data.CalcEpoch(height)
	size := c.data.DatasetSize(epoch)
	cache := c.data.GetCache(epoch)
	lookup := c.data.NewLookupFunc2048(cache, epoch)

	mix, digest := kawpow(hash, height, nonce, size, lookup, cache.L1())
	runtime.KeepAlive(cache)

	return mix, digest, nil
}
