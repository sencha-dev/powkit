package progpow

import (
	"github.com/sencha-dev/powkit/internal/dag"
)

type Config struct {
	DagCfg              *dag.Config
	PeriodLength        uint64
	DagLoads            int
	CacheBytes          uint32
	RoundCount          int
	RoundCacheAccesses  int
	RoundMathOperations int
}

var (
	ProgPow092 = &Config{
		DagCfg:              dag.Progpow092Cfg,
		PeriodLength:        50,
		DagLoads:            4,
		CacheBytes:          16 * 1024,
		RoundCount:          64,
		RoundCacheAccesses:  12,
		RoundMathOperations: 20,
	}

	ProgPow093 = &Config{
		DagCfg:              dag.Progpow093Cfg,
		PeriodLength:        10,
		DagLoads:            4,
		CacheBytes:          16 * 1024,
		RoundCount:          64,
		RoundCacheAccesses:  11,
		RoundMathOperations: 18,
	}

	ProgPow094 = &Config{
		DagCfg:              dag.Progpow094Cfg,
		PeriodLength:        10,
		DagLoads:            4,
		CacheBytes:          16 * 1024,
		RoundCount:          64,
		RoundCacheAccesses:  11,
		RoundMathOperations: 18,
	}

	Kawpow = &Config{
		DagCfg:              dag.RavencoinCfg,
		PeriodLength:        3,
		DagLoads:            4,
		CacheBytes:          16 * 1024,
		RoundCount:          64,
		RoundCacheAccesses:  11,
		RoundMathOperations: 18,
	}

	Firopow = &Config{
		DagCfg:              dag.FiroCfg,
		PeriodLength:        1,
		DagLoads:            4,
		CacheBytes:          16 * 1024,
		RoundCount:          64,
		RoundCacheAccesses:  11,
		RoundMathOperations: 18,
	}
)
