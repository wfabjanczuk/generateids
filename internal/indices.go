package internal

import (
	"math/rand"
)

type RandomIndicesGenerator struct {
	random  *rand.Rand
	array   []int
	length  int
	current int
}

func NewRandomIndicesGenerator(random *rand.Rand, total int) *RandomIndicesGenerator {
	indicesArray := make([]int, total)
	for i := 0; i < total; i++ {
		indicesArray[i] = i
	}

	return &RandomIndicesGenerator{
		random:  random,
		array:   indicesArray,
		length:  total,
		current: 0,
	}
}

func (ig *RandomIndicesGenerator) Next() int {
	if ig.current == 0 {
		ig.random.Shuffle(len(ig.array), ig.swap)
	}

	ig.current++
	if ig.current == ig.length {
		ig.current = 0
	}

	return ig.array[ig.current]
}

func (ig *RandomIndicesGenerator) swap(i, j int) {
	ig.array[i], ig.array[j] = ig.array[j], ig.array[i]
}
