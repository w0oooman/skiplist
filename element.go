package skiplist

import (
	"math/rand"
)

const (
	SkipListMaxLevel = 32
	SkipListBranch   = 2
)

type skiplistLevel struct {
	forward *Element
	span    int
}

type Element struct {
	key      interface{}
	Value    Interface
	backward *Element
	level    []skiplistLevel
}

// Next returns the next skiplist element or nil.
func (e *Element) Next() *Element {
	return e.level[0].forward
}

// Prev returns the previous skiplist element of nil.
func (e *Element) Prev() *Element {
	return e.backward
}

// newElement returns an initialized element.
func newElement(level int, k interface{}, v Interface) *Element {
	return &Element{
		key:      k,
		Value:    v,
		backward: nil,
		level:    make([]skiplistLevel, level),
	}
}

// randomLevel returns a random level.
func randomLevel() int {
	level := 1
	for {
		if rand.Int31n(SkipListBranch) == 0 {
			break
		}
		level++
		if level >= SkipListMaxLevel {
			break
		}
	}

	return level
}
