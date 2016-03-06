/*
 Package generator modifies your mains and generates their autoregistrations,
standalone versions and the bigbin to invoke them all from a single big binary.

Benefits:

1) A single binary to install multiple programs (using symbolic links to it for each program name)

2) A single binary means:

 - Several go programs get to share the same copy of runtime in disk and memory

 - Less disk space required

 - Less memory comsumption and faster loading on subsequent program invocations of the same bigbin

Source code changes automated by this Generate(bigBinDir, mainDirs...) are, for all "mainDirs":

1) Rename in all *.go files the package name from main to the last name of that directory path

2) Rename "func main()" to "func Main()"

3) Add an autoregistration file with a init() that registers this package Main to be invocable by the big binary:

	package {appname}

	import "github.com/josvazg/bigbin"

	func init() {
		bigbin.Add("{appname}", Main)
	}

4) A standalone main will be generated at {directory=appname}/{appname} wich code such as:

	package main

	import "{directory import path}"

	func main() {
		{appname}.Main()
	}

Also, if "bigBinDir" is not empty, the bigbin main is created at "bigBinDir" with code such as:

	package main

	import (
		"github.com/josvazg/bigbin"
		_ "{directory1 import path}"
		...
		_ "{directoryN import path}"
	)

	func main() {
		bigbin.Main()
	}

Note:

- That bigbin won't try to rename directories packages such as .../cmd or .../main but you might want to rename
them to something more useful and less confusing, specially if you end up with .../main/main after standalone
generation.

- Generate() does NOT apply changes, but returns all Sources that would need to be changed.
Use Sources.Apply() to enforce the changes to the file system.

- Restore() does the exact opposite to Generate() to help users undo their changes if needed.
*/
package generator
