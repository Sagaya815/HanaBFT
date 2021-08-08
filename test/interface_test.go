// -author: Iridescent -time: 2021/8/8
package test

import (
	"fmt"
	"testing"
)

type AA interface {
	call()
}

type A struct {
	name string
	AA
}

type B struct {
	age int
}

func (b *B) call() {
	fmt.Println("age is", b.age)
}

func Test(t *testing.T) {
	b := &B{age: 19}
	a := &A{
		name: "son",
		AA:   b,
	}
	a.call()
}
