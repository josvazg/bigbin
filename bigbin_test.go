package bigbin

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
)

var testData = []string{"a", "b", "app1", "app2"}

const (
	BIGBIN_APPNAME = "BIGBIN_APPNAME"
)

// register all names in the init, no matter what,as the real bigbin would do
func init() {
	for _, appName := range testData {
		AddApp(appName, mainMaker(appName))
	}
}

// mainMaker generates dumb MainFuncs that just return the appName
func mainMaker(appName string) func() {
	return func() {
		fmt.Println(appName)
	}
}

// TestMain to fake the proper BigBin main behavior
func TestMain(m *testing.M) {
	switch os.Getenv(BIGBIN_APPNAME) {
	case "":
		os.Exit(m.Run())
	default:
		BigBin()
	}
}

// BigBin simulates the BigBin Main behavior on a real binary
func BigBin() {
	// hack to properly test bigbin.Main behavior
	os.Args[0] = os.Getenv(BIGBIN_APPNAME)
	Main()
}

// TestBigBin invokes each app in testData as a subprocess and checks they reply as expected
// each app is created by mainMaker, so it just should output its own name
func TestBigBin(t *testing.T) {
	for _, appName := range testData {
		output, err := Run(appName)
		if err != nil {
			t.Fatalf("BigBin failed to start app %s with error: %s", appName, err)
		}
		if strings.Trim(string(output), " \n") != appName {
			t.Fatalf("BigBin failed to execute app %s correctly expected output was '%s' but got: '%s'",
				appName, appName, output)
		}
	}
}

// Run reinvokes the test to fake a run of appName
func Run(appName string) ([]byte, error) {
	cmd := exec.Command(os.Args[0])
	cmd.Env = append(os.Environ(), BIGBIN_APPNAME+"="+appName)
	return cmd.CombinedOutput()
}
