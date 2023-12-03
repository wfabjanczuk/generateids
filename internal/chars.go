package internal

type charJob struct {
	char            byte
	writesFinished  int
	writesScheduled int

	next *charJob
}

type UniformCharsGenerator struct {
	head           *charJob
	CurrentJobSize int
}

func NewUniformCharsGenerator(idsToGenerate int, charList []byte, uniformIndicesGen *UniformIndicesGenerator) *UniformCharsGenerator {
	totalChars := len(charList)
	charOccurrencesList := make([]int, totalChars)

	minCharOccurrences := idsToGenerate / totalChars
	for i := range charList {
		charOccurrencesList[i] = minCharOccurrences
	}

	capacityLeft := idsToGenerate - totalChars*minCharOccurrences
	for i := 0; i < capacityLeft; i++ {
		charOccurrencesList[uniformIndicesGen.next()]++
	}

	uniformCharsGen := &UniformCharsGenerator{}
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

func (cg *UniformCharsGenerator) Empty() bool {
	if cg == nil {
		return true
	}

	return cg.head == nil
}

func (cg *UniformCharsGenerator) Next() byte {
	char := cg.head.char
	cg.head.writesFinished++

	if cg.head.writesFinished == 1 {
		cg.CurrentJobSize = cg.head.writesScheduled
	}

	if cg.head.writesFinished == cg.head.writesScheduled {
		tmp := cg.head.next
		cg.head.next = nil
		cg.head = tmp
	}

	return char
}

func (cg *UniformCharsGenerator) push(job *charJob) {
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
