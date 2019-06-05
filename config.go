package gosmi

import (
	"os"

	"github.com/sleepinggenius2/gosmi/smi"
)

func Init() {
	if !smi.Init("gosmi") {
		panic("Failed to initialize")
	}
}

func Exit() { smi.Exit() }

func GetPath() string         { return smi.GetPath() }
func SetPath(path string)     { smi.SetPath(path) }
func AppendPath(path string)  { smi.SetPath(string(os.PathListSeparator) + path) }
func PrependPath(path string) { smi.SetPath(path + string(os.PathListSeparator)) }

func ReadConfig(filename string, tag ...string) error { return smi.ReadConfig(filename, tag...) }
