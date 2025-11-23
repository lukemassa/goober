package main

import (
	"log"
	"os"

	"github.com/lukemassa/goober/internal/goober"
)

func main() {
	f, err := os.Open("foo.txt")
	if err != nil {
		log.Fatal(err)
	}
	err = goober.Hexdump(f, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
