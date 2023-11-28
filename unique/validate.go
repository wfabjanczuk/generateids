package unique

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrTotalToGenerateInvalid = errors.New("number of requested IDs must be greater than zero")
	ErrEachLengthInvalid      = errors.New("length of requested IDs must be greater than zero")

	ErrCharListInvalid = errors.New("invalid character list")
	ErrCharListEmpty   = fmt.Errorf("%w: empty", ErrCharListInvalid)
)

func newCharacterDuplicatedError(duplicated byte) error {
	return fmt.Errorf("%w: duplicated character %s", ErrCharListInvalid, string(duplicated))
}

func newUniquenessError(totalToGenerate, eachLength, totalChars, maxToGenerate int) error {
	return fmt.Errorf(
		"impossible to generate %d unique IDs with %d length each and %d total chars; maximum of %d unique IDs can be generated",
		totalToGenerate, eachLength, totalChars, maxToGenerate,
	)
}

func validate(totalToGenerate, eachLength int, charList []byte) error {
	if totalToGenerate <= 0 {
		return ErrTotalToGenerateInvalid
	}

	if eachLength <= 0 {
		return ErrEachLengthInvalid
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

	maxToGenerate := pow(totalChars, eachLength)
	if totalToGenerate > maxToGenerate {
		return newUniquenessError(totalToGenerate, eachLength, totalChars, maxToGenerate)
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
