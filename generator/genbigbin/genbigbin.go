package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/josvazg/bigbin/generator"
)

func dieOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func main() {
	var bigBinDir string
	var apply, restore bool
	flag.StringVar(&bigBinDir, "to", "", "Directory in where to generate the big binary main "+
		"(by default is empty and does not create a big binary main)")
	flag.BoolVar(&apply, "apply", false, "Apply changes to the filesystem (false by default)")
	flag.BoolVar(&restore, "restore", false, "Restore files to before the big binary changes intead (false by default)")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage:")
		fmt.Fprintln(os.Stderr, os.Args[0], "[flags]", "mainDir1 [mainDir2...]")
		fmt.Fprintln(os.Stderr, "\nFlags:")
		flag.PrintDefaults()
	}
	flag.Parse()
	mainDirs := flag.Args()
	if mainDirs == nil || len(mainDirs) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	// Generate or Restore, depending on restore flag
	var sources *generator.Sources
	if !restore {
		sources = generator.Generate(bigBinDir, mainDirs...)
	} else {
		sources = generator.Restore(bigBinDir, mainDirs...)
	}
	dieOnError(sources.SingleError())
	// if code generation was successful, apply or print
	if apply {
		dieOnError(sources.Apply())
	} else {
		fmt.Print(sources.String())
	}
}
