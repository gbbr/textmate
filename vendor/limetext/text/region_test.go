// Copyright 2013 Fredrik Ehnbom
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package text

import (
	"reflect"
	"testing"
)

// Verified against ST3
func TestRegionIntersects(t *testing.T) {
	type Test struct {
		a, b Region
		c    bool
	}
	var tests = []Test{
		{Region{10, 20}, Region{25, 35}, false},
		{Region{25, 35}, Region{10, 20}, false},
		{Region{10, 25}, Region{20, 35}, true},
		{Region{20, 35}, Region{10, 25}, true},
		{Region{10, 25}, Region{15, 20}, true},
		{Region{15, 20}, Region{10, 25}, true},
		{Region{5, 10}, Region{10, 23}, false},
		{Region{5, 10}, Region{5, 10}, true},
		{Region{0, 0}, Region{0, 0}, true},
		{Region{1, 1}, Region{1, 1}, true},
		{Region{23, 24}, Region{10, 23}, false},
	}
	for _, test := range tests {
		if res := test.a.Intersects(test.b); res != test.c {
			t.Errorf("Expected %v, but got %v, %v", test.c, res, test)
		}
	}
}

// Verified against ST3
func TestRegionIntersection(t *testing.T) {
	var tests = [][]Region{
		{{10, 20}, {25, 35}, {0, 0}},
		{{25, 35}, {10, 20}, {0, 0}},
		{{10, 25}, {20, 35}, {20, 25}},
		{{20, 35}, {10, 25}, {20, 25}},
		{{10, 25}, {15, 20}, {15, 20}},
		{{15, 20}, {10, 25}, {15, 20}},
		{{5, 10}, {10, 23}, {0, 0}},
		{{5, 10}, {5, 10}, {5, 10}},
		{{1, 1}, {1, 1}, {0, 0}},
	}
	for _, test := range tests {
		if res := test[0].Intersection(test[1]); res != test[2] {
			t.Errorf("Expected intersection %v, but got %v, %v", test[2], res, test)
		}
	}
}

func TestClip(t *testing.T) {
	tests := [][]Region{
		{{10, 20}, {25, 35}, {10, 20}},
		{{10, 20}, {0, 5}, {10, 20}},
		{{10, 20}, {0, 11}, {11, 20}},
		{{10, 20}, {0, 15}, {15, 20}},
		{{10, 20}, {15, 30}, {10, 15}},
		{{10, 20}, {20, 30}, {10, 20}},
		{{10, 20}, {0, 30}, {10, 20}},
		{{10, 20}, {10, 20}, {10, 20}},
	}
	for i := range tests {
		a := tests[i][0]
		ignoreRegion := tests[i][1]
		a = a.Clip(ignoreRegion)
		if a != tests[i][2] {
			t.Errorf("Expected %v, got: %v", tests[i][2], a)
		}
	}
}

// Verified against ST3
func TestContains(t *testing.T) {
	type Test struct {
		r   Region
		pos int
		c   bool
	}
	tests := []Test{
		{Region{0, 0}, 0, true},
		{Region{10, 10}, 10, true},
		{Region{10, 11}, 10, true},
		{Region{10, 11}, 11, true},
		{Region{10, 11}, 12, false},
		{Region{10, 11}, 9, false},
	}
	for _, test := range tests {
		if res := test.r.Contains(test.pos); res != test.c {
			t.Errorf("Expected %v, but got %v, %v, %v", test.c, res, test.r, test.pos)
		}
	}
}

// Verified against ST3
func TestCover(t *testing.T) {
	tests := []struct {
		a, b Region
		out  Region
	}{
		{Region{0, 1}, Region{1, 0}, Region{0, 1}},
		{Region{1, 0}, Region{0, 1}, Region{1, 0}},
		{Region{1, 0}, Region{5, 10}, Region{10, 0}},
		{Region{5, 10}, Region{1, 0}, Region{0, 10}},
	}
	for _, test := range tests {
		if res := test.a.Cover(test.b); !reflect.DeepEqual(res, test.out) {
			t.Errorf("Expected %v, but got %v", test.out, res)
		}
	}

}
