package progpow

type Config struct {
	EpochLength         uint64
	PeriodLength        uint64
	DagLoads            int
	CacheBytes          uint32
	RoundCount          int
	RoundCacheAccesses  int
	RoundMathOperations int
}

var (
	ProgPow092 = &Config{
		EpochLength:         30000,
		PeriodLength:        50,
		DagLoads:            4,
		CacheBytes:          16 * 1024,
		RoundCount:          64,
		RoundCacheAccesses:  12,
		RoundMathOperations: 20,
	}

	ProgPow093 = &Config{
		EpochLength:         30000,
		PeriodLength:        10,
		DagLoads:            4,
		CacheBytes:          16 * 1024,
		RoundCount:          64,
		RoundCacheAccesses:  11,
		RoundMathOperations: 18,
	}

	ProgPow094 = &Config{
		EpochLength:         30000,
		PeriodLength:        10,
		DagLoads:            4,
		CacheBytes:          16 * 1024,
		RoundCount:          64,
		RoundCacheAccesses:  11,
		RoundMathOperations: 18,
	}

	Kawpow = &Config{
		EpochLength:         7500,
		PeriodLength:        3,
		DagLoads:            4,
		CacheBytes:          16 * 1024,
		RoundCount:          64,
		RoundCacheAccesses:  11,
		RoundMathOperations: 18,
	}

	Firopow = &Config{
		EpochLength:         1300,
		PeriodLength:        1,
		DagLoads:            4,
		CacheBytes:          16 * 1024,
		RoundCount:          64,
		RoundCacheAccesses:  11,
		RoundMathOperations: 18,
	}
)
