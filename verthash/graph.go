package verthash

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"golang.org/x/crypto/sha3"
)

func log2(x int64) int64 {
	r := int64(0)
	for ; x > 1; x >>= 1 {
		r++
	}

	return r
}

type node struct {
	H []byte // hash at the file
}

func (n *node) MarshalBinary() ([]byte, error) {
	return n.H, nil
}

func (n *node) UnmarshalBinary(data []byte) error {
	n.H = data
	return nil
}

type graph struct {
	pk    []byte
	fn    string
	db    *os.File
	index int64
	log2  int64
	pow2  int64
	size  int64
}

func writeGraph(path string) {
	pk := sha3.Sum256([]byte("Verthash Proof-of-Space Datafile"))
	newGraph(17, path, pk[:], nil)
}

func validateGraph(path string) error {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	if size := len(b); size != datasetSize {
		return fmt.Errorf("verthash graph size mismatch: have %d want %d", size, datasetSize)
	}

	hash := sha256.Sum256(b)
	expectedHash, _ := hex.DecodeString("a55531e843cd56b010114aaf6325b0d529ecf88f8ad47639b6ededafd721aa48")
	if !bytes.Equal(hash[:], expectedHash) {
		return fmt.Errorf("verthash graph checksum hash mismatch: have %x want %x", hash, expectedHash)
	}

	return nil
}

// Generate a new PoS graph of index
// Currently only supports the weaker PoS graph
// Note that this graph will have O(2^index) nodes
func newGraph(index int64, fn string, pk []byte, progress chan float64) *graph {
	var db *os.File
	_, err := os.Stat(fn)
	fileExists := err == nil
	if fileExists { //file exists
		db, err = os.OpenFile(fn, os.O_RDWR, 0666)
		if err != nil {
			panic(err)
		}
	} else {
		db, err = os.Create(fn)
		if err != nil {
			panic(err)
		}
	}

	size := numXi(index)
	log2 := log2(size) + 1
	pow2 := int64(1 << log2)

	g := &graph{
		pk:    pk,
		fn:    fn,
		db:    db,
		index: index,
		log2:  log2,
		size:  size,
		pow2:  pow2,
	}

	done := make(chan bool, 1)
	if !fileExists {
		go func() {
			g.xiGraphIter(index)
			done <- true
		}()

		for {
			genDone := false
			select {
			case genDone = <-done:
			case <-time.After(1 * time.Second):
			}

			if genDone {
				break
			}
			if progress != nil {
				s, err := g.db.Stat()
				if err == nil {
					progress <- float64(s.Size()) / float64(1283457024)
				}
			}
		}

	}

	return g
}

// compute parents of nodes
func (g *graph) getParents(node, index int64) []int64 {
	if node < int64(1<<uint64(index)) {
		return nil
	}

	offset0, offset1 := g.getGraph(node, index)

	var res []int64
	if offset0 != 0 {
		res = append(res, node-offset0)
	}
	if offset1 != 0 {
		res = append(res, node-offset1)
	}
	return res
}

// compute the offsets for the two parents in the butterfly graph
func (g *graph) butterflyParents(begin, node, index int64) (int64, int64) {
	pow2index_1 := int64(1 << uint64(index-1))
	level := (node - begin) / pow2index_1
	var prev int64
	shift := (index - 1) - level
	if level > (index - 1) {
		shift = level - (index - 1)
	}
	i := (node - begin) % pow2index_1
	if (i>>uint64(shift))&1 == 0 {
		prev = i + (1 << uint64(shift))
	} else {
		prev = i - (1 << uint64(shift))
	}
	parent0 := begin + (level-1)*pow2index_1 + prev
	parent1 := node - pow2index_1
	return parent0, parent1
}

