package skiplist

import "sync"

type Interface interface {
	Less(other interface{}) bool
}

type SkipList struct {
	header     *Element
	tail       *Element
	update     []*Element
	rank       []int
	length     int
	level      int
	elementMap map[interface{}]*Element
	mutex      sync.RWMutex
}

// New returns an initialized skiplist.
func New() *SkipList {
	return &SkipList{
		header:     newElement(SkipListMaxLevel, 0, nil),
		tail:       nil,
		update:     make([]*Element, SkipListMaxLevel),
		rank:       make([]int, SkipListMaxLevel),
		length:     0,
		level:      1,
		elementMap: make(map[interface{}]*Element),
	}
}

// Init initializes or clears skiplist sl.
func (sl *SkipList) Init() *SkipList {
	sl.header = newElement(SkipListMaxLevel, 0, nil)
	sl.tail = nil
	sl.update = make([]*Element, SkipListMaxLevel)
	sl.rank = make([]int, SkipListMaxLevel)
	sl.length = 0
	sl.level = 1
	sl.elementMap = make(map[interface{}]*Element)

	return sl
}

// Front returns the first elements of skiplist sl or nil.
func (sl *SkipList) Front() *Element {
	return sl.header.level[0].forward
}

// Back returns the last elements of skiplist sl or nil.
func (sl *SkipList) Back() *Element {
	sl.mutex.RLock()
	defer sl.mutex.RUnlock()

	return sl.tail
}

// Len returns the numbler of elements of skiplist sl.
func (sl *SkipList) Len() int {
	sl.mutex.RLock()
	defer sl.mutex.RUnlock()

	return sl.length
}

// Insert inserts v, increments sl.length, and returns a new element of wrap v.
func (sl *SkipList) Set(k interface{}, v Interface) *Element {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()
	sl.remove(k)
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		// store rank that is crossed to reach the insert position
		if i == sl.level-1 {
			sl.rank[i] = 0
		} else {
			sl.rank[i] = sl.rank[i+1]
		}
		for x.level[i].forward != nil && x.level[i].forward.Value.Less(v) {
			sl.rank[i] += x.level[i].span
			x = x.level[i].forward
		}
		sl.update[i] = x
	}

	// ensure that the v is unique, the re-insertion of v should never happen since the
	// caller of sl.Insert() should test in the hash table if the element is already inside or not.
	level := randomLevel()
	if level > sl.level {
		for i := sl.level; i < level; i++ {
			sl.rank[i] = 0
			sl.update[i] = sl.header
			sl.update[i].level[i].span = sl.length
		}
		sl.level = level
	}

	x = newElement(level, k, v)
	for i := 0; i < level; i++ {
		x.level[i].forward = sl.update[i].level[i].forward
		sl.update[i].level[i].forward = x

		// update span covered by update[i] as x is inserted here
		x.level[i].span = sl.update[i].level[i].span - sl.rank[0] + sl.rank[i]
		sl.update[i].level[i].span = sl.rank[0] - sl.rank[i] + 1
	}

	// increment span for untouched levels
	for i := level; i < sl.level; i++ {
		sl.update[i].level[i].span++
	}

	if sl.update[0] == sl.header {
		x.backward = nil
	} else {
		x.backward = sl.update[0]
	}
	if x.level[0].forward != nil {
		x.level[0].forward.backward = x
	} else {
		sl.tail = x
	}

	sl.length++
	sl.elementMap[k] = x

	return x
}

// deleteElement deletes e from its skiplist, and decrements sl.length.
func (sl *SkipList) deleteElement(e *Element, update []*Element) {
	for i := 0; i < sl.level; i++ {
		if update[i].level[i].forward == e {
			update[i].level[i].span += e.level[i].span - 1
			update[i].level[i].forward = e.level[i].forward
		} else {
			update[i].level[i].span--
		}
	}

	if e.level[0].forward != nil {
		e.level[0].forward.backward = e.backward
	} else {
		sl.tail = e.backward
	}

	for sl.level > 1 && sl.header.level[sl.level-1].forward == nil {
		sl.level--
	}

	sl.length--
	delete(sl.elementMap, e.key)
}

