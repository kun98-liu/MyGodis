package dict

import (
	"strconv"
	"sync"
	"testing"
)

func TestConcurrentPut(t *testing.T) {
	d := MakeConcurrentDict(0) // default 16

	count := 100
	var wg sync.WaitGroup
	wg.Add(count)
	for i := 0; i < count; i++ {
		go func(i int) {
			key := "k" + strconv.Itoa(i)
			ret := d.Put(key, i)
			if ret != 1 {
				t.Error("Put failed:[" + key + "]" + " -> return val : " + strconv.Itoa(ret))
			}

			val, ok := d.Get(key)
			if ok {
				intVal := val.(int)
				if intVal != i {
					t.Error("Put failed:["+key+"]"+" = "+strconv.Itoa(intVal), " | expected: "+strconv.Itoa(i))
				}
			} else {
				t.Error("Put failed:[" + key + "]" + "does not exist")
			}
			wg.Done()

		}(i)
	}

	wg.Wait()

}
