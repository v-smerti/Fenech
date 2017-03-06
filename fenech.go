package fenech

import (
	"errors"
	"os"
	"strconv"
	"sync"
)

const writerCount = 250

type Fenech struct {
	maps   []*concurrentMap
	binlog []*concurrentBinlog
	done   *sync.WaitGroup
	sync.RWMutex
	dir string
	cl  bool
}

type concurrentBinlog struct {
	file *os.File
	sync.Mutex
}

type concurrentMap struct {
	items map[string][]byte
	sync.RWMutex
}

//New creates a new concurrent map.
func New(dir string) (*Fenech, error) {
	m := make([]*concurrentMap, writerCount)
	b := make([]*concurrentBinlog, writerCount)
	for i := 0; i < writerCount; i++ {
		f, err := openFile(dir + "/" + strconv.Itoa(i) + ".binlog")
		if err != nil {
			return new(Fenech), err
		}
		m[i] = &concurrentMap{items: make(map[string][]byte)}
		b[i] = &concurrentBinlog{file: f}
	}

	f := new(Fenech)
	f.maps = m
	f.binlog = b
	f.done = new(sync.WaitGroup)
	f.cl = false
	f.dir = dir
	if err := f.restoreSnapshot(); err != nil {
		return new(Fenech), err
	}
	if err := f.restoreBinlog(); err != nil {
		return new(Fenech), err
	}
	if err := f.upSnapshot(); err != nil {
		return new(Fenech), err
	}
	if err := f.clearBinlog(); err != nil {
		return new(Fenech), err
	}

	return f, nil
}
func (f *Fenech) Wait() {
	f.done.Wait()
}

func (f *Fenech) Close() {
	f.Lock()
	f.cl = true
	f.Unlock()
}

func (f *Fenech) isClose() error {
	f.RLock()
	defer f.RUnlock()
	if f.cl {
		return errors.New("close db")
	} else {
		return nil
	}
}

func openFile(file string) (*os.File, error) {
	return os.OpenFile(file, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
}

// Returns shard under given key
func (f *Fenech) getShard(key string) *concurrentMap {
	return f.maps[getShardId(key)]
}

func getShardId(key string) uint {
	return uint(fnv32(key)) % uint(writerCount)
}

func (f *Fenech) MSet(data map[string][]byte) error {
	if err := f.isClose(); err != nil {
		return err
	}
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
	if err := f.isClose(); err != nil {
		return err
	}
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
	if err := f.isClose(); err != nil {
		return []byte{}, err
	}
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
	if err := f.isClose(); err != nil {
		return false, err
	}
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
func (f *Fenech) Get(key string) ([]byte, bool, error) {
	if err := f.isClose(); err != nil {
		return []byte{}, false, err
	}
	// Get shard
	shard := f.getShard(key)
	shard.RLock()
	// Get item from shard.
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok, nil
}

// Returns the number of elements within the map.
func (f *Fenech) Count() (int, error) {
	if err := f.isClose(); err != nil {
		return 0, err
	}
	count := 0
	for i := 0; i < writerCount; i++ {
		shard := f.maps[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count, nil
}

// Looks up an item under specified key
func (f *Fenech) Has(key string) (bool, error) {
	if err := f.isClose(); err != nil {
		return false, err
	}
	// Get shard
	shard := f.getShard(key)
	shard.RLock()
	// See if element is within shard.
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok, nil
}

// Removes an element from the map.
func (f *Fenech) Remove(key string) error {
	if err := f.isClose(); err != nil {
		return err
	}
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
	if err := f.isClose(); err != nil {
		return []byte{}, false, err
	}
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
func (f *Fenech) IsEmpty() (bool, error) {
	if err := f.isClose(); err != nil {
		return false, err
	}
	if i, err := f.Count(); err != nil {
		return false, err
	} else {
		return i == 0, nil
	}
}

// Used by the Iter & IterBuffered functions to wrap two variables together over a channel,
//easyjson:json
type Tuple struct {
	Key string
	Val []byte
}

// Returns an iterator which could be used in a for range loop.
//
// Deprecated: using IterBuffered() will get a better performence
func (f *Fenech) Iter() (<-chan Tuple, error) {
	ch := make(chan Tuple)
	if err := f.isClose(); err != nil {
		return ch, err
	}
	f.done.Add(1)
	go func() {
		defer f.done.Done()
		wg := sync.WaitGroup{}
		wg.Add(writerCount)
		// Foreach shard.
		for _, shard := range f.maps {
			f.done.Add(1)
			go func(shard *concurrentMap) {
				defer f.done.Done()
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
	return ch, nil
}

// Returns a buffered iterator which could be used in a for range loop.
func (f *Fenech) IterBuffered() (<-chan Tuple, error) {
	count, err := f.Count()
	if err != nil {
		return make(chan Tuple), err
	}
	ch := make(chan Tuple, count)
	f.done.Add(1)
	go func() {
		defer f.done.Done()
		wg := sync.WaitGroup{}
		wg.Add(writerCount)
		// Foreach shard.
		for _, shard := range f.maps {
			f.done.Add(1)
			go func(shard *concurrentMap) {
				defer f.done.Done()
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
	return ch, nil
}

// Returns all items as map[string]interface{}
func (f *Fenech) Items() (map[string][]byte, error) {
	tmp := make(map[string][]byte)
	items, err := f.IterBuffered()
	if err != nil {
		return tmp, err
	}

	// Insert items to temporary map.
	for item := range items {
		tmp[item.Key] = item.Val
	}

	return tmp, nil
}

// Iterator callback,called for every key,value found in
// maps. RLock is held for all calls for a given shard
// therefore callback sess consistent view of a shard,
// but not across the shards
type IterCb func(key string, v []byte)

// Callback based iterator, cheapest way to read
// all elements in a map.
func (f *Fenech) IterCb(fn IterCb) error {
	if err := f.isClose(); err != nil {
		return err
	}
	for idx := range f.maps {
		shard := (f.maps)[idx]
		shard.RLock()
		for key, value := range shard.items {
			fn(key, value)
		}
		shard.RUnlock()
	}
	return nil
}

//Keys return all keys as []string and error
func (f *Fenech) Keys() ([]string, error) {

	count, err := f.Count()
	if err != nil {
		return []string{}, err
	}
	ch := make(chan string, count)
	f.done.Add(1)
	go func() {
		defer f.done.Done()
		// Foreach shard.
		wg := sync.WaitGroup{}
		wg.Add(writerCount)
		for _, shard := range f.maps {
			f.done.Add(1)
			go func(shard *concurrentMap) {
				defer f.done.Done()
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
	return keys, nil
}

//DeleteAllKeys delete all key and clear storage
func (f *Fenech) DeleteAllKeys() error {
	iter, err := f.IterBuffered()
	if err != nil {
		return err
	}
	for item := range iter {
		if err := f.Remove(item.Key); err != nil {
			return err
		}
	}
	if err := f.clearBinlog(); err != nil {
		return err
	}
	return f.clearSnapshot()
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
