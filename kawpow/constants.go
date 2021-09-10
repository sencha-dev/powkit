package kawpow

const (
	fnvOffsetBasis uint32 = 0x811c9dc5

	// kawpow constants
	periodLength      uint64 = 3
	numRegs           uint32 = 32
	numLanes          uint32 = 16
	numCacheAccesses  int    = 11
	numMathOperations int    = 18
	dagLoads          int    = 4
	kawpowRounds      int    = 64
	l1CacheSize       uint32 = 4096 * 4
	l1CacheNumItems   uint32 = l1CacheSize / 4
)

var ravencoinKawpow [15]uint32 = [15]uint32{
	0x00000072, //R
	0x00000041, //A
	0x00000056, //V
	0x00000045, //E
	0x0000004E, //N
	0x00000043, //C
	0x0000004F, //O
	0x00000049, //I
	0x0000004E, //N
	0x0000004B, //K
	0x00000041, //A
	0x00000057, //W
	0x00000050, //P
	0x0000004F, //O
	0x00000057, //W
}
