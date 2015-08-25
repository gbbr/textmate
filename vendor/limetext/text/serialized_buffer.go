// Copyright 2013 Fredrik Ehnbom
// Use of this source code is governed by a 2-clause
// BSD-style license that can be found in the LICENSE file.

package text

import (
	"runtime/debug"

	"github.com/gbbr/textmate/vendor/limetext/log4go"
)

// A type that serializes all read/write operations from/to the inner buffer implementation
type (
	SerializedBuffer struct {
		inner   InnerBufferInterface
		ops     chan SerializedOperation
		lockret chan interface{}
	}
	SerializedOperation func() interface{}
)

func (s *SerializedBuffer) init(bi InnerBufferInterface) {
	s.inner = bi
	s.ops = make(chan SerializedOperation)
	s.lockret = make(chan interface{})
	go s.worker()
}

func (s *SerializedBuffer) Close() {
	if s.inner == nil {
		return
	}
	// Close s.ops indicating that we won't send new data on in this channel.
	close(s.ops)
	// The rest of the cleanup will happen once the worker has finished
	// receiving everything on the s.ops channel
}

func (s *SerializedBuffer) worker() {
	for o := range s.ops {
		func() {
			defer func() {
				if r := recover(); r != nil {
					log4go.Error("Recovered from panic: %v, %s", r, debug.Stack())
					s.lockret <- r
				}
			}()
			s.lockret <- o()
		}()
	}
	// Done processing all ops, so freeing the other resources here
	s.inner = nil
	close(s.lockret)
}

func (s *SerializedBuffer) Size() int {
	s.ops <- func() interface{} { return s.inner.Size() }
	r := <-s.lockret
	if r2, ok := r.(int); ok {
		return r2
	} else {
		return 0
	}
}

func (s *SerializedBuffer) SubstrR(re Region) []rune {
	s.ops <- func() interface{} { return s.inner.SubstrR(re) }
	r := <-s.lockret
	if r2, ok := r.([]rune); ok {
		return r2
	} else {
		log4go.Error("Error: %v", r)
		return nil
	}
}

func (s *SerializedBuffer) InsertR(point int, data []rune) error {
	s.ops <- func() interface{} { return s.inner.InsertR(point, data) }
	r := <-s.lockret
	if r2, ok := r.(error); ok {
		return r2
	}
	return nil
}

func (s *SerializedBuffer) Erase(point, length int) error {
	s.ops <- func() interface{} { return s.inner.Erase(point, length) }
	r := <-s.lockret
	if r2, ok := r.(error); ok {
		return r2
	}
	return nil
}

func (s *SerializedBuffer) Index(i int) rune {
	s.ops <- func() interface{} { return s.inner.Index(i) }
	r := <-s.lockret
	if r2, ok := r.(rune); ok {
		return r2
	} else {
		log4go.Error("Error: %v", r)
		return 0
	}
}

func (s *SerializedBuffer) RowCol(point int) (row, col int) {
	s.ops <- func() interface{} { r, c := s.inner.RowCol(point); return [2]int{r, c} }
	r := <-s.lockret
	if r2, ok := r.([2]int); ok {
		return r2[0], r2[1]
	} else {
		log4go.Error("Error: %v", r)
		return 0, 0
	}
}

func (s *SerializedBuffer) TextPoint(row, col int) (i int) {
	s.ops <- func() interface{} { return s.inner.TextPoint(row, col) }
	r := <-s.lockret
	if r2, ok := r.(int); ok {
		return r2
	} else {
		log4go.Error("Error: %v", r)
		return 0
	}
}
