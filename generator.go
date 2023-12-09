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

type Generator struct {
	random       *rand.Rand
	encoder      *internal.SymmetricEncoder
	charList     []byte
	idLength     int
	idsScheduled int

	mu              sync.Mutex
	used            bool
	interruptionErr error
}

func NewGenerator(idsToGenerate, idLength int, charList []byte) (*Generator, error) {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	return newGenerator(idsToGenerate, idLength, charList, random)
}

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
