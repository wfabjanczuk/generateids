# streamids

```
GUJCLYPH3DSBKUGTR5W13FK77555PPBRWD89PIWGLTBORT7ROC3YLGR6SGTUJR98
S9JXJGLMK94N20KWO52EINZSLUZX3KN2FQ1ABA4MLBF8N6QW93SA68DUL2FM9ZN3
APPIQDB2C25MCPVUDL14NGELTYTMU8J8LDXNORXZDCIHCGRWDHW6VJWQJKUFZ8AC
89ZIIR252QGY18P4WKQGX29Y0PJZOVRMO8APTPBL89YHISQNFK4IO1PB7L1DUWZN
2AJEVNITF3T307WXRLPHZV97J840GGGUT06TSKH4HZEX8B5DSFZK7LRBPECXJRBJ
...
```

Library for generating a large number of unique IDs in a randomized order.
Provides a channel for retrieving them one by one.

## Description

### Constructor

```go
func NewGenerator(idsToGenerate, idLength int, charList []byte) (*Generator, error)
```

**Generator** constructor requires:

* the number of ids to generate,
* length of each id,
* list of characters (bytes) to generate the ids from.

By default, random number generator is seeded with the current time in nanoseconds.
To provide custom seed, use an alternative constructor:

```go
func NewGeneratorWithSeed(idsToGenerate, idLength int, charList []byte, seed int64) (*Generator, error)
```

### Generating ids

To start generating ids, choose the method depending on your needs:

```go
func (g *Generator) ToArray(ctx context.Context) ([][]byte, error)
func (g *Generator) ToChannel(ctx context.Context) (<-chan []byte, error)
```

**Generator** stops generating new ids when the provided context is cancelled
and wrapped context error is available, for example:

```
stopped generating ids at 59079: context deadline exceeded
```

Checking the error depends on the method:

* in case of `ToArray` method, wrapped context error can be obtained directly from the returned values,
* in case of `ToChannel` method, wrapped context error will be available after
  the returned channel is closed with the `Err` method.

```go
func (g *Generator) Err() error
```

### Warning

**Generator** struct is designed for a one-time use. Running either `ToArray` or `ToChannel` methods again
will result in an error:

```
generator can be used only once: create a new instance for another stream of ids
```

## Examples

See working examples:
* [printing ids in the console](./examples/console/main.go),
* [saving ids to a file](./examples/file/main.go).
