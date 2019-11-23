package util

import (
	"github.com/emirpasic/gods/maps/hashmap"
	"time"
)

type ExpiredMap struct {
	Map      *hashmap.Map
	Duration time.Duration
}

func New(durationFormat string) *ExpiredMap {
	var mm, _ = time.ParseDuration(durationFormat)
	return &ExpiredMap{Map: hashmap.New(), Duration: mm}
}

func (m ExpiredMap) Test(data string) bool {
	value, found := m.Map.Get(data)
	if !found {
		return false
	}
	before := time.Now().Before(value.(time.Time))
	if !before {
		m.Map.Remove(data)
	}
	return before
}

func (m ExpiredMap) Add(data string) {
	m.Map.Put(data, time.Now().Add(m.Duration))
}

func (m ExpiredMap) Reset() {
	now := time.Now()
	for _, v := range m.Map.Keys() {
		value, _ := m.Map.Get(v)
		if now.After(value.(time.Time)) {
			m.Map.Remove(v)
		}
	}
}
