package dict

type SimpleDict struct {
	m map[string]interface{}
}

func MakeSimpleDick() *SimpleDict {
	return &SimpleDict{
		m: make(map[string]interface{}),
	}
}

func (d *SimpleDict) Get(key string) (val interface{}, exists bool) {
	val, ok := d.m[key]
	return val, ok
}

func (d *SimpleDict) Len() int {
	if d.m == nil {
		panic("dict is nil")
	}
	return len(d.m)

}
func (d *SimpleDict) Put(key string, val interface{}) (result int) {
	_, exist := d.m[key]
	d.m[key] = val
	if exist {
		return 0
	}
	return 1
}

func (d *SimpleDict) PutIfAbsent(key string, val interface{}) (result int) {
	_, exist := d.m[key]
	if exist {
		return 0
	}
	d.m[key] = val
	return 1
}
func (d *SimpleDict) PutIfExists(key string, val interface{}) (result int) {
	_, exist := d.m[key]
	if exist {
		d.m[key] = val
		return 1
	}
	return 0
}
func (d *SimpleDict) Remove(key string) (result int) {
	_, exist := d.m[key]
	delete(d.m, key)
	if exist {
		return 1
	}
	return 0
}
func (d *SimpleDict) ForEach(consumer Consumer) {
	for k, v := range d.m {
		if !consumer(k, v) {
			break
		}
	}
}
func (d *SimpleDict) Keys() []string {

	keys := make([]string, 0, len(d.m))
	for k := range d.m {
		keys = append(keys, k)
	}
	return keys

}

//allow duplicate keys
func (d *SimpleDict) RandomKeys(limit int) []string {
	result := make([]string, limit)
	i := 0
	for len(result) < limit {

		for k := range d.m {
			if i == limit {
				break
			}
			result = append(result, k)
			i++
		}

	}

	return result

}
func (d *SimpleDict) RandomDistinctKeys(limit int) []string {
	size := limit
	if size > len(d.m) {
		size = len(d.m)
	}
	result := make([]string, size)
	i := 0
	for k := range d.m {
		if i == size {
			break
		}
		result[i] = k
		i++
	}
	return result

}
func (d *SimpleDict) Clear() {
	*d = *MakeSimpleDick()
}
