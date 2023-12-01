package unique

import (
	"math/rand"
)

func GenerateArray(totalToGenerate, eachLength int, charList []byte) ([][]byte, error) {
	resultsChan, err := GenerateChannel(totalToGenerate, eachLength, charList)
	if err != nil {
		return nil, err
	}

	results := make([][]byte, 0, totalToGenerate)
	for id := range resultsChan {
		results = append(results, id)
	}

	return results, nil
}

func GenerateChannel(totalToGenerate, eachLength int, charList []byte) (<-chan []byte, error) {
	err := validate(totalToGenerate, eachLength, charList)
	if err != nil {
		return nil, err
	}

	resultsChan := make(chan []byte, 1000)
	go streamToChannel(resultsChan, totalToGenerate, eachLength, charList)

	return resultsChan, nil
}

func streamToChannel(resultsChan chan<- []byte, totalToGenerate, eachLength int, charList []byte) {
	totalChars := len(charList)
	charIndices := newCharIndicesGenerator(totalChars)

	columns := make([]*charWritingJobsList, eachLength)
	columns[0] = generateJobsList(totalToGenerate, charList, charIndices)

	r := 0
	for r < totalToGenerate {
		id := make([]byte, eachLength)
		id[0] = columns[0].nextChar()

		c := 1
		for c < eachLength {
			jobsList := columns[c]
			if jobsList.isEmpty() {
				jobsList = generateJobsList(columns[c-1].currentCharCount, charList, charIndices)
				columns[c] = jobsList
			}

			id[c] = jobsList.nextChar()
			c++
		}
		resultsChan <- id
		r++
	}

	close(resultsChan)
}

type charWritingJobsList struct {
	head             *charWritingJob
	currentCharCount int
}

func (l *charWritingJobsList) isEmpty() bool {
	if l == nil {
		return true
	}
	return l.head == nil
}

func (l *charWritingJobsList) push(job *charWritingJob) {
	if l.head == nil {
		l.head = job
		return
	}

	last := l.head
	for last.next != nil {
		last = last.next
	}

	last.next = job
}

func (l *charWritingJobsList) nextChar() byte {
	char := l.head.char
	l.head.written++

	if l.head.written == 1 {
		l.currentCharCount = l.head.count
	}

	if l.head.written == l.head.count {
		tmp := l.head.next
		l.head.next = nil
		l.head = tmp
	}

	return char
}

type charWritingJob struct {
	char    byte
	count   int
	written int

	next *charWritingJob
}

func generateJobsList(totalToGenerate int, charList []byte, charIndices *charIndicesGenerator) *charWritingJobsList {
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

	jobsList := &charWritingJobsList{}
	for charIndex, charOccurrences := range charOccurrencesList {
		if charOccurrences > 0 {
			jobsList.push(&charWritingJob{
				char:  charList[charIndex],
				count: charOccurrences,
			})
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
