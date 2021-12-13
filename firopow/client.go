package firopow

import (
	"runtime"

	"github.com/sencha-dev/go-pow/internal/crypto"
	"github.com/sencha-dev/go-pow/internal/dag"
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
	return New(dag.FiroCfg)
}

func (e *Firopow) Compute(height, nonce uint64, hash []byte) ([]byte, []byte) {
	epoch := dag.CalcEpoch(e.cfg, height)
	cache := e.dag.GetCache(epoch)

	keccak512Hasher := crypto.NewKeccak512Hasher()
	lookup := func(index uint32) []uint32 {
		return dag.GenerateDatasetItem2048(e.cfg, cache.Cache(), index, keccak512Hasher)
	}

	mix, digest := firopow(hash, height, nonce, lookup, cache.L1())
	runtime.KeepAlive(cache)

	return mix, digest
}
