package ethash

import (
	"runtime"

	"github.com/sencha-dev/go-pow/internal/crypto"
	"github.com/sencha-dev/go-pow/internal/dag"
)

type Ethash struct {
	epochLength uint64
	dag         *dag.LightDag
}

const (
	cachesOnDisk = 3
)

func New(name string, epochLength uint64) *Ethash {
	client := &Ethash{
		epochLength: epochLength,
		dag: &dag.LightDag{
			Name:            name,
			EpochLength:     epochLength,
			SeedEpochLength: 30000,
			DatasetParents:  256,
			NumCaches:       cachesOnDisk,
			NeedsL1:         false,
		},
	}

	return client
}

func (e *Ethash) Compute(hash []byte, height, nonce uint64) ([]byte, []byte) {
	epoch := dag.CalcEpoch(height, e.epochLength)
	size := dag.DatasetSize(epoch)
	cache := e.dag.GetCache(epoch)

	keccak512Hasher := crypto.NewKeccak512Hasher()
	lookup := func(index uint32) []uint32 {
		return dag.GenerateDatasetItem512(cache.Cache(), index, keccak512Hasher, e.dag.DatasetParents)
	}

	mix, digest := hashimoto(hash, nonce, size, lookup)
	runtime.KeepAlive(cache)

	return mix, digest
}
