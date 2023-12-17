package generateids

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/wfabjanczuk/generateids/internal"
)

var (
	ErrUsed       = errors.New("generator can be used only once: create a new instance for another set of ids")
	ErrValidation = errors.New("validation error")
)

const bufferSize = 100

// Generator is a one-time use structure for generating a set of unique ids given their number, length
// and list of characters (bytes). Provides Array and Channel methods that can be used depending on your needs.
// To generate another set of ids, create a new instance of the Generator.
type Generator struct {
	random          *rand.Rand
	encoder         *internal.SymmetricEncoder
	charList        []byte
	idLength        int
	idsScheduled    int
	used            bool
	interruptionErr error
	mu              sync.Mutex
}

// NewGenerator is a basic constructor that requires the number of ids to generate, length of each id
// and list of characters (bytes) to generate the ids from.
// By default, internal random number generator is seeded with the current time in nanoseconds.
func NewGenerator(idsToGenerate, idLength int, charList []byte) (*Generator, error) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return newGenerator(idsToGenerate, idLength, charList, random)
}

// NewGeneratorWithSeed is an alternative constructor that additionally requires custom seed
// for the internal random number generator.
func NewGeneratorWithSeed(idsToGenerate, idLength int, charList []byte, seed int64) (*Generator, error) {
	random := rand.New(rand.NewSource(seed))
	return newGenerator(idsToGenerate, idLength, charList, random)
}

func newGenerator(idsToGenerate, idLength int, charList []byte, random *rand.Rand) (*Generator, error) {
	err := internal.Validate(idsToGenerate, idLength, charList)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrValidation, err)
	}

	return &Generator{
		random:       random,
		encoder:      internal.NewSymmetricEncoder(random, idLength, charList),
		charList:     charList,
		idLength:     idLength,
		idsScheduled: idsToGenerate,
		used:         false,
	}, nil
}

// InterruptionErr will return wrapped context error, if the context passed to either Array or Channel
// is cancelled during the process of generating ids.
//   - in case of Array method, wrapped context error can be also obtained directly from the returned values,
//   - in case of Channel method, wrapped context error will be available only from the InterruptionErr method after
//     the returned channel is closed.
func (g *Generator) InterruptionErr() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.interruptionErr
}

func (g *Generator) setInterruptionErr(idsGenerated int, err error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.interruptionErr = fmt.Errorf("stopped generating ids at %d: %w", idsGenerated, err)
}

// Array method generates the set of ids specified in the Generator constructor and saves them into an array.
func (g *Generator) Array(ctx context.Context) ([][]byte, error) {
	err := g.start(ctx)
	if err != nil {
		return nil, err
	}

	idsChan := make(chan []byte, bufferSize)
	go g.streamToChannel(ctx, idsChan)

	results := make([][]byte, 0, g.idsScheduled)
	for id := range idsChan {
		results = append(results, id)
	}

	return results, g.InterruptionErr()
}

// Channel method starts generating the set of ids specified in the Generator constructor
// and returns a channel to retrieve them one by one.
func (g *Generator) Channel(ctx context.Context) (<-chan []byte, error) {
	err := g.start(ctx)
	if err != nil {
		return nil, err
	}

	idsChan := make(chan []byte, bufferSize)
	go g.streamToChannel(ctx, idsChan)

	return idsChan, nil
}

func (g *Generator) streamToChannel(ctx context.Context, idsChan chan<- []byte) {
	defer close(idsChan)

	uniformIndicesGen := internal.NewUniformIndicesGenerator(g.random, len(g.charList))
	columns := make([]*internal.UniformCharsGenerator, g.idLength)
	columns[0] = internal.NewUniformCharsGenerator(g.idsScheduled, g.charList, uniformIndicesGen)

	idsGenerated := 0
	for idsGenerated < g.idsScheduled {
		if err := ctx.Err(); err != nil {
			g.setInterruptionErr(idsGenerated, err)
			return
		}

		id := make([]byte, g.idLength)
		id[0] = columns[0].Next()

		columnIndex := 1
		for columnIndex < g.idLength {
			uniformCharsGen := columns[columnIndex]
			if uniformCharsGen.Empty() {
				previousColumnJobSize := columns[columnIndex-1].CurrentJobSize
				uniformCharsGen = internal.NewUniformCharsGenerator(previousColumnJobSize, g.charList, uniformIndicesGen)
				columns[columnIndex] = uniformCharsGen
			}

			id[columnIndex] = uniformCharsGen.Next()
			columnIndex++
		}

		g.encoder.Encode(id)
		idsChan <- id
		idsGenerated++
	}
}

func (g *Generator) start(ctx context.Context) error {
	err := g.markUsed()
	if err != nil {
		return err
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	return nil
}

func (g *Generator) markUsed() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.used {
		return ErrUsed
	}

	g.used = true
	return nil
}
