package unique

import (
	"fmt"
	"math"
	"math/rand"
)

func newUniquenessError(totalToGenerate, eachLength, totalChars, maxToGenerate int) error {
	return fmt.Errorf(
		"impossible to generate %d unique IDs with %d length each and %d total chars; maximum of %d unique IDs can be generated",
		totalToGenerate, eachLength, totalChars, maxToGenerate,
	)
}

func validate(totalToGenerate, eachLength, totalChars int) error {
	maxToGenerate := pow(totalChars, eachLength)
	if totalToGenerate > maxToGenerate {
		return newUniquenessError(totalToGenerate, eachLength, totalChars, maxToGenerate)
	}
	return nil
}

func pow(base, exponent int) int {
	n := 1
	for i := 0; i < exponent; i++ {
		n *= base
		if n <= 0 {
			return math.MaxInt
		}
	}
	return n
}

func Generate(totalToGenerate, eachLength int, charSet []byte) ([][]byte, error) {
	totalChars := len(charSet)
	charIndices := newCharIndicesGenerator(totalChars)

	err := validate(totalToGenerate, eachLength, totalChars)
	if err != nil {
		return nil, err
	}

	results := make([][]byte, totalToGenerate)
	for r := range results {
		results[r] = make([]byte, eachLength)
	}

	var currentColumnJobs, nextColumnJobs [][]int
	currentColumnJobs = append(nextColumnJobs, createJob(totalToGenerate, charSet, charIndices))

	for c := 0; c < eachLength; c++ {
		r := 0
		for _, job := range currentColumnJobs {
			for charIndex, charOccurrences := range job {
				if charOccurrences == 0 {
					continue
				}

				for j := 0; j < charOccurrences; j++ {
					results[r+j][c] = charSet[charIndex]
				}
				r += charOccurrences

				newJob := createJob(charOccurrences, charSet, charIndices)
				if newJob != nil {
					nextColumnJobs = append(nextColumnJobs, newJob)
				}
			}
		}

		currentColumnJobs = nextColumnJobs
		nextColumnJobs = make([][]int, 0)
	}

	return results, nil
}

func createJob(totalToGenerate int, charSet []byte, charIndices *charIndicesGenerator) []int {
	totalChars := len(charSet)
	minCharOccurrences := totalToGenerate / totalChars
	job := make([]int, totalChars)

	for i := range charSet {
		job[i] = minCharOccurrences
	}

	capacityLeft := totalToGenerate - totalChars*minCharOccurrences
	for i := 0; i < capacityLeft; i++ {
		job[charIndices.next()]++
	}

	return job
}

type charIndicesGenerator struct {
	arr     []int
	length  int
	current int
}

func newCharIndicesGenerator(totalChars int) *charIndicesGenerator {
	charIndices := make([]int, totalChars)
	for i := 0; i < totalChars; i++ {
		charIndices[i] = i
	}

	return &charIndicesGenerator{
		arr:     charIndices,
		length:  totalChars,
		current: 0,
	}
}

func (ci *charIndicesGenerator) swap(i, j int) {
	ci.arr[i], ci.arr[j] = ci.arr[j], ci.arr[i]
}

func (ci *charIndicesGenerator) next() int {
	if ci.current == 0 {
		rand.Shuffle(len(ci.arr), ci.swap)
	}

	ci.current++
	if ci.current == ci.length {
		ci.current = 0
	}

	return ci.arr[ci.current]
}
