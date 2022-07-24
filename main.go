package main

import (
	"fmt"
	"sync"
)

type hi struct {
	sync.Map
}

func main() {
	h := hi{}
	h.Store("hello", "world")
	fmt.Println(h.Load("hello"))
}
