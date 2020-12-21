package smi

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"

	"github.com/sleepinggenius2/gosmi/smi/internal"
	"github.com/sleepinggenius2/gosmi/types"
)

const (
	DefaultErrorLevel   = 3
	DefaultGlobalConfig = "/etc/smi.conf"
	DefaultUserConfig   = ".smirc"
)

var DefaultSmiPaths []string = []string{
	"/usr/local/share/mibs/ietf",
	"/usr/local/share/mibs/iana",
	"/usr/local/share/mibs/irtf",
	"/usr/local/share/mibs/site",
	"/usr/local/share/mibs/jacobs",
	"/usr/local/share/mibs/tubs",
}

func checkInit() {
	if !internal.Initialized() {
		Init()
	}
}

// int smiInit(const char *tag)
func Init(tag ...string) bool {
	var configTag, handleName string
	if len(tag) > 0 {
		configTag = tag[0]
		handleName = strings.Join(tag, ":")
	}
	if !internal.Init(handleName) {
		return false
	}

	// Set to built-in default path, if not Windows
	if runtime.GOOS != "windows" {
		internal.SetPath(DefaultSmiPaths...)
	}

	// Read global config file, if we can
	_ = ReadConfig(DefaultGlobalConfig, configTag)

	// Read user config file, if we can
	if homedir, err := os.UserHomeDir(); err == nil {
		_ = ReadConfig(filepath.Join(homedir, DefaultUserConfig), configTag)
	}

	// Use SMIPATH environment variable, if set
	SetPath(os.Getenv("SMIPATH"))

	return true
}

// void smiExit(void)
func Exit() {
	internal.Exit()
}

// void smiSetErrorLevel(int level)
func SetErrorLevel(level int) {
	checkInit()
	internal.SetErrorLevel(level)
}

// int smiGetFlags(void)
func GetFlags() int {
	checkInit()
	return internal.GetFlags()
}

// void smiSetFlags(int userflags)
func SetFlags(userflags int) {
	checkInit()
	internal.SetFlags(userflags)
}

// char *smiGetPath(void)
func GetPath() string {
	checkInit()
	return internal.GetPath()
}

// int smiSetPath(const char *path)
func SetPath(path string) {
	paths := filepath.SplitList(path)
	if len(paths) == 0 {
		return
	}
	if paths[0] == "" {
		internal.AppendPath(paths[1:]...)
	} else if paths[len(paths)-1] == "" {
		internal.PrependPath(paths[:len(paths)-1]...)
	} else {
		internal.SetPath(paths...)
	}
}

func AddModuleFile(name string, data []byte) {
	internal.AddModuleFile(name, data)
}

// void smiSetSeverity(char *pattern, int severity)
func SetSeverity(pattern string, severity int) {
	checkInit()
	internal.SetSeverity(pattern, severity)
}

// int smiReadConfig(const char *filename, const char *tag)
func ReadConfig(filename string, tag ...string) error {
	f, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, "Open file")
	}
	defer f.Close()
	// TODO: Parse file
	return nil
}

// void smiSetErrorHandler(SmiErrorHandler smiErrorHandler)
func SetErrorHandler(smiErrorHandler types.SmiErrorHandler) {
	checkInit()
	internal.SetErrorHandler(smiErrorHandler)
}
