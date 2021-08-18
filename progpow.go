package pow

type mixArray [numLanes][numRegs]uint32

type mixRngState struct {
	SrcCounter uint32
	DstCounter uint32
	SrcSeq     [numRegs]uint32
	DstSeq     [numRegs]uint32
	Rng        *kiss99
}

func (s *mixRngState) nextDst() uint32 {
	val := s.DstSeq[s.DstCounter%numRegs]
	s.DstCounter++

	return val
}

func (s *mixRngState) nextSrc() uint32 {
	val := s.SrcSeq[s.SrcCounter%numRegs]
	s.SrcCounter++

	return val
}

func initmixRngState(seed uint64) *mixRngState {
	var z, w, jsr, jcong uint32

	z = fnv1a(fnvOffsetBasis, uint32(seed))
	w = fnv1a(z, uint32(seed>>32))
	jsr = fnv1a(w, uint32(seed))
	jcong = fnv1a(jsr, uint32(seed>>32))

	rng := newKiss99(z, w, jsr, jcong)

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

	return &mixRngState{0, 0, srcSeq, dstSeq, rng}
}

type kiss99 struct {
	z     uint32
	w     uint32
	jsr   uint32
	jcong uint32
}

func newKiss99(z, w, jsr, jcong uint32) *kiss99 {
	return &kiss99{z, w, jsr, jcong}
}

func (k *kiss99) Next() uint32 {
	k.z = 36969*(k.z&65535) + (k.z >> 16)
	k.w = 18000*(k.w&65535) + (k.w >> 16)

	k.jcong = 69069*k.jcong + 1234567

	k.jsr = k.jsr ^ (k.jsr << 17)
	k.jsr = k.jsr ^ (k.jsr >> 13)
	k.jsr = k.jsr ^ (k.jsr << 5)

	return (((k.z << 16) + k.w) ^ k.jcong) + k.jsr
}

func initProgpowMix(seed uint64) mixArray {
	z := fnv1a(fnvOffsetBasis, uint32(seed))
	w := fnv1a(z, uint32(seed>>32))

	var mix mixArray
	for l := range mix {
		jsr := fnv1a(w, uint32(l))
		jcong := fnv1a(jsr, uint32(l))

		rng := newKiss99(z, w, jsr, jcong)

		for r := range mix[l] {
			mix[l][r] = rng.Next()
		}
	}

	return mix
}

func progpowRound(l1 []uint32, datasetSize uint64, r uint32, mix mixArray, seed uint64, lookup func(index uint32) []uint32) mixArray {
	state := initmixRngState(seed)
	numItems := uint32(datasetSize / (2 * 128))
	itemIndex := mix[r%numLanes][0] % numItems

	item := lookup(itemIndex)

	numWordsPerLane := len(item) / int(numLanes)
	maxOperations := max(numCacheAccesses, numMathOperations)

	for i := 0; i < maxOperations; i++ {

		if i < numCacheAccesses {
			src := state.nextSrc()
			dst := state.nextDst()
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

func hashProgpowMix(l1 []uint32, height uint64, seed uint64, lookup func(index uint32) []uint32) []byte {
	mix := initProgpowMix(seed)

	number := height / periodLength
	epoch := calcEpoch(height, 7500)
	datasetSize := datasetSize(epoch)

	for i := 0; i < kawpowRounds; i++ {
		mix = progpowRound(l1, datasetSize, uint32(i), mix, number, lookup)
	}

	var laneHash [numLanes]uint32
	for l := 0; l < int(numLanes); l++ {
		laneHash[l] = fnvOffsetBasis

		for i := 0; i < int(numRegs); i++ {
			laneHash[l] = fnv1a(laneHash[l], mix[l][i])
		}
	}

	numWords := 8
	mixHash := make([]uint32, numWords)
	for i := 0; i < numWords; i++ {
		mixHash[i] = fnvOffsetBasis
	}

	for l := 0; l < int(numLanes); l++ {
		mixHash[l%numWords] = fnv1a(mixHash[l%numWords], laneHash[l])
	}

	return uint32ArrayToBytes(mixHash)
}
