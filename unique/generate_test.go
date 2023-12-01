package unique

import (
	"fmt"
	"testing"
)

var alphanumericCharList = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func TestGenerate(t *testing.T) {
	runExpectedErrorTest(t, "returns error when number of requested IDs is negative", -10, 1, alphanumericCharList)
	runExpectedErrorTest(t, "returns error when number of requested IDs is zero", 0, 1, alphanumericCharList)

	runExpectedErrorTest(t, "returns error when length of requested IDs is negative", 10, -1, alphanumericCharList)
	runExpectedErrorTest(t, "returns error when length of requested IDs is zero", 10, 0, alphanumericCharList)

	runExpectedErrorTest(t, "returns error when character list is empty", 10, 1, nil)
	runExpectedErrorTest(t, "returns error when character is duplicated", 10, 1, []byte("AA"))
	runExpectedErrorTest(t, "returns error when not enough unique combinations", 10, 1, []byte("AB"))

	t.Run("returns no error when unique combinations are possible", func(t *testing.T) {
		results, err := GenerateArray(4, 2, []byte("AB"))

		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if len(results) != 4 {
			t.Fatalf("expected %d results, got %d", 4, len(results))
		}
	})

	runUniquenessTest(t, 8, 3, []byte("AB"))
	runUniquenessTest(t, 1024, 10, []byte("AB"))
	runUniquenessTest(t, 1048576, 20, []byte("AB"))

	runUniquenessTest(t, 1024, 10, []byte("ABC"))
	runUniquenessTest(t, 4096, 12, []byte("ABC"))

}

func runExpectedErrorTest(t *testing.T, testName string, totalToGenerate, eachLength int, charList []byte) {
	t.Run(testName, func(t *testing.T) {
		_, err := GenerateArray(10, 1, nil)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func runUniquenessTest(t *testing.T, totalToGenerate, eachLength int, charList []byte) {
	testName := fmt.Sprintf("returns only unique IDs for %d results with %d length each and %d total chars",
		totalToGenerate, eachLength, len(charList),
	)

	t.Run(testName, func(t *testing.T) {
		results, err := GenerateArray(totalToGenerate, eachLength, charList)

		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if len(results) != totalToGenerate {
			t.Fatalf("expected %d results, got %d", totalToGenerate, len(results))
		}

		uniqueIDs := make(map[string]struct{})
		for _, id := range results {
			_, exists := uniqueIDs[string(id)]
			if exists {
				t.Fatalf("expected unique IDs, got duplicated %s", id)
			}
			uniqueIDs[string(id)] = struct{}{}
		}
	})
}

func BenchmarkGenerate(b *testing.B) {
	runGenerateBenchmark(b, 1048576, 20, []byte("AB"))

	runGenerateBenchmark(b, 1, 128, alphanumericCharList)
	runGenerateBenchmark(b, 100, 128, alphanumericCharList)
	runGenerateBenchmark(b, 10000, 128, alphanumericCharList)
}

func runGenerateBenchmark(b *testing.B, totalToGenerate, eachLength int, charList []byte) {
	testName := fmt.Sprintf("generate %d unique IDs with %d length each from %d total chars",
		totalToGenerate, eachLength, len(charList),
	)

	b.Run(testName, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = GenerateArray(totalToGenerate, eachLength, charList)
		}
	})
}
