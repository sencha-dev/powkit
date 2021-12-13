package progpow

type Config struct {
	PeriodLength        uint64
	DagLoads            int
	CacheBytes          uint32
	RoundCount          int
	RoundCacheAccesses  int
	RoundMathOperations int
}

var (
	ProgPow092 = &Config{
		PeriodLength:        50,
		DagLoads:            4,
		CacheBytes:          16 * 1024,
		RoundCount:          64,
		RoundCacheAccesses:  12,
		RoundMathOperations: 20,
	}

	ProgPow093 = &Config{
		PeriodLength:        10,
		DagLoads:            4,
		CacheBytes:          16 * 1024,
		RoundCount:          64,
		RoundCacheAccesses:  11,
		RoundMathOperations: 18,
	}

	ProgPow094 = &Config{
		PeriodLength:        10,
		DagLoads:            4,
		CacheBytes:          16 * 1024,
		RoundCount:          64,
		RoundCacheAccesses:  11,
		RoundMathOperations: 18,
	}

	Kawpow = &Config{
		PeriodLength:        3,
		DagLoads:            4,
		CacheBytes:          16 * 1024,
		RoundCount:          64,
		RoundCacheAccesses:  11,
		RoundMathOperations: 18,
	}

	Firopow = &Config{
		PeriodLength:        1,
		DagLoads:            4,
		CacheBytes:          16 * 1024,
		RoundCount:          64,
		RoundCacheAccesses:  11,
		RoundMathOperations: 18,
	}
)
