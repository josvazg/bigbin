// BigBin allows you to reuse a single binary for multiple, normally related applications
// It uses the same pattern as BusyBox binaries, the program run depends on the command line first argument

package bigbin

import (
	"fmt"
	"os"
)

// app Main function type
type MainFunc func()

// apps the bigbin contain, that is can become
var apps = make(map[string]MainFunc)

// AddApp registers an appName to invoke the given appMain
func AddApp(appName string, appMain MainFunc) {
	apps[appName] = appMain
}

// Main runs the app named as the command line first argument
func Main() {
	appName := os.Args[0]
	if appMain, ok := apps[os.Args[0]]; ok {
		appMain()
	} else {
		fmt.Fprintf(os.Stderr, "%s app not added into this bigbin!\n", appName)
	}
}