// get graph that node belongs to, so i can find the parents
func (g *graph) getGraph(node, index int64) (int64, int64) {
	if index == 1 {
		if node < 2 {
			return 2, 0
		} else if node == 2 {
			return 1, 2
		} else if node == 3 {
			return 3, 2
		}
	}

	pow2index := int64(1 << uint64(index))
	pow2index_1 := int64(1 << uint64(index-1))
	sources := pow2index
	firstButter := sources + numButterfly(index-1)
	firstXi := firstButter + numXi(index-1)
	secondXi := firstXi + numXi(index-1)
	secondButter := secondXi + numButterfly(index-1)
	sinks := secondButter + sources

	if node < sources {
		return pow2index, 0
	} else if node >= sources && node < firstButter {
		if node < sources+pow2index_1 {
			return pow2index, pow2index_1
		} else {
			parent0, parent1 := g.butterflyParents(sources, node, index)
			return node - parent0, node - parent1
		}
	} else if node >= firstButter && node < firstXi {
		node = node - firstButter
		return g.getGraph(node, index-1)
	} else if node >= firstXi && node < secondXi {
		node = node - firstXi
		return g.getGraph(node, index-1)
	} else if node >= secondXi && node < secondButter {
		if node < secondXi+pow2index_1 {
			return pow2index_1, 0
		} else {
			parent0, parent1 := g.butterflyParents(secondXi, node, index)
			return node - parent0, node - parent1
		}
	} else if node >= secondButter && node < sinks {
		offset := (node - secondButter) % pow2index_1
		parent1 := sinks - numXi(index) + offset
		if offset+secondButter == node {
			return pow2index_1, node - parent1
		} else {
			return pow2index, node - parent1 - pow2index_1
		}
	} else {
		return 0, 0
	}
}

func (g *graph) newNodeById(id int64, hash []byte) {
	node := &node{
		H: hash,
	}
	g.writeId(node, id)
}

func (g *graph) newNode(id int64, hash []byte) {
	node := &node{
		H: hash,
	}
	g.writeNode(node, id)
}

func (g *graph) getId(id int64) *node {
	//fmt.Println("read id", id)
	node := new(node)
	data := make([]byte, nodeSize)
	num, err := g.db.ReadAt(data, id*nodeSize)
	if err != nil || num != nodeSize {
		panic(err)
	}
	node.H = data
	return node
}

func (g *graph) writeId(node *node, id int64) {
	//fmt.Println("write id", id)
	num, err := g.db.WriteAt(node.H, id*nodeSize)
	if err != nil || num != nodeSize {
		panic(err)
	}
}

func (g *graph) getNode(id int64) *node {
	idx := g.bfsToPost(id)
	//fmt.Println("read", idx)
	return g.getId(idx)
}

func (g *graph) writeNode(node *node, id int64) {
	idx := g.bfsToPost(id)
	//fmt.Println("write", idx)
	g.writeId(node, idx)
}

func (g *graph) close() {
	g.db.Close()
}

// post-order is better for disk than bfs~
func (g *graph) bfsToPost(node int64) int64 {
	return node & (^(g.pow2))
}

func numXi(index int64) int64 {
	return (1 << uint64(index)) * (index + 1) * index
}

func numButterfly(index int64) int64 {
	return 2 * (1 << uint64(index)) * index
}

func (g *graph) butterflyGraph(index int64, count *int64) {
	if index == 0 {
		index = 1
	}
	numLevel := 2 * index
	perLevel := int64(1 << uint64(index))
	begin := *count - perLevel // level 0 created outside
	// no parents at level 0
	var level, i int64
	for level = 1; level < numLevel; level++ {
		for i = 0; i < perLevel; i++ {
			var prev int64
			shift := index - level
			if level > numLevel/2 {
				shift = level - numLevel/2
			}
			if (i>>uint64(shift))&1 == 0 {
				prev = i + (1 << uint64(shift))
			} else {
				prev = i - (1 << uint64(shift))
			}
			parent0 := g.getNode(begin + (level-1)*perLevel + prev)
			parent1 := g.getNode(*count - perLevel)

			ph := append(parent0.H, parent1.H...)
			buf := make([]byte, hashSize)
			binary.PutVarint(buf, *count)
			val := append(g.pk, buf...)
			val = append(val, ph...)
			hash := sha3.Sum256(val)

			g.newNode(*count, hash[:])
			*count++
		}
	}
}

