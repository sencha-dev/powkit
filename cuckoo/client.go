package cuckoo

import (
	"encoding/binary"
	"fmt"

	"github.com/sencha-dev/powkit/internal/crypto"
)

type CuckooVariant int

const (
	Cuckoo CuckooVariant = iota
	Cuckatoo
	Cuckaroo
	Cuckarood
	Cuckaroom
	Cuckarooz
)

type Client struct {
	variant   CuckooVariant
	proofSize int
	edgeBits  int
	edgeMask  uint64
	sipnode   crypto.SipNodeFunc
	sipblock  crypto.SipBlockFunc
}

func newClient(variant CuckooVariant, edgeBits, proofSize int, sipnode crypto.SipNodeFunc, sipblock crypto.SipBlockFunc) *Client {
	c := &Client{
		variant:   variant,
		proofSize: proofSize,
		edgeBits:  edgeBits,
		edgeMask:  (uint64(1) << edgeBits) - 1,
		sipnode:   sipnode,
		sipblock:  sipblock,
	}

	return c
}

func NewCuckoo(edgeBits, proofSize int, sipnode crypto.SipNodeFunc, sipblock crypto.SipBlockFunc) *Client {
	return newClient(Cuckoo, edgeBits, proofSize, sipnode, sipblock)
}

func NewAeternity() *Client {
	return NewCuckoo(29, 42, crypto.SipNode24Legacy, nil)
}

func NewCuckaroo(edgeBits, proofSize int, sipnode crypto.SipNodeFunc, sipblock crypto.SipBlockFunc) *Client {
	return newClient(Cuckaroo, edgeBits, proofSize, sipnode, sipblock)
}

func NewCortex() *Client {
	return NewCuckaroo(30, 42, nil, crypto.SipBlock48)
}

func (c *Client) Verify(header []byte, sols []uint64) (bool, error) {
	if len(sols) != c.proofSize {
		return false, fmt.Errorf("sols must be %d uint64s", c.proofSize)
	}

	// create siphash keys
	hash := crypto.Blake2b256(header)
	keys := [4]uint64{
		binary.LittleEndian.Uint64(hash[0:8]),
		binary.LittleEndian.Uint64(hash[8:16]),
		binary.LittleEndian.Uint64(hash[16:24]),
		binary.LittleEndian.Uint64(hash[24:32]),
	}

	switch c.variant {
	case Cuckoo:
		return c.cuckoo(keys, sols)
	case Cuckaroo:
		return c.cuckaroo(keys, sols)
	default:
		return false, fmt.Errorf("unsupported cuckoo variant")
	}
}
