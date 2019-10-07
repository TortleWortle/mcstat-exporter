package internal_test

import (
	"fmt"
	"testing"

	"github.com/tortlewortle/mcstat-exporter/internal"
)

func TestCheckPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()
	internal.Check(fmt.Errorf("Some error"))
}

func TestCheckDoesntPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("The code did panic")
		}
	}()
	internal.Check(nil)
}
