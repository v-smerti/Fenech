package fenech

import (
	"encoding/json"
	"os"
	"strconv"
	"sync"
)

const SHARD_COUNT = 32

type Fenech struct {
	maps   ConcurrentMap
	binlog ConcurrentBinlog
	sync.Mutex
}
type ConcurrentBinlog []*ConcurrentBinlogShared
type ConcurrentBinlogShared struct {
	file *os.File
	sync.Mutex
}

type ConcurrentMap []*ConcurrentMapShared
type ConcurrentMapShared struct {
	items map[string][]byte
	sync.RWMutex
}

// Creates a new concurrent map.
func New(dir string) (*Fenech, error) {
	m := make(ConcurrentMap, SHARD_COUNT)
	b := make(ConcurrentBinlog, SHARD_COUNT)
	for i := 0; i < SHARD_COUNT; i++ {
		f, err := openFile(dir + "/" + strconv.Itoa(i) + ".binlog")
		if err != nil {
			return new(Fenech), err
		}
		m[i] = &ConcurrentMapShared{items: make(map[string][]byte)}
		b[i] = &ConcurrentBinlogShared{file: f}
	}

	f := new(Fenech)
	f.maps = m
	f.binlog = b
	if err := f.restoreBinlog(); err != nil {
		return new(Fenech), err
	}
	return f, nil
}

