package generator

import (
	"fmt"
	"math/rand"
)

func newValidationError(totalToGenerate, eachLength, totalChars, maxToGenerate int) error {
	return fmt.Errorf(
		"impossible to generate %d unique IDs with %d length each and %d total chars; maximum of %d unique IDs can be generated",
		totalToGenerate, eachLength, totalChars, maxToGenerate,
	)
}

func validate(totalToGenerate, eachLength, totalChars int) error {
	maxToGenerate := pow(totalChars, eachLength)
	if totalToGenerate > maxToGenerate {
		return newValidationError(totalToGenerate, eachLength, totalChars, maxToGenerate)
	}
	return nil
}

func pow(base, exponent int) int {
	n := 1
	for i := 0; i < exponent; i++ {
		n *= base
	}
	return n
}

func Generate(totalToGenerate, eachLength int, charset []byte) ([][]byte, error) {
	totalChars := len(charset)

	err := validate(totalToGenerate, eachLength, totalChars)
	if err != nil {
		return nil, err
	}

	results := make([][]byte, totalToGenerate)
	for r := range results {
		results[r] = make([]byte, eachLength)
	}

	charOccurrences := generateCharOccurrences(totalToGenerate, charset)
	rowStart, rowEnd := 0, 0
	for i, occurrences := range charOccurrences {
		rowEnd += occurrences
		generateSubmatrix(results,
			0, eachLength,
			rowStart, rowEnd,
			charset[i], charset,
		)
		rowStart = rowEnd
	}

	return results, nil
}

func generateCharOccurrences(totalToGenerate int, charset []byte) []int {
	totalChars := len(charset)
	minCharOccurrences := totalToGenerate / totalChars
	charOccurrences := make([]int, totalChars)

	for i := range charset {
		charOccurrences[i] = minCharOccurrences
	}

	capacityLeft := totalToGenerate - totalChars*minCharOccurrences
	if capacityLeft == 0 {
		return charOccurrences
	}

	charIndices := make([]int, totalChars)
	for i := 0; i < totalChars; i++ {
		charIndices[i] = i
	}
	rand.Shuffle(totalChars, func(i, j int) { charIndices[i], charIndices[j] = charIndices[j], charIndices[i] })

	for i := 0; i < capacityLeft; i++ {
		charOccurrences[charIndices[i]]++
	}

	return charOccurrences
}

func generateSubmatrix(results [][]byte, columnStart, columnEnd, rowStart, rowEnd int, charColumnStart byte, charset []byte) {
	totalToGenerate := rowEnd - rowStart
	eachLength := columnEnd - columnStart

	if eachLength == 0 {
		return
	}

	for i := 0; i < totalToGenerate; i++ {
		results[rowStart+i][columnStart] = charColumnStart
	}

	charOccurrences := generateCharOccurrences(totalToGenerate, charset)
	rowEnd = rowStart
	for i, occurrences := range charOccurrences {
		rowEnd += occurrences
		generateSubmatrix(results,
			columnStart+1, columnEnd,
			rowStart, rowEnd,
			charset[i], charset,
		)
		rowStart = rowEnd
	}
}
