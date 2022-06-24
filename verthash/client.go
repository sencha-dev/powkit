package verthash

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sencha-dev/powkit/internal/common"
)

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

type Client struct {
	data []byte
}

func New() (*Client, error) {
	path := filepath.Join(common.DefaultDir(".powcache"), "verthash.dat")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		writeGraph(path)
	}

	err := validateGraph(path)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	c := &Client{
		data: data,
	}

	return c, nil
}

func (c *Client) Compute(input []byte) []byte {
	return verthash(c.data, input)
}
