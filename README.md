# generateids

```
YMRMT2OX5WS7JUDSYEREBEFQR73Z407ZCPS2IJTV33ADVUAPPBTO0H39RK8LXSLFDL8SR17DUUGUORZV
W8Y271N4NGE6SA4XYYUX3A1609BRT11FCG9UFLBKJNLHTBH58DUTAORSCM23NZ2QEBMLU246CZ9KWCJT
497XIFS514YIWOFAPPH0KU56YOU1LFGZRPLB098GL07ZS4599PIEOVY0X3MDGSW6704VZ4JVHFWLCO1S
Z7QAUKXAQ0XTQ435JH32F20AJ9CDY8QMTIEEL709ZCIRP2PTWLAORJRBV8HGESEIE972M6KA8WC5YZPN
BTTPM4B9LKF1RT3SHJ36ARTLFTWQD4TAYLNTPNT6TA6CEBNQV4EBILNYU0J29KZT9SBI0KFYBIXSP1LG
...
```

Toy library for generating a large number of unique IDs in a randomized order.
Provides a channel for retrieving them one by one.

## Installation

```shell
go get github.com/wfabjanczuk/generateids 
```

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

To generate ids, choose the method depending on your needs:

```go
func (g *Generator) Array(ctx context.Context) ([][]byte, error)
func (g *Generator) Channel(ctx context.Context) (<-chan []byte, error)
```

`Channel` method returns a buffered channel, which is closed when the job is finished.

If the provided context is cancelled during the process of generating ids, 
wrapped context error is available from `InterruptionErr` method:

```go
func (g *Generator) InterruptionErr() error
// stopped generating ids at 59079: context deadline exceeded
```

* in case of `Array` method, wrapped context error can be also obtained directly from the returned values,
* in case of `Channel` method, wrapped context error will be available only from the `InterruptionErr` method after
  the returned channel is closed.

### Warning

**Generator** struct is designed for a one-time use. Running either `Array` or `Channel` methods again
will result in an error:

```
generator can be used only once: create a new instance for another set of ids
```

## Examples

See working examples:
* [printing ids in the console](./examples/console/main.go),
* [saving ids to a file](./examples/file/main.go).
