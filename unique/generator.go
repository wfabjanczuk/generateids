package unique

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var ErrUsed = errors.New("generator can be used only once")

const bufferSize = 100

type Generator struct {
	encoder      *symmetricEncoder
	charList     []byte
	idLength     int
	idsScheduled int

	mu   sync.Mutex
	used bool
	err  error
}

func NewGenerator(idsToGenerate, idLength int, charList []byte) (*Generator, error) {
	err := validate(idsToGenerate, idLength, charList)
	if err != nil {
		return nil, err
	}

	return &Generator{
		encoder:      newSymmetricEncoder(idLength, charList),
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

	randomIndicesGen := newRandomIndicesGenerator(len(g.charList))
	columns := make([]*uniformCharsGenerator, g.idLength)
	columns[0] = newUniformCharsGenerator(g.idsScheduled, g.charList, randomIndicesGen)

	idsCreated := 0
	for idsCreated < g.idsScheduled {
		if err := ctx.Err(); err != nil {
			g.mu.Lock()
			g.err = fmt.Errorf("stopped generating ids at %d: %w", idsCreated, err)
			g.mu.Unlock()
			return
		}

		id := make([]byte, g.idLength)
		id[0] = columns[0].next()

		c := 1
		for c < g.idLength {
			uniformCharsGen := columns[c]
			if uniformCharsGen.empty() {
				uniformCharsGen = newUniformCharsGenerator(columns[c-1].currentCharCount, g.charList, randomIndicesGen)
				columns[c] = uniformCharsGen
			}

			id[c] = uniformCharsGen.next()
			c++
		}

		g.encoder.encode(id)
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
