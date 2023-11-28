package unique

import (
	"math/rand"
)

func Generate(totalToGenerate, eachLength int, charList []byte) ([][]byte, error) {
	totalChars := len(charList)
	charIndices := newCharIndicesGenerator(totalChars)

	err := validate(totalToGenerate, eachLength, charList)
	if err != nil {
		return nil, err
	}

	results := make([][]byte, totalToGenerate)
	for r := range results {
		results[r] = make([]byte, eachLength)
	}

	var currentColumnJobsList, nextColumnJobsList []charWritingJob
	currentColumnJobsList = generateJobs(nextColumnJobsList, totalToGenerate, charList, charIndices)

	for c := 0; c < eachLength; c++ {
		r := 0
		for _, job := range currentColumnJobsList {
			for j := 0; j < job.count; j++ {
				results[r+j][c] = job.char
			}
			r += job.count

			if c < eachLength-1 {
				nextColumnJobsList = generateJobs(nextColumnJobsList, job.count, charList, charIndices)
			}
		}

		currentColumnJobsList = nextColumnJobsList
		nextColumnJobsList = make([]charWritingJob, 0)
	}

	return results, nil
}

type charWritingJob struct {
	char  byte
	count int
}

func generateJobs(jobsList []charWritingJob, totalToGenerate int, charList []byte, charIndices *charIndicesGenerator) []charWritingJob {
	totalChars := len(charList)
	minCharOccurrences := totalToGenerate / totalChars
	charOccurrencesList := make([]int, totalChars)

	for i := range charList {
		charOccurrencesList[i] = minCharOccurrences
	}

	capacityLeft := totalToGenerate - totalChars*minCharOccurrences
	for i := 0; i < capacityLeft; i++ {
		charOccurrencesList[charIndices.next()]++
	}

	for charIndex, charOccurrences := range charOccurrencesList {
		if charOccurrences > 0 {
			jobsList = append(jobsList, charWritingJob{
				char:  charList[charIndex],
				count: charOccurrences},
			)
		}
	}

	return jobsList
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
