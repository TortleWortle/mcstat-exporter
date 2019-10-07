package internal_test

import (
	"path/filepath"
	"testing"

	"github.com/tortlewortle/mcstat-exporter/internal"
)

func TestFileExists(t *testing.T) {
	path, _ := filepath.Abs("../test/fileThatExists")
	if !internal.FileExists(path) {
		t.Errorf("%s exists but returns false", path)
	}
	path, _ = filepath.Abs("../test/fileThatDoesntExists")
	if internal.FileExists(path) {
		t.Errorf("%s doesn't exist but returns true", path)
	}
}
