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

func newUniformCharsGenerator(idsToGenerate int, charList []byte, randomIndicesGen *randomIndicesGenerator) *uniformCharsGenerator {
	totalChars := len(charList)
	charOccurrencesList := make([]int, totalChars)

	minCharOccurrences := idsToGenerate / totalChars
	for i := range charList {
		charOccurrencesList[i] = minCharOccurrences
	}

	capacityLeft := idsToGenerate - totalChars*minCharOccurrences
	for i := 0; i < capacityLeft; i++ {
		charOccurrencesList[randomIndicesGen.next()]++
	}

	uniformCharsGen := &uniformCharsGenerator{}
	for charIndex, charOccurrences := range charOccurrencesList {
		if charOccurrences > 0 {
			uniformCharsGen.push(&charJob{
				char:            charList[charIndex],
				writesFinished:  0,
				writesScheduled: charOccurrences,
			})
		}
	}

	return uniformCharsGen
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

func (l *uniformCharsGenerator) next() byte {
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
