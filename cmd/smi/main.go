package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

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
		moduleName, err := gosmi.LoadModule(module)
		if err != nil {
			fmt.Printf("Init Error: %s\n", err)
			return
		}
		if debug {
			fmt.Printf("Loaded module %s\n", moduleName)
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
	node, err := gosmi.GetNode(oid)
	if err != nil {
		fmt.Printf("Subtree Error: %s\n", err)
		return
	}

	subtree := node.GetSubtree()

	jsonBytes, _ := json.Marshal(subtree)
	os.Stdout.Write(jsonBytes)
}

func ModuleTrees() {
	for _, module := range modules {
		m, err := gosmi.GetModule(module)
		if err != nil {
			fmt.Printf("ModuleTrees Error: %s\n", err)
			continue
		}

		nodes := m.GetNodes()
		types := m.GetTypes()

		jsonBytes, _ := json.Marshal(struct {
			Module gosmi.SmiModule
			Nodes  []gosmi.SmiNode
			Types  []gosmi.SmiType
		}{
			Module: m,
			Nodes:  nodes,
			Types:  types,
		})
		os.Stdout.Write(jsonBytes)
	}
}
