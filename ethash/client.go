//go:generate ../.bin/gen-lookup -package ethash -mixBytes 128 -cacheInit 16777216 -cacheGrowth 131072 -datasetInit 1073741824 -datasetGrowth 8388608

package ethash

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

func NewEthereum() *Client {
	var cfg = dag.Config{
		Name:       "ETH",
		Revision:   23,
		StorageDir: common.DefaultDir(".powcache"),

		DatasetInitBytes:   1 << 30,
		DatasetGrowthBytes: 1 << 23,
		CacheInitBytes:     1 << 24,
		CacheGrowthBytes:   1 << 17,

		CacheSizes:   dag.NewLookupTable(cacheSizes, 2048),
		DatasetSizes: dag.NewLookupTable(datasetSizes, 2048),

		MixBytes:        128,
		DatasetParents:  256,
		EpochLength:     30000,
		SeedEpochLength: 30000,

		CacheRounds:    3,
		CachesCount:    3,
		CachesLockMmap: false,

		L1Enabled: false,
	}

	return New(cfg)
}

func NewEthereumClassic() *Client {
	var cfg = dag.Config{
		Name:       "ETC",
		Revision:   23,
		StorageDir: common.DefaultDir(".powcache"),

		DatasetInitBytes:   1 << 30,
		DatasetGrowthBytes: 1 << 23,
		CacheInitBytes:     1 << 24,
		CacheGrowthBytes:   1 << 17,

		CacheSizes:   dag.NewLookupTable(cacheSizes, 2048),
		DatasetSizes: dag.NewLookupTable(datasetSizes, 2048),

		DatasetParents:  256,
		EpochLength:     60000,
		SeedEpochLength: 30000,

		CacheRounds:    3,
		CachesCount:    3,
		CachesLockMmap: false,

		L1Enabled: false,
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
	lookup := c.data.NewLookupFunc512(cache, epoch)

	mix, digest := hashimoto(hash, nonce, size, lookup)
	runtime.KeepAlive(cache)

	return mix, digest, nil
}
