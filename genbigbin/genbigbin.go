// Genbigbin is in itself a sample usage of bigbin
// Depending on the invocation name from the CLI, it will generate one type of another of template main code:
//
// genAppMain generates the template for an App to be added to the Big Binary
//
// genBigBinMain generates the initial Big Binary template code
//
// genStandAloneMain generates a stand alone main for the App that reuses the code in of the AppMain
//
// Use the above names as symbolic links to genbigbin to ensure the correct generation can be invoked
package main

import (
	"bigbin"
	"fmt"
	"os"
)

// init register all the genration Apps
func init() {
	bigbin.AddApp("genAppMain", GenAppMain)
	bigbin.AddApp("genBigBinMain", GenBigBinMain)
	bigbin.AddApp("genStandAloneMain", GenStandAloneMain)
}

// GenAppMain generates the template for an App to be added to the Big Binary
func GenAppMain() {
	fmt.Print(bigbin.GenerateAppMain(ensureAppPackageArgs()))
}

// GenBigBinMain generates the initial Big Binary template code
func GenBigBinMain() {
	fmt.Print(bigbin.GenerateBigBin())
}

// GenStandAloneMain generates a stand alone main for the App that reuses the code in of the AppMain
func GenStandAloneMain() {
	fmt.Print(bigbin.GenerateStandAloneMain(ensureAppPackageArgs()))
}

func main() {
	bigbin.Main()
}

// ensureAppPackageArgs checks only 2 arguments exists and returns the second, or fails the process
func ensureAppPackageArgs() string {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "Usage:\n%s full/path/to/app",
			os.Args[0])
		os.Exit(-1)
	}
	return os.Args[1]
}
