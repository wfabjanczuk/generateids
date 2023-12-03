package unique

import (
	"math/rand"
)

type symmetricEncoder struct {
	idLength      int
	pairEncodings map[pair]pair
}

type pair struct {
	c1 byte
	c2 byte
}

func newSymmetricEncoder(random *rand.Rand, idLength int, charList []byte) *symmetricEncoder {
	totalChars := len(charList)

	pairs := make([]pair, 0, totalChars*totalChars)
	for i := 0; i < totalChars; i++ {
		for j := 0; j < totalChars; j++ {
			pairs = append(pairs, pair{charList[i], charList[j]})
		}
	}

	shuffledPairs := make([]pair, totalChars*totalChars)
	copy(shuffledPairs, pairs)
	random.Shuffle(totalChars*totalChars, func(i, j int) {
		shuffledPairs[i], shuffledPairs[j] = shuffledPairs[j], shuffledPairs[i]
	})

	pairEncodings := make(map[pair]pair, totalChars*totalChars)
	for i, p := range pairs {
		pairEncodings[p] = shuffledPairs[i]
	}

	return &symmetricEncoder{
		idLength:      idLength,
		pairEncodings: pairEncodings,
	}
}

func (e *symmetricEncoder) encode(id []byte) {
	i, j := 0, e.idLength-1

	for i+1 < j {
		encoding := e.pairEncodings[pair{id[i], id[j]}]
		id[i] = encoding.c1
		id[j] = encoding.c2

		i++
		j--
	}
}
