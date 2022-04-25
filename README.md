skiplist
===============

reference from redis [zskiplist](https://github.com/antirez/redis)


Usage
===============

~~~Go

package main

import (
	"fmt"
	"github.com/gansidui/skiplist"
	"log"
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

func main() {
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
		fmt.Println("if sl.Len() != length {")
	}

	e := sl.Front()
	for e != nil {
		fmt.Printf("e.score=%d, e.key=%v\n", e.Value.(*User).score, e.key)
		e = e.Next()
	}

	if sl.Get(8).(*User).score != 888 {
		fmt.Println("sl.Get(8).(*User).score != 888")
	}

	if sl.Get(1).(*User).score != 123 {
		fmt.Println("sl.Get(1).(*User).score != 123")
	}

	sl.Set(1, &User{id: 1, score: 2555})
	if sl.Len() != length {
		fmt.Println("sl.Len() != length")
	}

	if sl.Get(1).(*User).score != 2555 {
		fmt.Println("sl.Get(1).(*User).score != 2555")
	}

	if sl.Get(4).(*User).score != 258 {
		fmt.Println("sl.Get(4).(*User).score != 258")
	}

	length--
	sl.Remove(4)
	if sl.Len() != length || sl.Get(4) != nil {
		fmt.Println("sl.Len() != length || sl.Get(4) != nil")
	}

	length--
	sl.Remove(8)
	if sl.Len() != length || sl.Get(8) != nil {
		fmt.Println("sl.Len() != length || sl.Get(8) != nil")
	}

	fmt.Println("------------------------------")
	e = sl.Front()
	for e != nil {
		fmt.Printf("e.score=%d, e.key=%v\n", e.Value.(*User).score, e.key)
		e = e.Next()
	}
}

~~~


License
===============

MIT