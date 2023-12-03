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

func (cg *uniformCharsGenerator) empty() bool {
	if cg == nil {
		return true
	}

	return cg.head == nil
}

func (cg *uniformCharsGenerator) push(job *charJob) {
	if cg.head == nil {
		cg.head = job
		return
	}

	last := cg.head
	for last.next != nil {
		last = last.next
	}

	last.next = job
}

func (cg *uniformCharsGenerator) next() byte {
	char := cg.head.char
	cg.head.writesFinished++

	if cg.head.writesFinished == 1 {
		cg.currentCharCount = cg.head.writesScheduled
	}

	if cg.head.writesFinished == cg.head.writesScheduled {
		tmp := cg.head.next
		cg.head.next = nil
		cg.head = tmp
	}

	return char
}
