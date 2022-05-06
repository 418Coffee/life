package life

import (
	"reflect"
	"testing"
)

func TestExample(t *testing.T) {
	l, err := LoadGame("./examples/glider.rle", true)
	if err != nil {
		t.Fatal(err)
	}
	if l.width != 3 || l.height != 3 {
		t.Fatalf("unexpected width height: %d %d", l.width, l.height)
	}
	if !reflect.DeepEqual(l.current.s[0], []bool{false, true, false}) {
		t.FailNow()
	}
	if !reflect.DeepEqual(l.current.s[1], []bool{false, false, true}) {
		t.FailNow()
	}
	if !reflect.DeepEqual(l.current.s[2], []bool{true, true, true}) {
		t.FailNow()
	}
}
