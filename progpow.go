package pow

type MixArray [numLanes][numRegs]uint32

func (m *MixArray) Copy() MixArray {
	var newMix MixArray

	for i, v1 := range m {
		for j, v2 := range v1 {
			newMix[i][j] = v2
		}
	}

	return newMix
}

type MixRngState struct {
	SrcCounter uint32
	DstCounter uint32
	SrcSeq     [numRegs]uint32
	DstSeq     [numRegs]uint32
	Rng        *Kiss99
}

func (s MixRngState) nextDst() uint32 {
	val := s.DstSeq[s.DstCounter%numRegs]

	return val
}

func (s MixRngState) nextSrc() uint32 {
	val := s.SrcSeq[s.SrcCounter%numRegs]

	return val
}

func initMixRngState(seed uint64) MixRngState {
	var z, w, jsr, jcong uint32

	z = fnv1a(FNV_OFFSET_BASIS, uint32(seed))
	w = fnv1a(z, uint32(seed>>32))
	jsr = fnv1a(w, uint32(seed))
	jcong = fnv1a(jsr, uint32(seed>>32))

	rng := NewKiss99(z, w, jsr, jcong)

	var srcSeq [numRegs]uint32
	var dstSeq [numRegs]uint32

	var i uint32
	for i = 0; i < numRegs; i++ {
		dstSeq[i] = i
		srcSeq[i] = i
	}

	for i = numRegs; i > 1; i-- {
		dstInd := rng.Next() % i
		dstSeq[i-1], dstSeq[dstInd] = dstSeq[dstInd], dstSeq[i-1]

		srcInd := rng.Next() % i
		srcSeq[i-1], srcSeq[srcInd] = srcSeq[srcInd], srcSeq[i-1]
	}

	return MixRngState{0, 0, srcSeq, dstSeq, rng}
}

type Kiss99 struct {
	z     uint32
	w     uint32
	jsr   uint32
	jcong uint32
}

func NewKiss99(z, w, jsr, jcong uint32) *Kiss99 {
	return &Kiss99{z, w, jsr, jcong}
}

func (k *Kiss99) Next() uint32 {
	k.z = 36969*(k.z&65535) + (k.z >> 16)
	k.w = 18000*(k.w&65535) + (k.w >> 16)

	k.jcong = 69069*k.jcong + 1234567

	k.jsr = k.jsr ^ (k.jsr << 17)
	k.jsr = k.jsr ^ (k.jsr >> 13)
	k.jsr = k.jsr ^ (k.jsr << 5)

	return (((k.z << 16) + k.w) ^ k.jcong) + k.jsr
}

func initMix(seed [2]uint32) MixArray {
	z := fnv1a(FNV_OFFSET_BASIS, uint32(seed[0]))
	w := fnv1a(z, uint32(seed[1]))

	var mix MixArray
	for l := range mix {
		jsr := fnv1a(w, uint32(l))
		jcong := fnv1a(jsr, uint32(l))

		rng := NewKiss99(z, w, jsr, jcong)

		for r := range mix[l] {
			mix[l][r] = rng.Next()
		}
	}

	return mix
}

func round(l1 []uint32, datasetSize uint64, r uint32, mix MixArray, seed uint64, lookup func(index uint32) []uint32) MixArray {
	state := initMixRngState(seed)
	numItems := uint32(datasetSize / (2 * 128))
	itemIndex := mix[r%numLanes][0] % numItems

	item := lookup(itemIndex)

	numWordsPerLane := len(item) / int(numLanes)
	maxOperations := max(numCacheAccesses, numMathOperations)

	for i := 0; i < maxOperations; i++ {

		if i < numCacheAccesses {
			src := state.nextSrc()
			state.SrcCounter++
			dst := state.nextDst()
			state.DstCounter++
			sel := state.Rng.Next()

			for l := 0; l < int(numLanes); l++ {
				offset := mix[l][src] % l1CacheNumItems
				mix[l][dst] = randomMerge(mix[l][dst], l1[offset], sel)
			}
		}

		if i < numMathOperations {
			srcRand := state.Rng.Next() % (numRegs * (numRegs - 1))
			src1 := srcRand % numRegs // O <= src1 < numRegs
			src2 := srcRand / numRegs // 0 <= src2 < numRegs - 1
			if src2 >= src1 {
				src2 += 1
			}

			sel1 := state.Rng.Next()
			dst := state.nextDst()
			state.DstCounter++
			sel2 := state.Rng.Next()

			for l := 0; l < int(numLanes); l++ {
				data := randomMath(mix[l][src1], mix[l][src2], sel1)
				mix[l][dst] = randomMerge(mix[l][dst], data, sel2)
			}
		}
	}

	// DAG access pattern.
	dsts := make([]uint32, numWordsPerLane)
	sels := make([]uint32, numWordsPerLane)
	for i := 0; i < numWordsPerLane; i++ {
		if i == 0 {
			dsts[i] = 0
		} else {
			dsts[i] = state.nextDst()
			state.DstCounter++
		}

		sels[i] = state.Rng.Next()
	}

	for l := 0; l < int(numLanes); l++ {
		offset := ((uint32(l) ^ r) % numLanes) * uint32(numWordsPerLane)
		for i := 0; i < numWordsPerLane; i++ {
			word := item[offset+uint32(i)]
			mix[l][dsts[i]] = randomMerge(mix[l][dsts[i]], word, sels[i])
		}
	}

	return mix
}

func hashMix(l1 []uint32, height uint64, seed [2]uint32, lookup func(index uint32) []uint32) []byte {
	mix := initMix(seed)

	number := height / periodLength
	epoch := calcEpoch(height, 7500)
	datasetSize := datasetSize(epoch)

	for i := 0; i < cntDag; i++ {
		mix = round(l1, datasetSize, uint32(i), mix, number, lookup)
	}

	var laneHash [numLanes]uint32
	for l := 0; l < int(numLanes); l++ {
		laneHash[l] = FNV_OFFSET_BASIS

		for i := 0; i < int(numRegs); i++ {
			laneHash[l] = fnv1a(laneHash[l], mix[l][i])
		}
	}

	numWords := 8
	mixHash := make([]uint32, numWords)
	for i := 0; i < numWords; i++ {
		mixHash[i] = FNV_OFFSET_BASIS
	}

	for l := 0; l < int(numLanes); l++ {
		mixHash[l%numWords] = fnv1a(mixHash[l%numWords], laneHash[l])
	}

	return uint32Array2ByteArray(mixHash)
}
