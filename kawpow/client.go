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

//go:generate ../.bin/gen-lookup -package kawpow -mixBytes 128 -cacheInit 16777216 -cacheGrowth 131072 -datasetInit 1073741824 -datasetGrowth 8388608

package kawpow

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

func (c *Client) Compute(height, nonce uint64, hash []byte) ([]byte, []byte) {
	epoch := c.CalcEpoch(height)
	size := c.DatasetSize(epoch)
	cache := c.GetCache(epoch)
	lookup := c.NewLookupFunc2048(cache, epoch)

	mix, digest := kawpow(hash, height, nonce, size, lookup, cache.L1())
	runtime.KeepAlive(cache)

	return mix, digest
}
