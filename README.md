# gosmi

Go wrapper around libsmi

## Usage
On Ubuntu: `$ sudo apt-get install libsmi2-dev`

### Example
```go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"os"

	"github.com/sleepinggenius2/gosmi"
)

type arrayStrings []string

var modules arrayStrings
var paths arrayStrings
var debug bool

func (a arrayStrings) String() string {
	return strings.Join(a, ",")
}

func (a *arrayStrings) Set(value string) error {
	*a = append(*a, value)
	return nil
}

func main() {
	flag.BoolVar(&debug, "d", false, "Debug")
	flag.Var(&modules, "m", "Module to load")
	flag.Var(&paths, "p", "Path to add")
	flag.Parse()

	Init()

	oid := flag.Arg(0)
	if oid == "" {
		ModuleTrees()
	} else {
		Subtree(oid)
	}

	Exit()
}

func Init() {
	gosmi.Init()

	for _, path := range paths {
		gosmi.AppendPath(path)
	}

	for _, module := range modules {
		moduleName, ok := gosmi.LoadModule(module)
		if !ok {
			fmt.Println("Failed to load module %s\n", module)
			return
		}
	}

	if debug {
		path := gosmi.GetPath()
		fmt.Printf("Search path: %s\n", path)
		loadedModules := gosmi.GetLoadedModules()
		fmt.Println("Loaded modules:")
		for _, loadedModule := range loadedModules {
			fmt.Printf("  %s (%s)\n", loadedModule.Name, loadedModule.Path)
		}
	}
}

func Exit() {
	gosmi.Exit()
}

func Subtree(oid string) {
	node, ok := gosmi.GetNode(oid)
	if !ok {
		fmt.Println("Invalid OID")
		return
	}

	subtree := node.GetSubtree()

	jsonBytes, _ := json.Marshal(subtree)
	os.Stdout.Write(jsonBytes)
}

func ModuleTrees() {
	for _, module := range modules {
		m, ok := gosmi.GetModule(module)
		if !ok {
			fmt.Printf("Could not display %s\n", module)
			continue
		}

		nodes := m.GetNodes()
		types := m.GetTypes()

		if jsonOutput {
			jsonBytes, _ := json.Marshal(struct{
				Module gosmi.Module
				Nodes []gosmi.Node
				Types []gosmi.Type
			}{
				Module: m,
				Nodes: nodes,
				Types: types,
			})
			os.Stdout.Write(jsonBytes)
		}
	}
}
```
