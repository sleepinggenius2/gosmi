package main

import (
	"log"
	"os"

	"github.com/alecthomas/repr"
	"github.com/sleepinggenius2/gosmi2/parser"
)

func main() {
	r, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	defer r.Close()
	module, err := parser.Parse(r)
	if err != nil {
		log.Fatalln(err)
	}
	_ = module
	repr.Println(module)
}
