package unique

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrIdsToGenerateInvalid = errors.New("idsToGenerate must be greater than zero")
	ErrIdLengthInvalid      = errors.New("idLength must be greater than zero")

	ErrCharListInvalid = errors.New("invalid character list")
	ErrCharListEmpty   = fmt.Errorf("%w: empty", ErrCharListInvalid)
)

func newCharacterDuplicatedError(duplicated byte) error {
	return fmt.Errorf("%w: duplicated character %s", ErrCharListInvalid, string(duplicated))
}

func newUniquenessError(idsToGenerate, idLength, totalChars, maxToGenerate int) error {
	return fmt.Errorf(
		"impossible to generate %d unique IDs with %d length each and %d total chars; maximum of %d unique IDs can be generated",
		idsToGenerate, idLength, totalChars, maxToGenerate,
	)
}

func validate(idsToGenerate, idLength int, charList []byte) error {
	if idsToGenerate <= 0 {
		return ErrIdsToGenerateInvalid
	}

	if idLength <= 0 {
		return ErrIdLengthInvalid
	}

	totalChars := len(charList)
	if totalChars == 0 {
		return ErrCharListEmpty
	}

	uniqueChars := make(map[byte]struct{})
	for _, char := range charList {
		_, exists := uniqueChars[char]
		if exists {
			return newCharacterDuplicatedError(char)
		}
		uniqueChars[char] = struct{}{}
	}

	maxToGenerate := pow(totalChars, idLength)
	if idsToGenerate > maxToGenerate {
		return newUniquenessError(idsToGenerate, idLength, totalChars, maxToGenerate)
	}
	return nil
}

func pow(base, exponent int) int {
	n := 1
	for i := 0; i < exponent; i++ {
		n *= base
		if n <= 0 {
			return math.MaxInt
		}
	}
	return n
}