func (sl *SkipList) remove(key interface{}) interface{} {
	e := sl.getElement(key)
	if e != nil {
		return sl.removeByElement(e)
	}

	return nil
}

// Remove removes e from sl if e is an element of skiplist sl.
// It returns the element value e.Value.
func (sl *SkipList) Remove(key interface{}) interface{} {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	return sl.remove(key)
}

func (sl *SkipList) removeByElement(e *Element) interface{} {
	x := sl.find(e.Value)                 // x.Value >= e.Value
	if x == e && !e.Value.Less(x.Value) { // e.Value >= x.Value
		sl.deleteElement(x, sl.update)

		return x.Value
	}

	return nil
}

// Remove removes e from sl if e is an element of skiplist sl.
// It returns the element value e.Value.
func (sl *SkipList) RemoveByElement(e *Element) interface{} {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()

	return sl.removeByElement(e)
}

// Delete deletes an element e that e.Value == v, and returns e.Value or nil.
func (sl *SkipList) RemoveByData(v Interface) interface{} {
	sl.mutex.Lock()
	defer sl.mutex.Unlock()
	x := sl.find(v)                   // x.Value >= v
	if x != nil && !v.Less(x.Value) { // v >= x.Value
		sl.deleteElement(x, sl.update)

		return x.Value
	}

	return nil
}

// Find finds an element e that e.Value == v, and returns e or nil.
func (sl *SkipList) Find(v Interface) *Element {
	sl.mutex.RLock()
	defer sl.mutex.RUnlock()
	x := sl.find(v)                   // x.Value >= v
	if x != nil && !v.Less(x.Value) { // v >= x.Value
		return x
	}

	return nil
}

// find finds the first element e that e.Value >= v, and returns e or nil.
func (sl *SkipList) find(v Interface) *Element {
	x := sl.header
	for i := sl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.Value.Less(v) {
			x = x.level[i].forward
		}
		sl.update[i] = x
	}

	return x.level[0].forward
}

// get data by key
func (sl *SkipList) Get(key interface{}) Interface {
	sl.mutex.RLock()
	defer sl.mutex.RUnlock()
	if elem, ok := sl.elementMap[key]; ok && elem != nil {
		return elem.Value
	}

	return nil
}

func (sl *SkipList) getElement(key interface{}) *Element {
	if elem, ok := sl.elementMap[key]; ok && elem != nil {
		return elem
	}

	return nil
}

// GetRank finds then rank for key
// O(lgn)
func (sl *SkipList) GetRank(key interface{}) int {
	sl.mutex.RLock()
	defer sl.mutex.RUnlock()
	elem := sl.getElement(key)
	if elem != nil {
		return sl.getRankByData(elem.Value)
	}

	return 0
}

func (sl *SkipList) getRankByData(v Interface) int {
	x := sl.header
	rank := 0
	for i := sl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && x.level[i].forward.Value.Less(v) {
			rank += x.level[i].span
			x = x.level[i].forward
		}
		if x.level[i].forward != nil && !x.level[i].forward.Value.Less(v) && !v.Less(x.level[i].forward.Value) {
			rank += x.level[i].span

			return rank
		}
	}

	return 0
}

// GetRankByData finds the rank for an element e that e.Value == v,
// Returns 0 when the element cannot be found, rank otherwise.
// Note that the rank is 1-based due to the span of sl.header to the first element.
// O(lgn)
func (sl *SkipList) GetRankByData(v Interface) int {
	sl.mutex.RLock()
	defer sl.mutex.RUnlock()

	return sl.getRankByData(v)
}

// GetDataByRank finds an element by ites rank. The rank argument needs bo be 1-based.
// Note that is the first element e that GetRank(e.Value) == rank, and returns e or nil.
// O(lgn)
func (sl *SkipList) GetElementByRank(rank int) *Element {
	sl.mutex.RLock()
	defer sl.mutex.RUnlock()
	x := sl.header
	traversed := 0
	for i := sl.level - 1; i >= 0; i-- {
		for x.level[i].forward != nil && traversed+x.level[i].span <= rank {
			traversed += x.level[i].span
			x = x.level[i].forward
		}
		if traversed == rank {
			return x
		}
	}

	return nil
}
