package unique

func GenerateArray(idsToGenerate, idLength int, charList []byte) ([][]byte, error) {
	resultsChan, err := GenerateChannel(idsToGenerate, idLength, charList)
	if err != nil {
		return nil, err
	}

	results := make([][]byte, 0, idsToGenerate)
	for id := range resultsChan {
		results = append(results, id)
	}

	return results, nil
}

func GenerateChannel(idsToGenerate, idLength int, charList []byte) (<-chan []byte, error) {
	err := validate(idsToGenerate, idLength, charList)
	if err != nil {
		return nil, err
	}

	resultsChan := make(chan []byte, 1000)
	go streamToChannel(resultsChan, idsToGenerate, idLength, charList)

	return resultsChan, nil
}

func streamToChannel(idsChan chan<- []byte, idsToGenerate, idLength int, charList []byte) {
	indicesGenerator := newRandomIndicesGenerator(len(charList))

	columns := make([]*uniformCharsGenerator, idLength)
	columns[0] = newUniformCharsGenerator(idsToGenerate, charList, indicesGenerator)

	r := 0
	for r < idsToGenerate {
		id := make([]byte, idLength)
		id[0] = columns[0].nextChar()

		c := 1
		for c < idLength {
			charsGenerator := columns[c]
			if charsGenerator.empty() {
				charsGenerator = newUniformCharsGenerator(columns[c-1].currentCharCount, charList, indicesGenerator)
				columns[c] = charsGenerator
			}

			id[c] = charsGenerator.nextChar()
			c++
		}
		idsChan <- id
		r++
	}

	close(idsChan)
}
