package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/wfabjanczuk/streamids"
)

func main() {
	now := time.Now()
	defer func() {
		fmt.Printf("duration: %v\n", time.Since(now))
	}()

	alphanumericCharList := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	generator, err := streamids.NewGenerator(1_000_000, 128, alphanumericCharList)
	check(err)

	ctx := context.Background()
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()

	resultsChan, err := generator.ToChannel(ctx)
	check(err)

	f, err := os.Create("results.txt")
	check(err)
	defer f.Close()

	bufWriter := bufio.NewWriter(f)
	defer bufWriter.Flush()

	for id := range resultsChan {
		_, err = bufWriter.Write(id)
		check(err)

		_, err = bufWriter.WriteString("\n")
		check(err)
	}

	if err := generator.Err(); err != nil {
		check(err)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
