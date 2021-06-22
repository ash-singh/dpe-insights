package testing

import (
	"log"
	"os"
	"path"
	"runtime"
	"testing"
)

func init() {
	changeWorkingDirectoryToRoot()
}

func changeWorkingDirectoryToRoot() {
	_, filePath, _, _ := runtime.Caller(0)
	dir := path.Join(path.Dir(filePath), "..")
	err := os.Chdir(dir)
	if err != nil {
		log.Fatal(err)
	}
}

// CheckResponseCode check http response code.
func CheckResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
