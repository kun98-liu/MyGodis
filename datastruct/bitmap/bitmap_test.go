package bitmap

import (
	"math/rand"
	"testing"
)

func TestSetValue(t *testing.T) {
	//生成随机数
	size := 1000
	offsets := make([]int64, size)

	for i := 0; i < size; i++ {
		offsets[i] = int64(i)
	}

	rand.Shuffle(size, func(i, j int) {
		offsets[i], offsets[j] = offsets[j], offsets[i]
	})

	offsets = offsets[0 : size/5]

	//set bit
	offsetMap := make(map[int64]struct{})
	bm := New()
	//保存预期值并设置BitMap
	for _, offset := range offsets {
		offsetMap[offset] = struct{}{}
		bm.SetValue(offset, 1)
	}

	//get Bit
	for i := 0; i < bm.Bitsize(); i++ {
		offset := int64(i)

		_, expectValue := offsetMap[offset]
		actualVal := bm.GetValue(offset) > 0
		if expectValue != actualVal {
			t.Errorf("wrong value at %d", offset)
		}
	}

}
