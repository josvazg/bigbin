package generator

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"testing"
)

const (
	BigBinDir = "src/somewhere.com/someones/sampleBigBin/"

	SampleDir                    = "src/somewhere.com/someones/sample/"
	SampleFilename               = SampleDir + "sample.go"
	ExpectedAutoRegisterFilename = SampleDir + "sample_autoregister.go"
	ExpectedStandAloneFilename   = SampleDir + "main/main.go"
	ExpectedBigBinFilename       = BigBinDir + "main.go"

	OriginalSample = `// Sample code
package main

import (
	"fmt"
	"os"
)
	
// sample shows the calling args
func main() {
	fmt.Println("Sample does some stuff, here args are", os.Args)
}
`

	ExpectedGeneratedSample = `// Sample code
package sample

import (
	"fmt"
	"os"
)
	
// sample shows the calling args
func Main() {
	fmt.Println("Sample does some stuff, here args are", os.Args)
}
`

	ExpectedAutoRegister = Header + `Autoregister code
package sample

import "github.com/josvazg/bigbin"

func init() {
	bigbin.Add("sample", Main)
}`

	ExpectedStandAlone = Header + `Standalone main for somewhere.com/someones/sample
package main 

import "somewhere.com/someones/sample"

func main() {
    sample.Main()
}
`

	ExpectedBigBin = Header + `bigbin main"

package main

import (
    "github.com/josvazg/bigbin"
    _ "somewhere.com/someones/sample"
)

func main() {
    bigbin.Main()
}
`
)

// TestGenerate validates that Generate creates proper sources
func TestGenerate(t *testing.T) {
	// setup
	parseDir = fakeParseDir
	absPath = fakeAbsPath
	gopath := os.Getenv("GOPATH")
	os.Setenv("GOPATH", "")
	// exercise
	sources := Generate(BigBinDir, SampleDir)
	// validate
	if sources.Errors() != nil {
		t.Fatalf("Generate failed:\n%v", sources.SingleError())
	}
	filenames := sources.Filenames()
	expectedSources := 4
	if len(filenames) != expectedSources {
		t.Fatalf("Expected %d generated sources but got %d: %v", expectedSources, len(filenames), filenames)
	}
	validate(t, sources, SampleFilename, ExpectedGeneratedSample)
	validate(t, sources, ExpectedAutoRegisterFilename, ExpectedAutoRegister)
	validate(t, sources, ExpectedStandAloneFilename, ExpectedStandAlone)
	validate(t, sources, ExpectedBigBinFilename, ExpectedBigBin)
	// restore
	os.Setenv("GOPATH", gopath)
}

//
// Helper functions and mocking infrastructure
//

// validate fails the test if actual contents are not as expected
func validate(t *testing.T, sources *Sources, filename, rawExpected string) {
	actual := sources.Source(filename)
	expectedBytes, err := format.Source(([]byte)(rawExpected))
	if err != nil {
		t.Fatalf("Fix Test contents for filename %s! expected contents do not validate!: %s", filename, err)
	}
	expected := string(expectedBytes)
	if actual == "" && expected != "" {
		t.Fatalf("Missing expected filename %s in generated sources!", filename)
	}
	if expected != actual {
		t.Fatalf("Source code generation failed for %s!\nExpected code was:\n%s\nBut got:\n%s",
			filename, expected, actual)
	}
}

// fakeParseDir parses this test code instead of filesystem directories
func fakeParseDir(fileset *token.FileSet, dir string) (map[string]*ast.Package, error) {
	if dir == SampleDir {
		packages := make(map[string]*ast.Package)
		src, err := parser.ParseFile(fileset, SampleFilename, OriginalSample, parser.ParseComments|parser.AllErrors)
		if err != nil {
			return nil, fmt.Errorf("Can't parse in memory test file: %v", err)
		}
		files := make(map[string]*ast.File)
		files[filepath.Base(SampleFilename)] = src
		pkg, err := ast.NewPackage(fileset, files, fakeImporter(), ast.NewScope(nil))
		if err != nil {
			return nil, fmt.Errorf("Can't create in memory test package: %v", err)
		}
		packages["main"] = pkg
		return packages, nil
	}
	return nil, fmt.Errorf("Can't generate test packages for unexpected directory %s", dir)
}

// fakeAbsPath returns an abs path for tests, basically all test path are in memory, imaginarious and absolute
func fakeAbsPath(dir string) (string, error) {
	return dir, nil
}

// fakeImporter just supports the imports for this tests
func fakeImporter() ast.Importer {
	return func(imports map[string]*ast.Object, path string) (pkg *ast.Object, err error) {
		if path == "fmt" || path == "os" {
			return ast.NewObj(ast.Pkg, path), nil
		}
		return nil, fmt.Errorf("Unsupported input imports=%v path=%s", imports, path)
	}
}
