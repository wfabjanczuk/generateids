package unique

import (
	"math/rand"
)

type randomIndicesGenerator struct {
	random  *rand.Rand
	array   []int
	length  int
	current int
}

func newRandomIndicesGenerator(random *rand.Rand, total int) *randomIndicesGenerator {
	indicesArray := make([]int, total)
	for i := 0; i < total; i++ {
		indicesArray[i] = i
	}

	return &randomIndicesGenerator{
		random:  random,
		array:   indicesArray,
		length:  total,
		current: 0,
	}
}

func (ig *randomIndicesGenerator) swap(i, j int) {
	ig.array[i], ig.array[j] = ig.array[j], ig.array[i]
}

func (ig *randomIndicesGenerator) next() int {
	if ig.current == 0 {
		ig.random.Shuffle(len(ig.array), ig.swap)
	}

	ig.current++
	if ig.current == ig.length {
		ig.current = 0
	}

	return ig.array[ig.current]
}
