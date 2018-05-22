package reexec

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

var registeredInitializers = make(map[string]func())

// Register adds an initialization func under the specified name
func Register(name string, initializer func()) {
	if _, exists := registeredInitializers[name]; exists {
		panic(fmt.Sprintf("reexec func already registred under name %q", name))
	}

	registeredInitializers[name] = initializer
}

// Init is called as the first part of the exec process and returns true if an
// initialization function was called.
func Init() bool {
	//args[0]获取的是相对路径，或者说，就是你使用什么命令启动的。
	//如果你用./a启动的话，args[0]就是./a，不是绝对路径。
	//如果你用./XXX/a启动的话，args[0]就是./XXX/a，不是绝对路径。
	//如果用/home/XXX/a启动，那么获取到的就是/home/XXX/a。
	//argv[0]的做法来自C语言，因此其他语言的argv[0]也就保持了和C语言一致。
	initializer, exists := registeredInitializers[os.Args[0]]
	if exists {
		initializer()

		return true
	}

	return false
}

// Self returns the path to the current processes binary
func Self() string {
	name := os.Args[0]

	if filepath.Base(name) == name {
		if lp, err := exec.LookPath(name); err == nil {
			name = lp
		}
	}

	return name
}
