package internal

import (
	"math/rand"
)

type UniformIndicesGenerator struct {
	random    *rand.Rand
	indices   []int
	length    int
	current   int
	generated int
}

func NewUniformIndicesGenerator(random *rand.Rand, total int) *UniformIndicesGenerator {
	indicesArray := make([]int, total)
	for i := 0; i < total; i++ {
		indicesArray[i] = i
	}

	ig := &UniformIndicesGenerator{
		random:  random,
		indices: indicesArray,
		length:  total,
		current: 0,
	}
	ig.shuffle()

	return ig
}

func (ig *UniformIndicesGenerator) generatedAll() bool {
	return ig.generated >= ig.length
}

func (ig *UniformIndicesGenerator) shuffle() {
	ig.generated = 0
	ig.random.Shuffle(len(ig.indices), func(i, j int) {
		ig.indices[i], ig.indices[j] = ig.indices[j], ig.indices[i]
	})
}

func (ig *UniformIndicesGenerator) next() int {
	ig.generated++
	ig.current++
	if ig.current == ig.length {
		ig.current = 0
	}

	return ig.indices[ig.current]
}
