package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/wfabjanczuk/generateids"
)

func main() {
	now := time.Now()
	defer func() {
		fmt.Printf("duration: %v\n", time.Since(now))
	}()

	alphanumericCharList := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	generator, err := generateids.NewGenerator(1_000_000, 128, alphanumericCharList)
	check(err)

	ctx := context.Background()
	//ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	//defer cancel()

	idsChan, err := generator.Channel(ctx)
	check(err)

	f, err := os.Create("results.txt")
	check(err)
	defer f.Close()

	bufWriter := bufio.NewWriter(f)
	for id := range idsChan {
		_, err = bufWriter.Write(id)
		check(err)

		_, err = bufWriter.WriteString("\n")
		check(err)
	}

	check(bufWriter.Flush())
	check(generator.InterruptionErr())
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
