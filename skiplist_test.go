package skiplist

import (
	"fmt"
	"math/rand"
	"sort"
	"testing"
	"time"
)

type User struct {
	id    int
	score int
}

// greater to less
func (l User) Less(other interface{}) bool {
	r := other.(*User)
	return l.score > r.score || (l.score == r.score && l.id < r.id)
}

func TestKey(t *testing.T) {
	sl := New()

	users := []*User{
		{id: 1, score: 123},
		{id: 2, score: 1234},
		{id: 3, score: 12345},
		{id: 4, score: 258},
		{id: 5, score: 147},
		{id: 6, score: 369},
		{id: 7, score: 888},
		{id: 8, score: 888},
	}

	length := len(users)
	for _, user := range users {
		sl.Set(user.id, user)
	}

	if sl.Len() != length {
		t.Fatal()
	}

	e := sl.Front()
	for e != nil {
		fmt.Printf("e.score=%d, e.key=%v\n", e.Value.(*User).score, e.key)
		e = e.Next()
	}

	if sl.Get(8).(*User).score != 888 {
		t.Fatal()
	}

	if sl.Get(1).(*User).score != 123 {
		t.Fatal()
	}

	sl.Set(1, &User{id: 1, score: 2555})
	if sl.Len() != length {
		t.Fatal()
	}

	if sl.Get(1).(*User).score != 2555 {
		t.Fatal()
	}

	if sl.Get(4).(*User).score != 258 {
		t.Fatal()
	}

	length--
	sl.Remove(4)
	if sl.Len() != length || sl.Get(4) != nil {
		t.Fatal()
	}

	length--
	sl.Remove(8)
	if sl.Len() != length || sl.Get(8) != nil {
		t.Fatal()
	}

	fmt.Println("------------------------------")
	e = sl.Front()
	for e != nil {
		fmt.Printf("e.score=%d, e.key=%v\n", e.Value.(*User).score, e.key)
		e = e.Next()
	}

	output(sl)
}

// ---------------------------------------------------------------------------------------
type Int int

func (i Int) Less(other interface{}) bool {
	return i < other.(Int)
}

func TestInt(t *testing.T) {
	sl := New()
	if sl.Len() != 0 || sl.Front() != nil && sl.Back() != nil {
		t.Fatal()
	}

	testData := []Int{Int(1), Int(2), Int(3)}

	sl.Set(1, testData[0])
	if sl.Len() != 1 || sl.Front().Value.(Int) != testData[0] || sl.Back().Value.(Int) != testData[0] {
		t.Fatal()
	}

	sl.Set(2, testData[2])
	if sl.Len() != 2 || sl.Front().Value.(Int) != testData[0] || sl.Back().Value.(Int) != testData[2] {
		t.Fatal()
	}

	sl.Set(3, testData[1])
	if sl.Len() != 3 || sl.Front().Value.(Int) != testData[0] || sl.Back().Value.(Int) != testData[2] {
		t.Fatal()
	}

	sl.Set(4, Int(-999))
	sl.Set(5, Int(-888))
	sl.Set(6, Int(888))
	sl.Set(7, Int(999))
	sl.Set(8, Int(1000))

	expect := []Int{Int(-999), Int(-888), Int(1), Int(2), Int(3), Int(888), Int(999), Int(1000)}
	ret := make([]Int, 0)

	for e := sl.Front(); e != nil; e = e.Next() {
		ret = append(ret, e.Value.(Int))
	}
	for i := 0; i < len(ret); i++ {
		if ret[i] != expect[i] {
			t.Fatal()
		}
	}

	e := sl.Find(Int(2))
	if e == nil || e.Value.(Int) != 2 {
		t.Fatal()
	}

	ret = make([]Int, 0)
	for ; e != nil; e = e.Next() {
		ret = append(ret, e.Value.(Int))
	}
	for i := 0; i < len(ret); i++ {
		if ret[i] != expect[i+3] {
			t.Fatal()
		}
	}

	sl.RemoveByElement(sl.Find(Int(2)))
	sl.RemoveByData(Int(888))
	sl.RemoveByData(Int(1000))

	expect = []Int{Int(-999), Int(-888), Int(1), Int(3), Int(999)}
	ret = make([]Int, 0)

	for e := sl.Back(); e != nil; e = e.Prev() {
		ret = append(ret, e.Value.(Int))
	}

	for i := 0; i < len(ret); i++ {
		if ret[i] != expect[len(ret)-i-1] {
			t.Fatal()
		}
	}

	if sl.Front().Value.(Int) != -999 {
		t.Fatal()
	}

	sl.RemoveByElement(sl.Front())
	if sl.Front().Value.(Int) != -888 || sl.Back().Value.(Int) != 999 {
		t.Fatal()
	}

	sl.RemoveByElement(sl.Back())
	if sl.Front().Value.(Int) != -888 || sl.Back().Value.(Int) != 3 {
		t.Fatal()
	}

	if e = sl.Set(100, Int(2)); e.Value.(Int) != 2 {
		t.Fatal()
	}
	sl.RemoveByData(Int(-888))

	if r := sl.RemoveByData(Int(123)); r != nil {
		t.Fatal()
	}

	if sl.Len() != 3 {
		t.Fatal()
	}

	sl.Set(200, Int(2))
	sl.Set(201, Int(2))
	sl.Set(202, Int(1))

	if e = sl.Find(Int(2)); e == nil {
		t.Fatal()
	}

	expect = []Int{Int(2), Int(2), Int(2), Int(3)}
	ret = make([]Int, 0)
	for ; e != nil; e = e.Next() {
		ret = append(ret, e.Value.(Int))
	}
	for i := 0; i < len(ret); i++ {
		if ret[i] != expect[i] {
			t.Fatal()
		}
	}

	sl2 := sl.Init()
	if sl.Len() != 0 || sl.Front() != nil || sl.Back() != nil ||
		sl2.Len() != 0 || sl2.Front() != nil || sl2.Back() != nil {
		t.Fatal()
	}

	// for i := 0; i < 100; i++ {
	// 	sl.Insert(Int(rand.Intn(200)))
	// }
	// output(sl)
}

