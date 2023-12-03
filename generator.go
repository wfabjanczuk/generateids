package streamids

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/wfabjanczuk/streamids/internal"
)

var (
	ErrUsed       = errors.New("generator can be used only once")
	ErrValidation = errors.New("validation error")
)

const bufferSize = 100

type Generator struct {
	random       *rand.Rand
	encoder      *internal.SymmetricEncoder
	charList     []byte
	idLength     int
	idsScheduled int

	mu   sync.Mutex
	used bool
	err  error
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

func (g *Generator) Err() error {
	g.mu.Lock()
	defer g.mu.Unlock()

	return g.err
}

func (g *Generator) ToArray(ctx context.Context) ([][]byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	err := g.markUsed()
	if err != nil {
		return nil, err
	}

	idsChan := make(chan []byte, bufferSize)
	go g.streamToChannel(ctx, idsChan)

	results := make([][]byte, 0, g.idsScheduled)
	for id := range idsChan {
		results = append(results, id)
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	return results, g.err
}

func (g *Generator) ToChannel(ctx context.Context) (<-chan []byte, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	err := g.markUsed()
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

	idsCreated := 0
	for idsCreated < g.idsScheduled {
		if err := ctx.Err(); err != nil {
			g.mu.Lock()
			g.err = fmt.Errorf("stopped generating ids at %d: %w", idsCreated, err)
			g.mu.Unlock()
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
		idsCreated++
	}
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
