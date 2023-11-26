package generator

import (
	"fmt"
	"math/rand"
)

func Generate(n, l int, charset []byte) {
	results := make([][]byte, n)

	for r := range results {
		results[r] = make([]byte, l)
	}

	for c := 0; c < n; c++ {
		for r := 0; r < n; r++ {
			i := rand.Intn(len(charset))
			results[r][c] = charset[i]
		}
	}

	for _, id := range results {
		fmt.Println(string(id))
	}
}
