package unique

import (
	"fmt"
	"testing"
)

const alphanumericCharSet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func TestGenerate(t *testing.T) {
	t.Run("returns error when not enough unique combinations", func(t *testing.T) {
		_, err := Generate(10, 1, []byte("AB"))

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

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

func runUniquenessTest(t *testing.T, totalToGenerate, eachLength int, charset []byte) {
	testName := fmt.Sprintf("returns only unique IDs for %d results with %d length each and %d total chars",
		totalToGenerate, eachLength, len(charset),
	)

	t.Run(testName, func(t *testing.T) {
		results, err := Generate(totalToGenerate, eachLength, charset)

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

	runGenerateBenchmark(b, 1, 128, []byte(alphanumericCharSet))
	runGenerateBenchmark(b, 100, 128, []byte(alphanumericCharSet))
	runGenerateBenchmark(b, 10000, 128, []byte(alphanumericCharSet))
}

func runGenerateBenchmark(b *testing.B, totalToGenerate, eachLength int, charset []byte) {
	testName := fmt.Sprintf("generate %d unique IDs with %d length each from %d total chars",
		totalToGenerate, eachLength, len(charset),
	)

	b.Run(testName, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = Generate(totalToGenerate, eachLength, charset)
		}
	})
}
