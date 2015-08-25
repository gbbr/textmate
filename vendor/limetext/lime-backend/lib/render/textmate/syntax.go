// Copyright 2013 The lime Authors.
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

// +build nope

package textmate

import (
	"github.com/gbbr/textmate/vendor/limetext/lime-backend/lib/render"
	"github.com/gbbr/textmate/vendor/limetext/lime-backend/lib/textmate"
	"github.com/gbbr/textmate/vendor/limetext/text"
)

type (
	TextmateSyntax struct {
		textmate.LanguageParser
	}
)

func (ts *TextmateSyntax) PrepareForRendering(text.Region) []render.ViewRegions {
	ts.Lock()
	defer ts.Unlock()
	return ts.flatten([]render.ViewRegions{}, "", ts.RootNode())
}