func openFile(file string) (*os.File, error) {
	/*if _, err := os.Stat(file); os.IsNotExist(err) {
		if _, err := os.Create(file); err != nil {
			return new(os.File), err
		}
	} else if err != nil {
		return new(os.File), err
	}*/
	//O_RDONLY, O_WRONLY, or O_RDWR
	return os.OpenFile(file, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
}

// Returns shard under given key
func (f *Fenech) getShard(key string) *ConcurrentMapShared {
	return f.maps[getShardId(key)]
}

func getShardId(key string) uint {
	return uint(fnv32(key)) % uint(SHARD_COUNT)
}

func (f *Fenech) MSet(data map[string][]byte) error {

	for key, value := range data {
		if err := putBinlog(f, key, value); err != nil {
			return err
		}
		shard := f.getShard(key)
		shard.Lock()
		shard.items[key] = value
		shard.Unlock()
	}
	return nil
}

// Sets the given value under the specified key.
func (f *Fenech) Set(key string, value []byte) error {
	if err := putBinlog(f, key, value); err != nil {
		return err
	}
	// Get map shard.
	shard := f.getShard(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
	return nil
}

// Callback to return new element to be inserted into the map
// It is called while lock is held, therefore it MUST NOT
// try to access other keys in same map, as it can lead to deadlock since
// Go sync.RWLock is not reentrant
type UpsertCb func(exist bool, valueInMap []byte, newValue []byte) []byte

// Insert or Update - updates existing element or inserts a new one using UpsertCb
func (f *Fenech) Upsert(key string, value []byte, cb UpsertCb) ([]byte, error) {

	shard := f.getShard(key)
	shard.Lock()
	defer shard.Unlock()
	v, ok := shard.items[key]
	res := cb(ok, v, value)
	if err := putBinlog(f, key, res); err != nil {
		return []byte{}, err
	}
	shard.items[key] = res

	return res, nil
}

// Sets the given value under the specified key if no value was associated with it.
func (f *Fenech) SetIfAbsent(key string, value []byte) (bool, error) {
	// Get map shard.
	shard := f.getShard(key)
	shard.Lock()
	_, ok := shard.items[key]
	shard.Unlock()
	if !ok {
		if err := putBinlog(f, key, value); err != nil {
			return false, err
		}
		shard.Lock()
		shard.items[key] = value
		shard.Unlock()
	}

	return !ok, nil
}

// Retrieves an element from map under given key.
func (f *Fenech) Get(key string) ([]byte, bool) {
	// Get shard
	shard := f.getShard(key)
	shard.RLock()
	// Get item from shard.
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

// Returns the number of elements within the map.
func (f *Fenech) Count() int {
	count := 0
	for i := 0; i < SHARD_COUNT; i++ {
		shard := f.maps[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

// Looks up an item under specified key
func (f *Fenech) Has(key string) bool {
	// Get shard
	shard := f.getShard(key)
	shard.RLock()
	// See if element is within shard.
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

// Removes an element from the map.
func (f *Fenech) Remove(key string) error {
	if err := delBinlog(f, key); err != nil {
		return err
	}
	// Try to get shard.
	shard := f.getShard(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
	return nil
}

// Removes an element from the map and returns it
func (f *Fenech) Pop(key string) ([]byte, bool, error) {
	// Try to get shard.
	shard := f.getShard(key)
	shard.Lock()
	v, exists := shard.items[key]
	if exists {
		if err := delBinlog(f, key); err != nil {
			return []byte{}, false, err
		}
		delete(shard.items, key)
	}
	shard.Unlock()
	return v, exists, nil
}

// Checks if map is empty.
func (f *Fenech) IsEmpty() bool {
	return f.Count() == 0
}

// Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type Tuple struct {
	Key string
	Val []byte
}

// Returns an iterator which could be used in a for range loop.
//
// Deprecated: using IterBuffered() will get a better performence
func (f *Fenech) Iter() <-chan Tuple {
	ch := make(chan Tuple)
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(SHARD_COUNT)
		// Foreach shard.
		for _, shard := range f.maps {
			go func(shard *ConcurrentMapShared) {
				// Foreach key, value pair.
				shard.RLock()
				for key, val := range shard.items {
					ch <- Tuple{key, val}
				}
				shard.RUnlock()
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()
	return ch
}

// Returns a buffered iterator which could be used in a for range loop.
func (f *Fenech) IterBuffered() <-chan Tuple {
	ch := make(chan Tuple, f.Count())
	go func() {
		wg := sync.WaitGroup{}
		wg.Add(SHARD_COUNT)
		// Foreach shard.
		for _, shard := range f.maps {
			go func(shard *ConcurrentMapShared) {
				// Foreach key, value pair.
				shard.RLock()
				for key, val := range shard.items {
					ch <- Tuple{key, val}
				}
				shard.RUnlock()
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()
	return ch
}

// Returns all items as map[string]interface{}
func (f *Fenech) Items() map[string][]byte {
	tmp := make(map[string][]byte)

	// Insert items to temporary map.
	for item := range f.IterBuffered() {
		tmp[item.Key] = item.Val
	}

	return tmp
}

// Iterator callback,called for every key,value found in
// maps. RLock is held for all calls for a given shard
// therefore callback sess consistent view of a shard,
// but not across the shards
type IterCb func(key string, v []byte)

// Callback based iterator, cheapest way to read
// all elements in a map.
func (f *Fenech) IterCb(fn IterCb) {
	for idx := range f.maps {
		shard := (f.maps)[idx]
		shard.RLock()
		for key, value := range shard.items {
			fn(key, value)
		}
		shard.RUnlock()
	}
}

// Return all keys as []string
func (f *Fenech) Keys() []string {
	count := f.Count()
	ch := make(chan string, count)
	go func() {
		// Foreach shard.
		wg := sync.WaitGroup{}
		wg.Add(SHARD_COUNT)
		for _, shard := range f.maps {
			go func(shard *ConcurrentMapShared) {
				// Foreach key, value pair.
				shard.RLock()
				for key := range shard.items {
					ch <- key
				}
				shard.RUnlock()
				wg.Done()
			}(shard)
		}
		wg.Wait()
		close(ch)
	}()

	// Generate keys
	keys := make([]string, 0, count)
	for k := range ch {
		keys = append(keys, k)
	}
	return keys
}

//Reviles ConcurrentMap "private" variables to json marshal.
func (f *Fenech) MarshalJSON() ([]byte, error) {
	// Create a temporary map, which will hold all item spread across shards.
	tmp := make(map[string][]byte)

	// Insert items to temporary map.
	for item := range f.IterBuffered() {
		tmp[item.Key] = item.Val
	}
	return json.Marshal(tmp)
}

// Concurrent map uses Interface{} as its value, therefor JSON Unmarshal
// will probably won't know which to type to unmarshal into, in such case
// we'll end up with a value of type map[string]interface{}, In most cases this isn't
// out value type, this is why we've decided to remove this functionality.

func (f *Fenech) UnmarshalJSON(b []byte) (err error) {
	// Reverse process of Marshal.
	tmp := make(map[string][]byte)

	// Unmarshal into a single map.
	if err := json.Unmarshal(b, &tmp); err != nil {
		return nil
	}

	// foreach key,value pair in temporary map insert into our concurrent map.
	for key, val := range tmp {
		f.Set(key, val)
	}
	return nil
}

func fnv32(key string) uint32 {
	hash := uint32(2166136261)
	const prime32 = uint32(16777619)
	for i := 0; i < len(key); i++ {
		hash *= prime32
		hash ^= uint32(key[i])
	}
	return hash
}
