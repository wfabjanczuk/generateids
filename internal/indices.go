package internal

import (
	"math/rand"
)

type UniformIndicesGenerator struct {
	random  *rand.Rand
	indices []int
	length  int
	current int
}

func NewUniformIndicesGenerator(random *rand.Rand, total int) *UniformIndicesGenerator {
	indicesArray := make([]int, total)
	for i := 0; i < total; i++ {
		indicesArray[i] = i
	}

	random.Shuffle(len(indicesArray), func(i, j int) {
		indicesArray[i], indicesArray[j] = indicesArray[j], indicesArray[i]
	})

	return &UniformIndicesGenerator{
		random:  random,
		indices: indicesArray,
		length:  total,
		current: 0,
	}
}

func (ig *UniformIndicesGenerator) next() int {
	ig.current++
	if ig.current == ig.length {
		ig.current = 0
	}

	return ig.indices[ig.current]
}
