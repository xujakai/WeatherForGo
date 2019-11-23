package util

import (
	"github.com/pmylund/go-bloom"
)

type IFilter interface {
	Test(data string) bool
	Add(data string)
	Reset()
}

type MyFilter struct {
	filter *bloom.Filter
}

func NewFilter() MyFilter {
	return MyFilter{filter: bloom.New(10000, 0.001)}
}

func (m MyFilter) Test(data string) bool {
	return m.filter.Test([]byte(data))
}

func (m MyFilter) Add(data string) {
	m.filter.Add([]byte(data))
}

func (m MyFilter) Reset() {
	m.filter.Reset()
}


