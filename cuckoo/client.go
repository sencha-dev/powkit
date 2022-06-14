package cuckoo

import (
	"encoding/base64"
	"encoding/binary"

	"github.com/sencha-dev/powkit/internal/crypto"
)

type Config struct {
	proofSize int
	edgeBits  int
	edgeMask  uint64
	nodeBits  int
	nodeMask  uint64
}

func New(edgeBits, nodeBits, proofSize int) *Config {
	edgeNum := uint64(1) << edgeBits
	nodeNum := uint64(1) << nodeBits

	cfg := &Config{
		proofSize: proofSize,
		edgeBits:  edgeBits,
		edgeMask:  edgeNum - 1,
		nodeBits:  nodeBits,
		nodeMask:  nodeNum - 1,
	}

	return cfg
}

func NewAeternity() *Config {
	return New(29, 29, 42)
}

func (cfg *Config) Verify(hash []byte, nonce uint64, sols []uint64) bool {
	// encode header
	nonceBytes := make([]uint8, 8)
	binary.LittleEndian.PutUint64(nonceBytes, nonce)
	hashEncoded := []byte(base64.StdEncoding.EncodeToString(hash))
	nonceEncoded := []byte(base64.StdEncoding.EncodeToString(nonceBytes))
	header := append(hashEncoded, append(nonceEncoded, make([]byte, 24)...)...)

	// create siphash keys
	h := crypto.Blake2b256(header)
	keys := [4]uint64{
		binary.LittleEndian.Uint64(h[0:8]),
		binary.LittleEndian.Uint64(h[8:16]),
		binary.LittleEndian.Uint64(h[16:24]),
		binary.LittleEndian.Uint64(h[24:32]),
	}

	return cfg.verify(keys, sols)
}
