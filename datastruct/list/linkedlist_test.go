package list

import (
	"strconv"
	"testing"
)

func TestAdd(t *testing.T) {

	list := Make()
	for i := 0; i < 10; i++ {
		list.Add(i)
	}
	list.ForEach(func(i int, v interface{}) bool {
		intVal, _ := v.(int)

		if intVal != i {
			t.Error("add test failed: expected" + strconv.Itoa(i) + ", obtained:" + strconv.Itoa(intVal))
		}
		return true
	})

}

func TestGet(t *testing.T) {
	list := Make()
	for i := 0; i < 10; i++ {
		list.Add(i)
	}

	for i := 0; i < 10; i++ {
		v := list.Get(i)
		intVal, _ := v.(int)

		if intVal != i {
			t.Error("add test failed: expected" + strconv.Itoa(i) + ", obtained:" + strconv.Itoa(intVal))
		}
	}

}
