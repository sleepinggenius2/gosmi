package main

import (
	"embed"
	"encoding/json"
	"flag"
	"os"

	"github.com/min-oc/gosmi"
)

//go:embed mibs
var fs embed.FS

func main() {
	module := flag.String("m", "FIZBIN-MIB", "Module to load")
	flag.Parse()

	gosmi.Init()

	gosmi.SetFS(gosmi.NamedFS("mibs", fs))

	m, err := gosmi.GetModule(*module)
	if err != nil {
		panic(err)
	}

	nodes := m.GetNodes()
	types := m.GetTypes()

	enc := json.NewEncoder(os.Stdout)
	enc.Encode(struct {
		Module gosmi.SmiModule
		Nodes  []gosmi.SmiNode
		Types  []gosmi.SmiType
	}{
		Module: m,
		Nodes:  nodes,
		Types:  types,
	})

	gosmi.Exit()
}
