package py_test

import (
	"testing"

	"github.com/gbbr/textmate/vendor/limetext/gopy/lib"
)

func TestRunString(t *testing.T) {
	py.Initialize()
	defer py.Finalize()
	if main, err := py.NewDict(); err != nil {
		t.Fatal(err)
	} else if g, err := py.GetBuiltins(); err != nil {
		t.Fatal(err)
	} else if err := main.SetItemString("__builtins__", g); err != nil {
		t.Fatal(err)
	} else if _, err := py.RunString("a = 'hello world!'", py.FileInput, main, nil); err != nil {
		t.Fatal(err)
	} else if a, err := main.GetItemString("a"); err != nil {
		t.Fatal(err)
	} else if b, ok := a.(*py.Unicode); !ok || b.String() != "hello world!" {
		t.Error(b, err)
	}
}
