package verthash

const (
	// graph constants
	hashSize    = 32
	nodeSize    = hashSize
	datasetSize = 1283457024

	// verthash constants
	verthashHeaderSize    uint32 = 80
	verthashHashOutSize   uint32 = 32
	verthashP0Size        uint32 = 64
	verthashIter          uint32 = 8
	verthashSubset        uint32 = verthashP0Size * verthashIter
	verthashRotations     uint32 = 32
	verthashIndexes       uint32 = 4096
	verthashByteAlignment uint32 = 16
)