func TestRank(t *testing.T) {
	t.Parallel()
	sl := New()

	for i := 1; i <= 10; i++ {
		sl.Set(i, Int(i))
	}

	for i := 1; i <= 10; i++ {
		if sl.GetRankByData(Int(i)) != i {
			t.Fatal()
		}
	}

	for i := 1; i <= 10; i++ {
		if sl.GetElementByRank(i).Value != Int(i) {
			t.Fatal()
		}
	}

	if sl.GetRankByData(Int(0)) != 0 || sl.GetRankByData(Int(11)) != 0 {
		t.Fatal()
	}

	if sl.GetElementByRank(11) != nil || sl.GetElementByRank(12) != nil {
		t.Fatal()
	}

	expect := []Int{Int(7), Int(8), Int(9), Int(10)}
	for e, i := sl.GetElementByRank(7), 0; e != nil; e, i = e.Next(), i+1 {
		if e.Value != expect[i] {
			t.Fatal()
		}
	}

	sl = sl.Init()
	mark := make(map[int]bool)
	ss := make([]int, 0)

	for i := 1; i <= 100000; i++ {
		x := rand.Int()
		if !mark[x] {
			mark[x] = true
			sl.Set(i, Int(x))
			ss = append(ss, x)
		}
	}
	sort.Ints(ss)

	for i := 0; i < len(ss); i++ {
		if sl.GetElementByRank(i+1).Value != Int(ss[i]) || sl.GetRankByData(Int(ss[i])) != i+1 {
			t.Fatal()
		}
	}

	// output(sl)
}

func TestGoroutine(t *testing.T) {
	t.Parallel()
	sl := New()
	for range 1000 {
		go testGoroutine(t, sl)
	}

	time.Sleep(3 * time.Second)
}

func testGoroutine(t *testing.T, sl *SkipList) {
	for i := 1; i <= 15; i++ {
		sl.Set(i, Int(i))
	}

	for i := 1; i <= 10; i++ {
		if sl.GetRankByData(Int(i)) != i {
			t.Fatal()
		}
	}

	for i := 1; i <= 10; i++ {
		if sl.GetElementByRank(i).Value != Int(i) {
			t.Fatal()
		}
	}

	sl.Set(1, Int(1))
	res := sl.Get(1)
	if res == nil {
		t.Fatal()
	}
	if int(res.(Int)) != 1 {
		t.Fatal()
	}

	sl.Len()
	sl.Remove(13)
	sl.Remove(99999)
	sl.Back()
	sl.Front()
	sl.Find(Int(666))
	sl.Find(Int(1))
}

func BenchmarkIntInsertOrder(b *testing.B) {
	b.StopTimer()
	sl := New()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sl.Set(i, Int(i))
	}
}

func BenchmarkIntInsertRandom(b *testing.B) {
	b.StopTimer()
	sl := New()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sl.Set(i, Int(rand.Int()))
	}
}

func BenchmarkIntDeleteOrder(b *testing.B) {
	b.StopTimer()
	sl := New()
	for i := 0; i < 1000000; i++ {
		sl.Set(i, Int(i))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sl.RemoveByData(Int(i))
	}
}

func BenchmarkIntDeleteRandome(b *testing.B) {
	b.StopTimer()
	sl := New()
	for i := 0; i < 1000000; i++ {
		sl.Set(i, Int(rand.Int()))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sl.Set(i, Int(rand.Int()))
	}
}

func BenchmarkIntFindOrder(b *testing.B) {
	b.StopTimer()
	sl := New()
	for i := 0; i < 1000000; i++ {
		sl.Set(i, Int(i))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sl.Find(Int(i))
	}
}

func BenchmarkIntFindRandom(b *testing.B) {
	b.StopTimer()
	sl := New()
	for i := 0; i < 1000000; i++ {
		sl.Set(i, Int(rand.Int()))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sl.Find(Int(rand.Int()))
	}
}

func BenchmarkIntRankOrder(b *testing.B) {
	b.StopTimer()
	sl := New()
	for i := 0; i < 1000000; i++ {
		sl.Set(i, Int(i))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sl.GetRank(Int(i))
	}
}

func BenchmarkIntRankRandom(b *testing.B) {
	b.StopTimer()
	sl := New()
	for i := 0; i < 1000000; i++ {
		sl.Set(i, Int(rand.Int()))
	}
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		sl.GetRank(Int(rand.Int()))
	}
}

func output(sl *SkipList) {
	var x *Element
	for i := 0; i < SKIPLIST_MAXLEVEL; i++ {
		fmt.Printf("LEVEL[%v]: ", i)
		count := 0
		x = sl.header.level[i].forward
		for x != nil {
			// fmt.Printf("%v -> ", x.Value)
			count++
			x = x.level[i].forward
		}
		// fmt.Println("NIL")
		fmt.Println("count==", count)
	}
}
