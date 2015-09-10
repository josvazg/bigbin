// BigBin allows you to reuse a single binary for multiple, normally related applications
// It uses the same pattern as BusyBox binaries, the program run depends on the command line first argument
//
// Benefits:
// 1) A single binary to install multiple programs (plus symbolic links to it with each program name)
// 2) A single binary means a single copy of the runtime and libraries used to store and load into memory
//  - Less disk space reuired in tight systems
//  - Less memory to load and faster loading on subsequent invocations of other programs in teh same bigbin
//
// How to use:
//
// 1) Rename your main (if yo have one) to xxxMain (where xxx is the name of your app normally)
//
// 2) Change the package from main to xxx_main, rename the package dir to match
//
// 3) Import the bigbin package and include an init function before such as:
//
//   import "github.com/josvazg/bigbin"
//
//   func init() {
//             bigbin.Add(xxx, xxxMain)
//   }
//
// 4) For the main use this template:
//
//   package main
//
//   import (
//              "github.com/josvazg/bigbin"
//              // add each app's xxx_main here so that hey auto register themselves
//   )
//
//   func main() {
//		bigbin.Main()
//   }
//
// 5) If you need to still have a separate main for single deployments use this template:
//
//   package main
//
//   import ".../xxx_main"
//
//   func main() {
//              xxx_main.XXXMain()
//   }

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
