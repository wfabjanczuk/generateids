package unique

import (
	"fmt"
	"testing"
)

var alphanumericCharSet = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func TestGenerate(t *testing.T) {
	runExpectedErrorTest(t, "returns error when number of requested IDs is negative", -10, 1, alphanumericCharSet)
	runExpectedErrorTest(t, "returns error when number of requested IDs is zero", 0, 1, alphanumericCharSet)

	runExpectedErrorTest(t, "returns error when length of requested IDs is negative", 10, -1, alphanumericCharSet)
	runExpectedErrorTest(t, "returns error when length of requested IDs is zero", 10, 0, alphanumericCharSet)

	runExpectedErrorTest(t, "returns error when character set is empty", 10, 1, nil)
	runExpectedErrorTest(t, "returns error when character is duplicated", 10, 1, []byte("AA"))
	runExpectedErrorTest(t, "returns error when not enough unique combinations", 10, 1, []byte("AB"))

	t.Run("returns no error when unique combinations are possible", func(t *testing.T) {
		results, err := Generate(4, 2, []byte("AB"))

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

func runExpectedErrorTest(t *testing.T, testName string, totalToGenerate, eachLength int, charSet []byte) {
	t.Run(testName, func(t *testing.T) {
		_, err := Generate(10, 1, nil)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func runUniquenessTest(t *testing.T, totalToGenerate, eachLength int, charSet []byte) {
	testName := fmt.Sprintf("returns only unique IDs for %d results with %d length each and %d total chars",
		totalToGenerate, eachLength, len(charSet),
	)

	t.Run(testName, func(t *testing.T) {
		results, err := Generate(totalToGenerate, eachLength, charSet)

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

	runGenerateBenchmark(b, 1, 128, alphanumericCharSet)
	runGenerateBenchmark(b, 100, 128, alphanumericCharSet)
	runGenerateBenchmark(b, 10000, 128, alphanumericCharSet)
}

func runGenerateBenchmark(b *testing.B, totalToGenerate, eachLength int, charSet []byte) {
	testName := fmt.Sprintf("generate %d unique IDs with %d length each from %d total chars",
		totalToGenerate, eachLength, len(charSet),
	)

	b.Run(testName, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = Generate(totalToGenerate, eachLength, charSet)
		}
	})
}
