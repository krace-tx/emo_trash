package test

import "testing"

type Interface interface{}

func TestType(t *testing.T) {
	var i Interface
	i = make([]int, 0)
	i = append(i.([]int), 1)
	i = append(i.([]int), 2)
	i = append(i.([]int), 3)
	t.Log(i.([]int)[:0:1])
	t.Log(i)
	t.Logf("%T", i)
	i = "1"
	t.Log(i)
	t.Logf("%T", i)
	i = true
	t.Log(i)
	t.Logf("%T", i)
}
