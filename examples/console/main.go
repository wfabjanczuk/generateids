package main

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/wfabjanczuk/streamids"
)

func main() {
	now := time.Now()

	alphanumericCharList := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	generator, err := streamids.NewGenerator(math.MaxInt, 128, alphanumericCharList)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()

	idsChan, err := generator.ToChannel(ctx)
	if err != nil {
		panic(err)
	}

	total, mod1k := 0, 0
	for id := range idsChan {
		total++
		if mod1k++; mod1k == 1000 {
			fmt.Println(string(id), "| total:", total, "| duration:", time.Since(now))
			mod1k = 0
		}
	}

	if err := generator.Err(); err != nil {
		panic(err)
	}
}
