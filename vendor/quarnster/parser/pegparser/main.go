/*
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

package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gbbr/textmate/vendor/quarnster/parser"
	"github.com/gbbr/textmate/vendor/quarnster/parser/peg"
)

func main() {
	var (
		pegfile    = ""
		testfile   = ""
		bench      = false
		debug      = 0
		dumptree   = false
		notest     = false
		heatmap    = false
		ignore     = ""
		generator  = "go"
		outpath    = ""
		outfile    = ""
		typename   = ""
		header     = "default"
		gogenerate = false
	)
	flag.StringVar(&ignore, "ignore", ignore, "List of definitions to ignore (not generate nodes for)")
	flag.StringVar(&pegfile, "peg", pegfile, "Pegfile for which to generate a parser for")
	flag.StringVar(&testfile, "testfile", testfile, "Test file to be used in testing")
	flag.StringVar(&outpath, "outpath", outpath, "Destination directory path")
	flag.StringVar(&outfile, "outfile", outfile, "Destination file")
	flag.BoolVar(&bench, "bench", bench, "Whether to run a benchmark test or not")
	flag.IntVar(&debug, "debug", debug, "The desired debug level the generated parser will use")
	flag.BoolVar(&dumptree, "dumptree", dumptree, "Whether to make the generated parser spit out the generated tree")
	flag.BoolVar(&notest, "notest", notest, "Whether to test the generated parser")
	flag.BoolVar(&heatmap, "heatmap", heatmap, "Whether to generate a heatmap or not")
	flag.StringVar(&generator, "generator", generator, "Which generator to use")
	flag.StringVar(&header, "header", header, "Header to put at the top of the generated source code")
	flag.StringVar(&typename, "name", typename, "Name of the generated type/namespace/package. By default it'll be based on the name of the .peg-file")
	flag.BoolVar(&gogenerate, "gogenerate", gogenerate, "Add a Go 1.4 \"//go:generate\" line to the generated code")
	flag.Parse()
	if pegfile == "" {
		flag.Usage()
		os.Exit(1)
	}
	p := peg.Peg{}
	if data, err := ioutil.ReadFile(pegfile); err != nil {
		log.Fatalf("%s", err)
	} else {
		if !p.Parse(string(data)) {
			log.Fatalf("Didn't parse correctly\n")
		} else {
			back := p.RootNode().Children[len(p.RootNode().Children)-1]
			if back.Name != "EndOfFile" {
				log.Println(p.RootNode())
				log.Println("File didn't finish parsing")
			}
			name := outfile
			if name == "" {
				name = filepath.Base(pegfile)
				name = name[:len(name)-len(filepath.Ext(name))]
			}

			ignoreFunc := func(g parser.Generator, in string) string {
				return g.Ignore(in)
			}
			var customActions []parser.CustomAction
			for _, action := range strings.Split(ignore, ",") {
				action = strings.TrimSpace(action)
				customActions = append(customActions, parser.CustomAction{action, ignoreFunc})
			}

			var gen parser.Generator
			switch generator {
			case "go":
				gen = &parser.GoGenerator{RootNode: p.RootNode()}
			case "c":
				gen = &parser.CGenerator{}
			case "cpp":
				gen = &parser.CPPGenerator{}
			case "java":
				gen = &parser.JavaGenerator{}
			case "py":
				gen = &parser.PyGenerator{}
			default:
				panic(generator)
			}

			//			gen.AddDebugLogging = debug
			root := outpath
			if root == "" {
				root = filepath.Dir(pegfile)
				if generator != "go" {
					root += "_" + generator
				}
			}
			root += "/"
			gen.SetCustomActions(customActions)
			if header == "default" {
				header = ""
				gogenerate = true
			}
			if gogenerate {
				header += "//go:generate"
				for _, a := range os.Args {
					header += " \"" + strings.Replace(strings.Replace(a, "\n", "\\n", -1), "\"", "\\\"", -1) + `"`
				}
				header += "\n"
			}
			if typename == "" {
				typename = filepath.Base(pegfile)
				typename = strings.ToTitle(typename[:len(typename)-len(filepath.Ext(typename))])
			}
			s := parser.GeneratorSettings{
				Header:     header,
				Name:       typename,
				Testname:   testfile,
				FileName:   outfile,
				Debug:      dumptree,
				DebugLevel: parser.DebugLevel(debug),
				Bench:      bench,
				Heatmap:    heatmap,
				WriteFile: func(name, data string) error {
					if err := os.Mkdir(root, 0755); err != nil && !os.IsExist(err) {
						return err
					}
					if err := ioutil.WriteFile(root+name, []byte(data), 0644); err != nil {
						return err
					}
					return nil
				},
			}
			if err := parser.GenerateParser(p.RootNode(), gen, s); err != nil {
				log.Fatalln(err)
			} else if !notest {
				cmd := gen.TestCommand()
				c := exec.Command(cmd[0], cmd[1:]...)
				c.Dir = root
				data, _ := c.CombinedOutput()
				log.Println(string(data))
			}
		}
	}
}
