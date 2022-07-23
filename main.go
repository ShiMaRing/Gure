package main

import (
	"fmt"
	"strings"
)

type A interface {
	hello()
}

type C interface {
	A
	hi()
}

type D interface {
	A
	gi()
}
type as struct {
}

func (a as) hello() {
}

type cs struct {
}

func (c cs) hello() {
}

func (c cs) hi() {
}

type ds struct {
}

func (d ds) hello() {
	//TODO implement me
	panic("implement me")
}

func (d ds) gi() {
	//TODO implement me
	panic("implement me")
}

func main() {
	var mid = "hello|1929|"
	split := strings.Split(mid, "|")
	fmt.Println(len(split))
	var c C = cs{}
	var d D = ds{}
	fmt.Println(CheckType(c))
	fmt.Println(CheckType(d))
}

func CheckType(a A) string {
	switch a.(type) {
	case C:
		return "c"
	case D:
		return "d"
	default:
		return "unknown"
	}
}
