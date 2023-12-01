package main

import (
	"fmt"
	"log"
	"math"

	"github.com/wfabjanczuk/id/unique"
)

func main() {
	alphanumericCharList := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	resultsChan, err := unique.GenerateChannel(math.MaxInt, 128, alphanumericCharList)
	if err != nil {
		log.Fatal(err)
	}

	total, mod10k := 0, 0
	for id := range resultsChan {
		total++

		if mod10k++; mod10k == 10000 {
			fmt.Println(string(id), "| total:", total)
			mod10k = 0
		}
	}
}
