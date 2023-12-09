package generateids

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
		runExpectedValidationErrorTest(t, tc.name, tc.args.idsToGenerate, tc.args.idLength, tc.args.charList)
	}
}

func runExpectedValidationErrorTest(t *testing.T, testName string, idsToGenerate, idLength int, charList []byte) {
	t.Run(testName, func(t *testing.T) {
		_, err := NewGenerator(idsToGenerate, idLength, charList)

		if !errors.Is(err, ErrValidation) {
			t.Errorf("expected validation error, got %v", err)
		}
	})
}

func TestGenerator_Seed(t *testing.T) {
	t.Run("generators with the same seed return the same results", func(t *testing.T) {
		seed := int64(0)
		idsToGenerate := 1024
		idLength := 128

		idsArray1 := generateIdsWithSeed(t, idsToGenerate, idLength, charsAlphanumeric, seed)
		idsArray2 := generateIdsWithSeed(t, idsToGenerate, idLength, charsAlphanumeric, seed)

		for index, id := range idsArray1 {
			if string(id) != string(idsArray2[index]) {
				t.Errorf("expected %s, got %s", id, idsArray2[index])
			}
		}
	})
}

func generateIdsWithSeed(t *testing.T, idsToGenerate, idLength int, charList []byte, seed int64) [][]byte {
	generator, err := NewGeneratorWithSeed(idsToGenerate, idLength, charList, seed)
	if err != nil {
		t.Fatalf("unexpected constructor error: %s", err)
	}

	idsArray, err := generator.Array(context.Background())
	if err != nil {
		t.Fatalf("unexpected array method error: %s", err)
	}
	if generator.InterruptionErr() != nil {
		t.Errorf("expected no interruptionErr, got %v", err)
	}

	if len(idsArray) != idsToGenerate {
		t.Errorf("expected %d results, got %d", idsToGenerate, len(idsArray))
	}

	return idsArray
}

func TestGenerator_Array(t *testing.T) {
	t.Run("returns no error when unique combinations are possible", func(t *testing.T) {
		idsToGenerate := 4

		generator, err := NewGenerator(idsToGenerate, 2, charsAB)
		if err != nil {
			t.Fatalf("unexpected constructor error: %s", err)
		}

		idsArray, err := generator.Array(context.Background())
		if err != nil {
			t.Fatalf("unexpected array method error: %s", err)
		}
		if generator.InterruptionErr() != nil {
			t.Errorf("expected no interruptionErr, got %v", err)
		}

		if len(idsArray) != idsToGenerate {
			t.Errorf("expected %d results, got %d", idsToGenerate, len(idsArray))
		}
	})
}

func TestGenerator_Channel(t *testing.T) {
	t.Run("returns no error when unique combinations are possible", func(t *testing.T) {
		generator, err := NewGenerator(4, 2, charsAB)
		if err != nil {
			t.Fatalf("unexpected constructor error: %s", err)
		}

		idsChan, err := generator.Channel(context.Background())
		if err != nil {
			t.Fatalf("unexpected channel method error: %s", err)
		}
		if generator.InterruptionErr() != nil {
			t.Errorf("expected no interruptionErr, got %v", err)
		}

		idsCount := 0
		for range idsChan {
			idsCount++
		}

		if idsCount != 4 {
			t.Errorf("expected %d results, got %d", 4, idsCount)
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
			t.Fatalf("unexpected constructor error: %s", err)
		}

		results, err := generator.Array(context.Background())
		if err != nil {
			t.Fatalf("unexpected array method error: %s", err)
		}
		if generator.InterruptionErr() != nil {
			t.Errorf("expected no interruptionErr, got %v", err)
		}

		if len(results) != idsToGenerate {
			t.Errorf("expected %d results, got %d", idsToGenerate, len(results))
		}

		uniqueIDs := make(map[string]struct{})
		for _, id := range results {
			_, exists := uniqueIDs[string(id)]
			if exists {
				t.Errorf("expected unique IDs, got duplicated %s", id)
			}
			uniqueIDs[string(id)] = struct{}{}
		}
	})
}

func TestGenerator_Context(t *testing.T) {
	t.Run("array method returns no ids when given cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		generator, err := NewGenerator(1024, 128, charsAlphanumeric)
		if err != nil {
			t.Fatalf("unexpected constructor error: %s", err)
		}

		idsArray, err := generator.Array(ctx)
		if !errors.Is(err, ctx.Err()) {
			t.Errorf("expected context error returned from array method, got %v", err)
		}
		if generator.InterruptionErr() != nil {
			t.Errorf("expected no interruptionErr, got %v", err)
		}

		if len(idsArray) != 0 {
			t.Errorf("expected no ids, got %d", len(idsArray))
		}
	})

	t.Run("channel method returns no channel when given cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()

		generator, err := NewGenerator(1024, 128, charsAlphanumeric)
		if err != nil {
			t.Fatalf("unexpected constructor error: %s", err)
		}

		idsChannel, err := generator.Channel(ctx)
		if !errors.Is(err, ctx.Err()) {
			t.Errorf("expected context error returned from array method, got %v", err)
		}
		if generator.InterruptionErr() != nil {
			t.Errorf("expected no interruptionErr, got %v", err)
		}

		if idsChannel != nil {
			t.Errorf("expected nil channel, got %v", idsChannel)
		}
	})

	t.Run("stops generating ids when context is cancelled", func(t *testing.T) {
		generator, err := NewGenerator(10*bufferSize, 128, charsAlphanumeric)
		if err != nil {
			t.Fatalf("unexpected constructor error: %s", err)
		}

		timeout := 50 * time.Millisecond
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		idsChan, err := generator.Channel(ctx)
		if err != nil {
			t.Fatalf("unexpected channel method error: %s", err)
		}

		time.Sleep(2 * timeout)
		idsCount := 0
		for range idsChan {
			idsCount++
		}

		if !errors.Is(generator.InterruptionErr(), ctx.Err()) {
			t.Errorf("expected interruptionErr to be context error, got %v", generator.InterruptionErr())
		}
		if idsCount > bufferSize+1 {
			t.Errorf("expected at most %d ids, got %d", bufferSize+1, idsCount)
		}
	})
}

func TestGenerator_OneTimeUse(t *testing.T) {
	t.Run("can be used only once", func(t *testing.T) {
		idsToGenerate := 1
		generator, err := NewGenerator(idsToGenerate, 1, charsAB)
		if err != nil {
			t.Fatalf("unexpected constructor error: %s", err)
		}

		goroutines := 10
		resultChan := make(chan [][]byte, goroutines)
		errChan := make(chan error, goroutines)

		wg := &sync.WaitGroup{}
		wg.Add(goroutines)

		for i := 0; i < goroutines; i++ {
			go func() {
				defer wg.Done()

				idsArray, err := generator.Array(context.Background())
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
			t.Errorf("expected %d errors, got %d", goroutines-1, errorsCount)
		}

		var result [][]byte
		resultsCount := 0
		for idsArray := range resultChan {
			resultsCount++
			result = idsArray
		}

		if resultsCount != 1 {
			t.Errorf("expected one array of ids, got %d", resultsCount)
		}

		if len(result) != idsToGenerate {
			t.Errorf("expected %d results, got %d", idsToGenerate, len(result))
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
			_, _ = generator.Array(context.Background())
		}
	})
}
