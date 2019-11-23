package main

import (
	"./util"
	"fmt"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	m  :=util.New("3s")
	c:= "test"

	m.Add( c)

	fmt.Println(m.Test( c))
	time.Sleep(time.Duration(5)*time.Second)
	m.Reset()
	fmt.Println(m.Test( c))

}
