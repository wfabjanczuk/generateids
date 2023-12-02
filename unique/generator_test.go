package unique

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

var (
	charsAB           = []byte("AB")
	charsABC          = []byte("ABC")
	charsAlphanumeric = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

type testCase struct {
	name string
	args constructorArguments
}

type constructorArguments struct {
	idsToGenerate int
	idLength      int
	charList      []byte
}

func TestNewGenerator_Validation(t *testing.T) {
	testCases := []testCase{
		{
			name: "returns error when idsToGenerate is negative",
			args: constructorArguments{-10, 1, charsAlphanumeric},
		},
		{
			name: "returns error when idsToGenerate is zero",
			args: constructorArguments{0, 1, charsAlphanumeric},
		},
		{
			name: "returns error when idLength is negative",
			args: constructorArguments{10, -1, charsAlphanumeric},
		},
		{
			name: "returns error when idLength is zero",
			args: constructorArguments{10, 0, charsAlphanumeric},
		},
		{
			name: "returns error when not enough unique combinations",
			args: constructorArguments{10, 1, charsAB},
		},
		{
			name: "returns error when character list is empty",
			args: constructorArguments{10, 1, nil},
		},
		{
			name: "returns error when character is duplicated",
			args: constructorArguments{10, 1, []byte("AA")},
		},
	}

	for _, tc := range testCases {
		runExpectedErrorTest(t, tc.name, tc.args.idsToGenerate, tc.args.idLength, tc.args.charList)
	}
}

func runExpectedErrorTest(t *testing.T, testName string, idsToGenerate, idLength int, charList []byte) {
	t.Run(testName, func(t *testing.T) {
		_, err := NewGenerator(idsToGenerate, idLength, charList)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestGenerator_ToArray(t *testing.T) {
	t.Run("returns no error when unique combinations are possible", func(t *testing.T) {
		idsToGenerate := 4

		generator, err := NewGenerator(idsToGenerate, 2, charsAB)
		idsArray, err := generator.ToArray(context.Background())

		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if len(idsArray) != idsToGenerate {
			t.Fatalf("expected %d results, got %d", idsToGenerate, len(idsArray))
		}
	})
}

func TestGenerator_ToChannel(t *testing.T) {
	t.Run("returns no error when unique combinations are possible", func(t *testing.T) {
		generator, err := NewGenerator(4, 2, charsAB)
		idsChan, err := generator.ToChannel(context.Background())

		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		idsCount := 0
		for range idsChan {
			idsCount++
		}

		if idsCount != 4 {
			t.Fatalf("expected %d results, got %d", 4, idsCount)
		}
	})
}

func TestGenerator_Uniqueness(t *testing.T) {
	constructorArgumentSets := []constructorArguments{
		{8, 3, charsAB},
		{1024, 10, charsAB},
		{1048576, 20, charsAB},
		{1024, 10, charsABC},
		{4096, 12, charsABC},
	}

	for _, args := range constructorArgumentSets {
		runUniquenessTest(t, args.idsToGenerate, args.idLength, args.charList)
	}
}

func runUniquenessTest(t *testing.T, idsToGenerate, idLength int, charList []byte) {
	testName := fmt.Sprintf("returns only unique IDs for %d idsToGenerate with %d idLength and %d total chars",
		idsToGenerate, idLength, len(charList),
	)

	t.Run(testName, func(t *testing.T) {
		generator, err := NewGenerator(idsToGenerate, idLength, charList)

		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		results, err := generator.ToArray(context.Background())
		if len(results) != idsToGenerate {
			t.Fatalf("expected %d results, got %d", idsToGenerate, len(results))
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

func TestGenerator_Context(t *testing.T) {
	t.Run("returns no ids when given cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		generator, err := NewGenerator(1024, 128, charsAlphanumeric)
		idsArray, err := generator.ToArray(ctx)

		if !errors.Is(err, ctx.Err()) {
			t.Fatalf("expected context error, got %v", err)
		}

		if len(idsArray) != 0 {
			t.Fatalf("expected no ids, got %d", len(idsArray))
		}
	})

	t.Run("stops streaming ids when context is cancelled", func(t *testing.T) {
		generator, err := NewGenerator(10*bufferSize, 128, charsAlphanumeric)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		timeout := 50 * time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		idsChan, err := generator.ToChannel(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		time.Sleep(2 * timeout)
		idsCount := 0
		for range idsChan {
			idsCount++
		}

		if idsCount > bufferSize+1 {
			t.Fatalf("expected at most %d ids, got %d", bufferSize+1, idsCount)
		}

		if !errors.Is(generator.Err(), ctx.Err()) {
			t.Fatalf("expected context error, got %v", generator.Err())
		}
	})
}

func TestGenerator_OneTimeUse(t *testing.T) {
	t.Run("can be used only once", func(t *testing.T) {
		idsToGenerate := 1
		generator, err := NewGenerator(idsToGenerate, 1, charsAB)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		goroutines := 10
		resultChan := make(chan [][]byte, goroutines)
		errChan := make(chan error, goroutines)

		wg := &sync.WaitGroup{}
		wg.Add(goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()

				idsArray, err := generator.ToArray(context.Background())
				if err != nil {
					errChan <- err
					return
				}

				resultChan <- idsArray
			}()
		}

		wg.Wait()
		close(resultChan)
		close(errChan)

		errorsCount := 0
		for err := range errChan {
			if err != nil {
				errorsCount++
				if !errors.Is(err, ErrUsed) {
					t.Fatalf("expected %v, got %v", ErrUsed, err)
				}
			}
		}

		if errorsCount != goroutines-1 {
			t.Fatalf("expected %d errors, got %d", goroutines-1, errorsCount)
		}

		var result [][]byte
		resultsCount := 0
		for idsArray := range resultChan {
			resultsCount++
			result = idsArray
		}

		if resultsCount != 1 {
			t.Fatalf("expected one array of ids, got %d", resultsCount)
		}

		if len(result) != idsToGenerate {
			t.Fatalf("expected %d results, got %d", idsToGenerate, len(result))
		}
	})
}

func BenchmarkGenerator(b *testing.B) {
	constructorArgumentSets := []constructorArguments{
		{1048576, 20, charsAB},
		{1, 128, charsAlphanumeric},
		{100, 128, charsAlphanumeric},
		{10000, 128, charsAlphanumeric},
	}

	for _, args := range constructorArgumentSets {
		runGeneratorBenchmark(b, args.idsToGenerate, args.idLength, args.charList)
	}
}

func runGeneratorBenchmark(b *testing.B, idsToGenerate, idLength int, charList []byte) {
	testName := fmt.Sprintf("generate %d unique IDs with %d length each from %d total chars",
		idsToGenerate, idLength, len(charList),
	)

	b.Run(testName, func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			generator, _ := NewGenerator(idsToGenerate, idLength, charList)
			_, _ = generator.ToArray(context.Background())
		}
	})
}
