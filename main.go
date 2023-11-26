package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/wfabjanczuk/id/unique"
)

func main() {
	now := time.Now()
	defer func() {
		fmt.Printf("duration: %v\n", time.Since(now))
	}()

	results, err := unique.Generate(1000000, 128, []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"))
	check(err)

	f, err := os.Create("results.txt")
	check(err)
	defer f.Close()

	for _, id := range results {
		_, err = f.Write(id)
		check(err)

		_, err = f.WriteString("\n")
		check(err)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
