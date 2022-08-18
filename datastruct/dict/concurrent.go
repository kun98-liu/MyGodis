package dict

import (
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
)

type shard struct {
	m     map[string]interface{}
	mutex sync.RWMutex
}

type ConcurrentDict struct {
	table       []*shard
	count       int32
	shardCounts int
}

func MakeConcurrentDict(shardCounts int) *ConcurrentDict {
	shardCounts = computeCapacity(shardCounts)
	table := make([]*shard, shardCounts)
	for i := 0; i < shardCounts; i++ {
		table[i] = &shard{
			m: make(map[string]interface{}, 0),
		}
	}

	d := &ConcurrentDict{
		count:       0,
		table:       table,
		shardCounts: shardCounts,
	}
	return d
}

//make sure the dict contains 2 in power of n shards  ( 17 -> 32 )
func computeCapacity(args int) int {
	if args <= 16 {
		return 16
	}
	//suppose: n = 17
	n := args - 1 // 16 0001 0000
	n |= n >> 1   // 0001 1000
	n |= n >> 2   // 0001 1110
	n |= n >> 4   // 0001 1111
	n |= n >> 8   // 0001 1111
	n |= n >> 16  // 0001 1111

	if n < 0 {
		return math.MaxInt32
	}

	return n + 1 // 0010 0000  -> 32
}

const hashseed = uint32(16777619)

//get hashcode for a key
func hash(key string) uint32 {
	hashcode := uint32(2166136261)
	for i := 0; i < len(key); i++ {
		hashcode *= hashseed
		hashcode ^= uint32(key[i])
	}
	return hashcode
}

//return the shard_id based on hashcode
func (dict *ConcurrentDict) spread(hashcode uint32) uint32 {
	if dict == nil {
		panic("dict is nil")
	}

	shardCounts := uint32(len(dict.table))

	return (shardCounts - 1) & hashcode
}

func (dict *ConcurrentDict) getShard(shard_id uint32) *shard {
	if dict == nil {
		panic("dict is nil")
	}
	return dict.table[shard_id]
}

func (dict *ConcurrentDict) Get(key string) (val interface{}, exists bool) {
	if dict == nil {
		panic("dict is nil")
	}

	hashcode := hash(key)
	shard_id := dict.spread(hashcode)
	shard := dict.getShard(shard_id)

	shard.mutex.RLock()
	defer shard.mutex.RUnlock()

	val, exists = shard.m[key]
	return
}

func (dict *ConcurrentDict) Len() int {
	if dict == nil {
		panic("dict is nil")
	}
	return int(atomic.LoadInt32(&dict.count))
}

func (dict *ConcurrentDict) Put(key string, val interface{}) (result int) {
	if dict == nil {
		panic("dict is nil")
	}

	hashcode := hash(key)
	shard_id := dict.spread(hashcode)
	shard := dict.getShard(shard_id)

	shard.mutex.Lock()
	defer shard.mutex.Unlock()

	_, ok := shard.m[key]
	if ok {
		shard.m[key] = val
		return 0
	} else {
		shard.m[key] = val
		dict.addCount()
		return 1
	}
}

func (dict *ConcurrentDict) PutIfAbsent(key string, val interface{}) (result int) {
	if dict == nil {
		panic("dict is nil")
	}

	hashcode := hash(key)
	shard_id := dict.spread(hashcode)
	shard := dict.getShard(shard_id)

	shard.mutex.Lock()
	defer shard.mutex.Unlock()

	_, ok := shard.m[key]

	//already exists -> do not put the new val
	if ok {
		return 0
	} else {
		shard.m[key] = val
		dict.addCount()
		return 1
	}
}

func (dict *ConcurrentDict) PutIfExists(key string, val interface{}) (result int) {
	if dict == nil {
		panic("dict is nil")
	}

	hashcode := hash(key)
	shard_id := dict.spread(hashcode)
	shard := dict.getShard(shard_id)

	shard.mutex.Lock()
	defer shard.mutex.Unlock()

	_, ok := shard.m[key]

	//already exists -> put the new val
	if ok {
		shard.m[key] = val
		return 1
	} else {
		return 0
	}
}

func (dict *ConcurrentDict) Remove(key string) (result int) {
	if dict == nil {
		panic("dict is nil")
	}

	hashcode := hash(key)
	shard_id := dict.spread(hashcode)
	shard := dict.getShard(shard_id)

	shard.mutex.Lock()
	defer shard.mutex.Unlock()

	_, ok := shard.m[key]

	if ok {
		delete(shard.m, key)
		dict.decreaseCount()
		return 1
	} else {
		return 0
	}
}
func (dict *ConcurrentDict) ForEach(consumer Consumer) {
	if dict == nil {
		panic("dict is nil")
	}

	for _, shard := range dict.table {
		shard.mutex.RLock()

		defer shard.mutex.RUnlock()
		for k, v := range shard.m {
			ok := consumer(k, v)
			if !ok {
				return
			}
		}

	}
}

func (dict *ConcurrentDict) Keys() []string {
	if dict == nil {
		panic("dict is nil")
	}

	keys := make([]string, dict.Len())
	dict.ForEach(func(key string, val interface{}) bool {
		keys = append(keys, key)
		return true
	})

	return keys
}

func (dict *ConcurrentDict) RandomKeys(limit int) []string {
	result := make([]string, limit)
	shards := dict.shardCounts
	for len(result) < limit {
		//random shard_id
		randomShardIndex := uint32(rand.Intn(shards))
		shard := dict.getShard(randomShardIndex)
		if shard == nil {
			panic("shard is nil")
		}

		//get a key from the random shard
		key := shard.randomKey()

		if key != "" {
			result = append(result, key)
		}
	}

	return result
}

func (dict *ConcurrentDict) RandomDistinctKeys(limit int) []string {
	size := dict.Len()
	if limit >= size {
		return dict.Keys()
	}

	shards := dict.shardCounts
	res := make(map[string]bool)

	for len(res) < limit {
		randomShardIndex := uint32(rand.Intn(shards))
		shard := dict.getShard(randomShardIndex)

		if shard == nil {
			panic("shard is nil")
		}

		key := shard.randomKey()
		if key != "" {
			res[key] = true
		}
	}

	result := make([]string, 0, len(res))
	for k := range res {
		result = append(result, k)
	}

	return result
}

func (dict *ConcurrentDict) Clear() {
	*dict = *MakeConcurrentDict(dict.shardCounts)
}

func (dict *ConcurrentDict) addCount() int32 {
	return atomic.AddInt32(&dict.count, 1)
}

func (dict *ConcurrentDict) decreaseCount() int32 {
	return atomic.AddInt32(&dict.count, -1)
}

func (shard *shard) randomKey() string {
	if shard == nil {
		panic("shard is nil")
	}

	shard.mutex.RLock()
	defer shard.mutex.RUnlock()
	for key := range shard.m {
		return key
	}
	return ""
}
