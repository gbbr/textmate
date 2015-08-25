package peg

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/gbbr/textmate/vendor/quarnster/parser"
)

// Diff cut'n'paste from http://golang.org/src/cmd/gofmt/gofmt.go
func diff(b1, b2 []byte) (data []byte, err error) {
	f1, err := ioutil.TempFile("", "parser")
	if err != nil {
		return
	}
	defer os.Remove(f1.Name())
	defer f1.Close()

	f2, err := ioutil.TempFile("", "parser")
	if err != nil {
		return
	}
	defer os.Remove(f2.Name())
	defer f2.Close()

	f1.Write(b1)
	f2.Write(b2)

	data, err = exec.Command("diff", "-u", f1.Name(), f2.Name()).CombinedOutput()
	if len(data) > 0 {
		// diff exits with a non-zero status when the files don't match.
		// Ignore that failure as long as we get output.
		err = nil
	}
	return

}

func TestParser(t *testing.T) {
	var p Peg
	if data, err := ioutil.ReadFile("./peg.peg"); err != nil {
		t.Fatalf("%s", err)
	} else {
		//		if !p.Grammar() {
		if !p.Parse(string(data)) {
			t.Fatalf("Didn't parse correctly")
		} else {
			if data, err := ioutil.ReadFile("./peg.go"); err != nil {
				t.Fatalf("%s", err)
			} else {
				// t.Log(p.Root)
				t.Log(p.RootNode())
				gen := parser.GoGenerator{}
				ignore := func(g parser.Generator, in string) string {
					return gen.Ignore(in)
				}
				justcall := func(g parser.Generator, in string) string {
					return g.Call(in)
				}
				gen = parser.GoGenerator{
					CustomActions: []parser.CustomAction{
						{"Spacing", ignore},
						{"Space", ignore},
						{"EndOfLine", ignore},
						{"IdentStart", justcall},
						{"IdentCont", justcall},
						{"SLASH", ignore},
						{"LEFTARROW", ignore},
						{"OPEN", ignore},
						{"CLOSE", ignore},
						{"Comment", ignore},
						{"Grammar", ignore},
					},
				}
				var data2 []byte
				s := parser.GeneratorSettings{
					Header: `/*
Copyright (c) 2012-2013 Fredrik Ehnbom
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

1. Redistributions of source code must retain the above copyright notice, this
		list of conditions and the following disclaimer.
2. Redistributions in binary form must reproduce the above copyright notice,
		this list of conditions and the following disclaimer in the documentation
		and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
(INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/
`,
					Name: "Peg",
					WriteFile: func(name, data string) error {
						if !strings.Contains(name, "test") {
							data2 = []byte(data)
						}
						return nil
					},
				}

				if err := parser.GenerateParser(p.RootNode(), &gen, s); err != nil {
					t.Fatal(err.Error())
				}
				if cmp := bytes.Compare(data, data2); cmp != 0 {
					d, _ := diff(data, data2)
					t.Log(p.RootNode())
					t.Fatalf("Generated parser isn't equal to self: %d\n%s\n", cmp, string(d))
				}
			}
		}
	}
}

func BenchmarkParser(b *testing.B) {
	var p Peg
	if data, err := ioutil.ReadFile("./peg.peg"); err != nil {
		b.Fatalf("%s", err)
	} else {
		d := string(data)
		for i := 0; i < b.N; i++ { //use b.N for looping
			p.Parse(d)
		}
	}
}
