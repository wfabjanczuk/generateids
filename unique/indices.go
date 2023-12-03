package unique

import (
	"math/rand"
)

type randomIndicesGenerator struct {
	array   []int
	length  int
	current int
}

func newRandomIndicesGenerator(total int) *randomIndicesGenerator {
	indicesArray := make([]int, total)
	for i := 0; i < total; i++ {
		indicesArray[i] = i
	}

	return &randomIndicesGenerator{
		array:   indicesArray,
		length:  total,
		current: 0,
	}
}

func (ci *randomIndicesGenerator) swap(i, j int) {
	ci.array[i], ci.array[j] = ci.array[j], ci.array[i]
}

func (ci *randomIndicesGenerator) next() int {
	if ci.current == 0 {
		rand.Shuffle(len(ci.array), ci.swap)
	}

	ci.current++
	if ci.current == ci.length {
		ci.current = 0
	}

	return ci.array[ci.current]
}
