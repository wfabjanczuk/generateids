package internal

import (
	"math/rand"
)

type SymmetricEncoder struct {
	end           int
	pairEncodings map[pair]pair

	odd             bool
	mid             int
	singleEncodings map[byte]byte
}

type pair struct {
	c1 byte
	c2 byte
}

func NewSymmetricEncoder(random *rand.Rand, idLength int, charList []byte) *SymmetricEncoder {
	e := &SymmetricEncoder{}

	e.setupPairEncodings(random, idLength, charList)
	if e.odd = idLength%2 == 1; e.odd {
		e.setupMidEncoding(random, idLength, charList)
	}

	return e
}

func (e *SymmetricEncoder) setupPairEncodings(random *rand.Rand, idLength int, charList []byte) {
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

	e.end = idLength - 1
	e.pairEncodings = pairEncodings
}

func (e *SymmetricEncoder) setupMidEncoding(random *rand.Rand, idLength int, charList []byte) {
	totalChars := len(charList)

	shuffledChars := make([]byte, totalChars)
	copy(shuffledChars, charList)
	random.Shuffle(totalChars, func(i, j int) {
		shuffledChars[i], shuffledChars[j] = shuffledChars[j], shuffledChars[i]
	})

	singleEncodings := make(map[byte]byte, totalChars)
	for i, c := range charList {
		singleEncodings[c] = shuffledChars[i]
	}

	e.odd = true
	e.mid = idLength / 2
	e.singleEncodings = singleEncodings
}

func (e *SymmetricEncoder) Encode(id []byte) {
	i, j := 0, e.end
	for i < j {
		encoding := e.pairEncodings[pair{id[i], id[j]}]
		id[i] = encoding.c1
		id[j] = encoding.c2

		i++
		j--
	}

	if e.odd {
		id[e.mid] = e.singleEncodings[id[e.mid]]
	}
}
