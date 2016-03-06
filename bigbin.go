// BigBin allows you to reuse a single binary for multiple, normally related applications
// It uses the same pattern as BusyBox binaries, the program run depends on the command line first argument

package bigbin

import (
	"fmt"
	"os"
	"path/filepath"
)

// app Main function type
type MainFunc func()

// apps the bigbin contain, that is can become
var apps = make(map[string]MainFunc)

// AddApp registers an appName to invoke the given appMain
func AddApp(appName string, appMain MainFunc) {
	apps[appName] = appMain
}

func dieOnError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

// Main runs the app named as the command line first argument
func Main() {
	cmd := os.Args[0]
	processFilename, err := filepath.EvalSymlinks(cmd)
	dieOnError(err)
	rootName := filepath.Base(processFilename)
	appName := filepath.Base(cmd)
	// Invoke the appName, if registered
	if appMain, ok := apps[appName]; ok {
		appMain()
	} else if appName == rootName { // Otherwise, if it is the root process filename, rebuild the symslinks
		fmt.Println("Rebuilding symlinks in current directory:")
		for app, _ := range apps {
			fmt.Printf(" %s -> %s\n", processFilename, app)
			os.Symlink(processFilename, app)
		}
	} else { // if all above fails, then output an error with some help and exit
		fmt.Fprintf(os.Stderr, "%s app not added into this bigbin!\n", appName)
		fmt.Fprintf(os.Stderr, "Registered apps are:\n")
		for app, _ := range apps {
			fmt.Fprintf(os.Stderr, " %s\n", app)
		}
		os.Exit(2)
	}
}
