package kawpow

import (
	"runtime"

	"github.com/sencha-dev/go-pow/internal/crypto"
	"github.com/sencha-dev/go-pow/internal/dag"
)

type Kawpow struct {
	epochLength uint64
	dag         *dag.LightDag
}

const (
	cachesOnDisk = 3
)

func New(name string, epochLength uint64) *Kawpow {
	client := &Kawpow{
		epochLength: epochLength,
		dag: &dag.LightDag{
			Name:            name,
			EpochLength:     epochLength,
			SeedEpochLength: epochLength,
			DatasetParents:  512,
			NumCaches:       cachesOnDisk,
			NeedsL1:         true,
		},
	}

	return client
}

func (k *Kawpow) Compute(height, nonce uint64, hash []byte) ([]byte, []byte) {
	epoch := dag.CalcEpoch(height, k.epochLength)
	cache := k.dag.GetCache(epoch)

	keccak512Hasher := crypto.NewKeccak512Hasher()
	lookup := func(index uint32) []uint32 {
		return dag.GenerateDatasetItem2048(cache.Cache(), index, keccak512Hasher, k.dag.DatasetParents)
	}

	mix, digest := kawpow(cache.L1(), hash, height, nonce, lookup)
	runtime.KeepAlive(cache)

	return mix, digest
}
