package kits

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
)

// MultipleReader 实现多重读取器
type MultipleReader interface {
	// Reader 获取一个可关闭实例
	Reader() io.ReadCloser
}

//多重读取器，返回多个reader
type gureMultipleReader struct {
	data []byte
}

func (g *gureMultipleReader) Reader() io.ReadCloser {
	//TODO implement me
	return ioutil.NopCloser(bytes.NewReader(g.data))
}

// NewMultipleReader 传入reader进行封装
func NewMultipleReader(reader io.Reader) (MultipleReader, error) {
	var data []byte
	var err error
	if reader != nil {
		data, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("create multipleReader fail with %v", err)
		}
	} else {
		data = []byte{}
	}
	return &gureMultipleReader{data: data}, nil
}
