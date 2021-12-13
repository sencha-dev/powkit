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
	"math/big"

	"github.com/sencha-dev/powkit/internal/crypto"
)

func CalcEpoch(cfg *Config, height uint64) uint64 {
	epoch := height / cfg.EpochLength

	return epoch
}

func SeedHash(cfg *Config, height uint64) []byte {
	seed := make([]byte, 32)
	if height < cfg.SeedEpochLength {
		return seed
	}

	keccak256Hasher := crypto.NewKeccak256Hasher()
	for i := 0; i < int(height/cfg.SeedEpochLength); i++ {
		keccak256Hasher(seed, seed)
	}

	return seed
}

func CacheSize(cfg *Config, epoch uint64) uint64 {
	if cfg.CacheSizes != nil && epoch < cfg.CacheSizes.maxEpoch {
		return cfg.CacheSizes.table[epoch]
	}

	return calcCacheSize(cfg, epoch)
}

func DatasetSize(cfg *Config, epoch uint64) uint64 {
	if cfg.DatasetSizes != nil && epoch < cfg.DatasetSizes.maxEpoch {
		return cfg.DatasetSizes.table[epoch]
	}

	return calcDatasetSize(cfg, epoch)
}

func calcCacheSize(cfg *Config, epoch uint64) uint64 {
	size := cfg.CacheInitBytes + cfg.CacheGrowthBytes*epoch - hashBytes

	// Always accurate for n < 2^64
	for !new(big.Int).SetUint64(size / hashBytes).ProbablyPrime(1) {
		size -= 2 * hashBytes
	}

	return size
}

func calcDatasetSize(cfg *Config, epoch uint64) uint64 {
	size := cfg.DatasetInitBytes + cfg.DatasetGrowthBytes*epoch - mixBytes

	// Always accurate for n < 2^64
	for !new(big.Int).SetUint64(size / mixBytes).ProbablyPrime(1) {
		size -= 2 * mixBytes
	}

	return size
}
