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
	"reflect"
	"testing"

	"github.com/sencha-dev/powkit/internal/common/testutil"
)

func TestEpochNumber(t *testing.T) {
	tests := []struct {
		cfg    *Config
		height uint64
		epoch  uint64
	}{
		{
			cfg:    RavencoinCfg,
			height: 0,
			epoch:  0,
		},
		{
			cfg:    RavencoinCfg,
			height: 1,
			epoch:  0,
		},
		{
			cfg:    RavencoinCfg,
			height: 7499,
			epoch:  0,
		},
		{
			cfg:    RavencoinCfg,
			height: 7500,
			epoch:  1,
		},
		{
			cfg:    RavencoinCfg,
			height: 7501,
			epoch:  1,
		},
		{
			cfg:    RavencoinCfg,
			height: 7502,
			epoch:  1,
		},
		{
			cfg:    RavencoinCfg,
			height: 1245000,
			epoch:  166,
		},
	}

	for i, tt := range tests {
		epoch := CalcEpoch(tt.cfg, tt.height)

		if epoch != tt.epoch {
			t.Errorf("failed on %d: epoch mismatch: have %d want %d", i, epoch, tt.epoch)
		}
	}
}

func TestSeedHash(t *testing.T) {
	tests := []struct {
		cfg   *Config
		epoch uint64
		seed  []byte
	}{
		{
			cfg:   RavencoinCfg,
			epoch: 0,
			seed:  testutil.MustDecodeHex("0x0000000000000000000000000000000000000000000000000000000000000000"),
		},
		{
			cfg:   RavencoinCfg,
			epoch: 1,
			seed:  testutil.MustDecodeHex("0x290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563"),
		},
		{
			cfg:   RavencoinCfg,
			epoch: 171,
			seed:  testutil.MustDecodeHex("0xa9b0e0c9aca72c07ba06b5bbdae8b8f69e61878301508473379bb4f71807d707"),
		},
		{
			cfg:   RavencoinCfg,
			epoch: 2048,
			seed:  testutil.MustDecodeHex("0x20a7678ca7b50829183baac2e1e3c43fa3c4bcbc171b11cf5a9f30bebd172920"),
		},
		{
			cfg:   RavencoinCfg,
			epoch: 29998,
			seed:  testutil.MustDecodeHex("0x1222b1faed7f93098f8ae498621fb3479805a664b70186063861c46596c66164"),
		},
		{
			cfg:   RavencoinCfg,
			epoch: 29999,
			seed:  testutil.MustDecodeHex("0xee1d0f61b054dff0f3025ebba821d405c8dc19a983e582e9fa5436fc3e7a07d8"),
		},
	}

	for i, tt := range tests {
		seed := SeedHash(tt.cfg, uint64(tt.epoch)*tt.cfg.EpochLength+1)
		if !reflect.DeepEqual(seed, tt.seed) {
			t.Errorf("failed on %d: seed mismatch: have %x, want %x", i, seed, tt.seed)
		}
	}
}

func TestCalcCacheSize(t *testing.T) {
	tests := []struct {
		cfg   *Config
		epoch uint64
		size  uint64
	}{
		{
			cfg:   EthereumCfg,
			epoch: 0,
			size:  16776896,
		},
		{
			cfg:   EthereumCfg,
			epoch: 1,
			size:  16907456,
		},
		{
			cfg:   EthereumCfg,
			epoch: 2,
			size:  17039296,
		},
		{
			cfg:   EthereumCfg,
			epoch: 3,
			size:  17170112,
		},
		{
			cfg:   EthereumCfg,
			epoch: 4,
			size:  17301056,
		},
		{
			cfg:   EthereumCfg,
			epoch: 5,
			size:  17432512,
		},
		{
			cfg:   EthereumCfg,
			epoch: 6,
			size:  17563072,
		},
		{
			cfg:   EthereumCfg,
			epoch: 7,
			size:  17693888,
		},
		{
			cfg:   EthereumCfg,
			epoch: 8,
			size:  17824192,
		},
		{
			cfg:   EthereumCfg,
			epoch: 9,
			size:  17955904,
		},
	}

	for i, tt := range tests {
		size := calcCacheSize(tt.cfg, tt.epoch)
		if size != tt.size {
			t.Errorf("failed on %d: size mismatch: have %d, want %d", i, size, tt.size)
		}
	}
}

func TestCalcDatasetSize(t *testing.T) {
	tests := []struct {
		cfg   *Config
		epoch uint64
		size  uint64
	}{
		{
			cfg:   EthereumCfg,
			epoch: 0,
			size:  1073739904,
		},
		{
			cfg:   EthereumCfg,
			epoch: 1,
			size:  1082130304,
		},
		{
			cfg:   EthereumCfg,
			epoch: 2,
			size:  1090514816,
		},
		{
			cfg:   EthereumCfg,
			epoch: 3,
			size:  1098906752,
		},
		{
			cfg:   EthereumCfg,
			epoch: 4,
			size:  1107293056,
		},
		{
			cfg:   EthereumCfg,
			epoch: 5,
			size:  1115684224,
		},
		{
			cfg:   EthereumCfg,
			epoch: 6,
			size:  1124070016,
		},
		{
			cfg:   EthereumCfg,
			epoch: 7,
			size:  1132461952,
		},
		{
			cfg:   EthereumCfg,
			epoch: 8,
			size:  1140849536,
		},
		{
			cfg:   EthereumCfg,
			epoch: 9,
			size:  1149232768,
		},
	}

	for i, tt := range tests {
		size := calcDatasetSize(tt.cfg, tt.epoch)
		if size != tt.size {
			t.Errorf("failed on %d: size mismatch: have %d, want %d", i, size, tt.size)
		}
	}
}
