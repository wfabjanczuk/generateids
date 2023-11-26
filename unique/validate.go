package unique

import (
	"errors"
	"fmt"
	"math"
)

var (
	ErrTotalToGenerateInvalid = errors.New("number of requested IDs must be greater than zero")
	ErrEachLengthInvalid      = errors.New("length of requested IDs must be greater than zero")

	ErrCharSetInvalid = errors.New("invalid character set")
	ErrCharSetEmpty   = fmt.Errorf("%w: empty", ErrCharSetInvalid)
)

func newCharacterDuplicatedError(duplicated byte) error {
	return fmt.Errorf("%w: duplicated character %s", ErrCharSetInvalid, string(duplicated))
}

func newUniquenessError(totalToGenerate, eachLength, totalChars, maxToGenerate int) error {
	return fmt.Errorf(
		"impossible to generate %d unique IDs with %d length each and %d total chars; maximum of %d unique IDs can be generated",
		totalToGenerate, eachLength, totalChars, maxToGenerate,
	)
}

func validate(totalToGenerate, eachLength int, charSet []byte) error {
	if totalToGenerate <= 0 {
		return ErrTotalToGenerateInvalid
	}

	if eachLength <= 0 {
		return ErrEachLengthInvalid
	}

	totalChars := len(charSet)
	if totalChars == 0 {
		return ErrCharSetEmpty
	}

	uniqueChars := make(map[byte]struct{})
	for _, char := range charSet {
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