// Iterative generation of the graph
func (g *graph) xiGraphIter(index int64) {
	count := g.pow2

	stack := []int64{index, index, index, index, index}
	graphStack := []int{4, 3, 2, 1, 0}

	var i int64
	graph := 0
	pow2index := int64(1 << uint64(index))
	for i = 0; i < pow2index; i++ { //sources at this level
		buf := make([]byte, hashSize)
		binary.PutVarint(buf, count)
		val := append(g.pk, buf...)
		hash := sha3.Sum256(val)

		g.newNode(count, hash[:])
		count++
	}

	if index == 1 {
		g.butterflyGraph(index, &count)
		return
	}

	for len(stack) != 0 && len(graphStack) != 0 {
		index, stack = stack[len(stack)-1], stack[:len(stack)-1]
		graph, graphStack = graphStack[len(graphStack)-1], graphStack[:len(graphStack)-1]

		indices := []int64{index - 1, index - 1, index - 1, index - 1, index - 1}
		graphs := []int{4, 3, 2, 1, 0}

		pow2index := int64(1 << uint64(index))
		pow2index_1 := int64(1 << uint64(index-1))

		if graph == 0 {
			sources := count - pow2index
			// sources to sources of first butterfly
			// create sources of first butterly
			for i = 0; i < pow2index_1; i++ {
				parent0 := g.getNode(sources + i)
				parent1 := g.getNode(sources + i + pow2index_1)

				ph := append(parent0.H, parent1.H...)
				buf := make([]byte, hashSize)
				binary.PutVarint(buf, count)
				val := append(g.pk, buf...)
				val = append(val, ph...)
				hash := sha3.Sum256(val)

				g.newNode(count, hash[:])
				count++
			}
		} else if graph == 1 {
			firstXi := count
			// sinks of first butterfly to sources of first xi graph
			for i = 0; i < pow2index_1; i++ {
				nodeId := firstXi + i
				// index is the last level; i.e., sinks
				parent := g.getNode(firstXi - pow2index_1 + i)

				buf := make([]byte, hashSize)
				binary.PutVarint(buf, nodeId)
				val := append(g.pk, buf...)
				val = append(val, parent.H...)
				hash := sha3.Sum256(val)

				g.newNode(nodeId, hash[:])
				count++
			}
		} else if graph == 2 {
			secondXi := count
			// sinks of first xi to sources of second xi
			for i = 0; i < pow2index_1; i++ {
				nodeId := secondXi + i
				parent := g.getNode(secondXi - pow2index_1 + i)

				buf := make([]byte, hashSize)
				binary.PutVarint(buf, nodeId)
				val := append(g.pk, buf...)
				val = append(val, parent.H...)
				hash := sha3.Sum256(val)

				g.newNode(nodeId, hash[:])
				count++
			}
		} else if graph == 3 {
			secondButter := count
			// sinks of second xi to sources of second butterfly
			for i = 0; i < pow2index_1; i++ {
				nodeId := secondButter + i
				parent := g.getNode(secondButter - pow2index_1 + i)

				buf := make([]byte, hashSize)
				binary.PutVarint(buf, nodeId)
				val := append(g.pk, buf...)
				val = append(val, parent.H...)
				hash := sha3.Sum256(val)

				g.newNode(nodeId, hash[:])
				count++
			}
		} else {
			sinks := count
			sources := sinks + pow2index - numXi(index)
			for i = 0; i < pow2index_1; i++ {
				nodeId0 := sinks + i
				nodeId1 := sinks + i + pow2index_1
				parent0 := g.getNode(sinks - pow2index_1 + i)
				parent1_0 := g.getNode(sources + i)
				parent1_1 := g.getNode(sources + i + pow2index_1)

				ph := append(parent0.H, parent1_0.H...)
				buf := make([]byte, hashSize)
				binary.PutVarint(buf, nodeId0)
				val := append(g.pk, buf...)
				val = append(val, ph...)
				hash1 := sha3.Sum256(val)

				ph = append(parent0.H, parent1_1.H...)
				binary.PutVarint(buf, nodeId1)
				val = append(g.pk, buf...)
				val = append(val, ph...)
				hash2 := sha3.Sum256(val)

				g.newNode(nodeId0, hash1[:])
				g.newNode(nodeId1, hash2[:])
				count += 2
			}
		}

		if (graph == 0 || graph == 3) ||
			((graph == 1 || graph == 2) && index == 2) {
			g.butterflyGraph(index-1, &count)
		} else if graph == 1 || graph == 2 {
			stack = append(stack, indices...)
			graphStack = append(graphStack, graphs...)
		}
	}
}

