package unique

type charJob struct {
	char            byte
	writesFinished  int
	writesScheduled int

	next *charJob
}

type uniformCharsGenerator struct {
	head             *charJob
	currentCharCount int
}

func newUniformCharsGenerator(idsToGenerate int, charList []byte, indicesGenerator *randomIndicesGenerator) *uniformCharsGenerator {
	totalChars := len(charList)
	minCharOccurrences := idsToGenerate / totalChars
	charOccurrencesList := make([]int, totalChars)

	for i := range charList {
		charOccurrencesList[i] = minCharOccurrences
	}

	capacityLeft := idsToGenerate - totalChars*minCharOccurrences
	for i := 0; i < capacityLeft; i++ {
		charOccurrencesList[indicesGenerator.nextIndex()]++
	}

	jobsList := &uniformCharsGenerator{}
	for charIndex, charOccurrences := range charOccurrencesList {
		if charOccurrences > 0 {
			jobsList.push(&charJob{
				char:            charList[charIndex],
				writesFinished:  0,
				writesScheduled: charOccurrences,
			})
		}
	}

	return jobsList
}

func (l *uniformCharsGenerator) empty() bool {
	if l == nil {
		return true
	}
	return l.head == nil
}

func (l *uniformCharsGenerator) push(job *charJob) {
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

func (l *uniformCharsGenerator) nextChar() byte {
	char := l.head.char
	l.head.writesFinished++

	if l.head.writesFinished == 1 {
		l.currentCharCount = l.head.writesScheduled
	}

	if l.head.writesFinished == l.head.writesScheduled {
		tmp := l.head.next
		l.head.next = nil
		l.head = tmp
	}

	return char
}
