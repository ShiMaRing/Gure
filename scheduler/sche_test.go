package scheduler

import (
	"fmt"
	"testing"
)

func TestDataArgs_Check(t *testing.T) {
	var d = DataArgs{
		ReqBufferCap:      10,
		ReqBufferMaxNum:   010,
		RespBufferCap:     010,
		RespBufferMaxNum:  010,
		ItemBufferCap:     010,
		ItemBufferMaxNum:  010,
		ErrorBufferCap:    010,
		ErrorBufferMaxNum: 010,
	}
	err := d.Check()
	fmt.Println(err)
}

func TestModuleArgs_Check(t *testing.T) {
	var m = ModuleArgs{}
	err := m.Check()
	fmt.Println(err)
}