func (g *graph) xiGraph(index int64, count *int64) {
	// recursively generate graphs
	// compute hashes along the way

	pow2index := int64(1 << uint64(index))

	// the first sources
	// if index == 1, then this will generate level 0 of the butterfly
	var i int64

	if *count == g.pow2 {
		for i = 0; i < pow2index; i++ {
			buf := make([]byte, hashSize)
			binary.PutVarint(buf, *count)
			val := append(g.pk, buf...)
			hash := sha3.Sum256(val)

			g.newNode(*count, hash[:])
			*count++
		}
	}

	if index == 1 {
		g.butterflyGraph(index, count)
		return
	}

	sources := *count - pow2index
	firstButter := sources + pow2index
	firstXi := firstButter + numButterfly(index-1)
	secondXi := firstXi + numXi(index-1)
	secondButter := secondXi + numXi(index-1)
	sinks := secondButter + numButterfly(index-1)
	pow2index_1 := int64(1 << uint64(index-1))

	// sources to sources of first butterfly
	// create sources of first butterly
	for i = 0; i < pow2index_1; i++ {
		parent0 := g.getNode(sources + i)
		parent1 := g.getNode(sources + i + pow2index_1)

		ph := append(parent0.H, parent1.H...)
		buf := make([]byte, hashSize)
		binary.PutVarint(buf, *count)
		val := append(g.pk, buf...)
		val = append(val, ph...)
		hash := sha3.Sum256(val)

		g.newNode(*count, hash[:])
		*count++
	}

	g.butterflyGraph(index-1, count)
	// sinks of first butterfly to sources of first xi graph
	for i = 0; i < pow2index_1; i++ {
		nodeId := firstXi + i
		parent := g.getNode(firstXi - pow2index_1 + i)

		buf := make([]byte, hashSize)
		binary.PutVarint(buf, nodeId)
		val := append(g.pk, buf...)
		val = append(val, parent.H...)
		hash := sha3.Sum256(val)

		g.newNode(nodeId, hash[:])
		*count++
	}

	g.xiGraph(index-1, count)
	// sinks of first xi to sources of second xi
	for i = 0; i < pow2index_1; i++ {
		nodeId := secondXi + i
		parent := g.getNode(secondXi - pow2index_1 + i)

		buf := make([]byte, hashSize)
		binary.PutVarint(buf, nodeId)
		val := append(g.pk, buf...)
		val = append(val, parent.H...)
		hash := sha3.Sum256(val)

		g.newNode(nodeId, hash[:])
		*count++
	}

	g.xiGraph(index-1, count)
	// sinks of second xi to sources of second butterfly
	for i = 0; i < pow2index_1; i++ {
		nodeId := secondButter + i
		parent := g.getNode(secondButter - pow2index_1 + i)

		buf := make([]byte, hashSize)
		binary.PutVarint(buf, nodeId)
		val := append(g.pk, buf...)
		val = append(val, parent.H...)
		hash := sha3.Sum256(val)

		g.newNode(nodeId, hash[:])
		*count++
	}

	// generate sinks
	// sinks of second butterfly to sinks
	// and sources to sinks directly
	g.butterflyGraph(index-1, count)
	for i = 0; i < pow2index_1; i++ {
		nodeId0 := sinks + i
		nodeId1 := sinks + i + pow2index_1
		parent0 := g.getNode(sinks - pow2index_1 + i)
		parent1_0 := g.getNode(sources + i)
		parent1_1 := g.getNode(sources + i + pow2index_1)

		ph := append(parent0.H, parent1_0.H...)
		buf := make([]byte, hashSize)
		binary.PutVarint(buf, nodeId0)
		val := append(g.pk, buf...)
		val = append(val, ph...)
		hash1 := sha3.Sum256(val)

		ph = append(parent0.H, parent1_1.H...)
		binary.PutVarint(buf, nodeId1)
		val = append(g.pk, buf...)
		val = append(val, ph...)
		hash2 := sha3.Sum256(val)

		g.newNode(nodeId0, hash1[:])
		g.newNode(nodeId1, hash2[:])
		*count += 2
	}
}
