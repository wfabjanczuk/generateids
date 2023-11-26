package main

import (
	"fmt"
	"log"

	"github.com/wfabjanczuk/id/generator"
)

func main() {
	results, err := generator.Generate(1000, 10, []byte("ABCD"))

	if err != nil {
		log.Fatal(err)
	}

	for _, id := range results {
		fmt.Println(string(id))
	}
}
