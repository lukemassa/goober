package main

import (
	"log"
	"os"

	"github.com/lukemassa/goober/internal/goober"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	err = goober.Hexdump2(f, os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
