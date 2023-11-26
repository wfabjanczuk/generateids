package main

import "github.com/wfabjanczuk/id/generator"

func main() {
	generator.Generate(10, 10, []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"))
}
