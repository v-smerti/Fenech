package fenech

import (
	"hash/fnv"
	"sort"
	"strconv"
	"testing"
)

type Animal []byte

func isAnimal(sl1, sl2 Animal) bool {
	for key, value := range sl1 {
		if sl2[key] != value {
			return false
		}
	}
	return true
}

var m *Fenech

func init() {
	M, err := New("./test")
	if err != nil {
		panic(err)
	}
	m = M

}
func TestMapCreation(t *testing.T) {
	removeAllKey()
	if m == nil {
		t.Error("map is null.")
	}
	i, err := m.Count()
	if err != nil {
		t.Error(err)
	}
	if i != 0 {
		t.Error("new map should be empty.")
	}
}

func TestInsert(t *testing.T) {
	elephant := Animal("elephant")
	monkey := Animal("monkey")

	m.Set("elephant", elephant)
	m.Set("monkey", monkey)
	i, err := m.Count()
	if err != nil {
		t.Error(err)
	}
	if i != 2 {
		t.Error("map should contain exactly two elements.")
	}
}

func TestInsertAbsent(t *testing.T) {
	elephant := Animal("elephant")
	monkey := Animal("monkey")

	m.SetIfAbsent("elephant", elephant)
	ok, err := m.SetIfAbsent("elephant", monkey)
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Error("map set a new value even the entry is already present")
	}
}

func TestGet(t *testing.T) {

	// Get a missing element.
	val, ok, err := m.Get("Money")
	if err != nil {
		t.Error(err)
	}

	if ok == true {
		t.Error("ok should be false when item is missing from map.")
	}

	if val != nil {
		t.Error("Missing values should return as null.")
	}

	elephant := Animal("elephant")
	m.Set("elephant", elephant)

	// Retrieve inserted element.

	m1, ok, err := m.Get("elephant")
	if err != nil {
		t.Error(err)
	}
	if ok == false {
		t.Error("ok should be true for item stored within the map.")
	}

	if &elephant == nil {
		t.Error("expecting an element, not null.")
	}

	if !isAnimal(elephant, m1) {
		t.Log(elephant)
		t.Log(m1)
		t.Error("item was modified.")
	}
}

func TestHas(t *testing.T) {
	ok, err := m.Has("Money")
	if err != nil {
		t.Error(err)
	}
	// Get a missing element.
	if ok == true {
		t.Error("element shouldn't exists")
	}

	elephant := Animal("elephant")
	m.Set("elephant", elephant)
	ok, err = m.Has("elephant")
	if err != nil {
		t.Error(err)
	}
	if ok == false {
		t.Error("element exists, expecting Has to return True.")
	}
}

func TestRemove(t *testing.T) {
	removeAllKey()
	monkey := Animal("monkey")
	m.Set("monkey", monkey)

	m.Remove("monkey")
	i, err := m.Count()
	if err != nil {
		t.Error(err)
	}
	if i != 0 {
		t.Error("Expecting count to be zero once item was removed.")
	}

	temp, ok, err := m.Get("monkey")
	if err != nil {
		t.Error(err)
	}
	if ok != false {
		t.Error("Expecting ok to be false for missing items.")
	}

	if temp != nil {
		t.Error("Expecting item to be nil after its removal.")
	}

	// Remove a none existing element.
	if err := m.Remove("noone"); err != nil {
		t.Error(err)
	}
}

func TestPop(t *testing.T) {

	monkey := Animal("monkey")
	m.Set("monkey", monkey)

	ml, exists, err := m.Pop("monkey")
	if err != nil {
		t.Fatal(err)
	}

	if !exists {
		t.Error("Pop didn't find a monkey.")
	}

	_, exists2, err := m.Pop("monkey")
	if err != nil {
		t.Fatal(err)
	}
	if exists2 || !isAnimal(ml, monkey) {
		t.Error("Pop keeps finding monkey")
	}
	i, err := m.Count()
	if err != nil {
		t.Error(err)
	}
	if i != 0 {
		t.Error("Expecting count to be zero once item was Pop'ed.")
	}

	temp, ok, err := m.Get("monkey")
	if err != nil {
		t.Error(err)
	}
	if ok != false {
		t.Error("Expecting ok to be false for missing items.")
	}

	if temp != nil {
		t.Error("Expecting item to be nil after its removal.")
	}
}

func TestCount(t *testing.T) {
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal(strconv.Itoa(i)))
	}
	i, err := m.Count()
	if err != nil {
		t.Error(err)
	}
	if i != 100 {
		t.Error("Expecting 100 element within map.")
	}
}

func TestIsEmpty(t *testing.T) {
	removeAllKey()
	ok, err := m.IsEmpty()
	if err != nil {
		t.Error(err)
	}
	if ok == false {
		t.Error("new map should be empty")
	}

	m.Set("elephant", Animal("elephant"))
	ok, err = m.IsEmpty()
	if err != nil {
		t.Error(err)
	}
	if ok != false {
		t.Error("map shouldn't be empty.")
	}
}

func TestIterator(t *testing.T) {
	removeAllKey()
	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal(strconv.Itoa(i)))
	}

	counter := 0
	iter, err := m.Iter()
	if err != nil {
		t.Error(err)
	}
	// Iterate over elements.
	for item := range iter {
		val := item.Val

		if val == nil {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestBufferedIterator(t *testing.T) {

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal(strconv.Itoa(i)))
	}

	counter := 0
	iter, err := m.IterBuffered()
	if err != nil {
		t.Error(err)
	}
	// Iterate over elements.
	for item := range iter {
		val := item.Val

		if val == nil {
			t.Error("Expecting an object.")
		}
		counter++
	}

	if counter != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestIterCb(t *testing.T) {

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal(strconv.Itoa(i)))
	}

	counter := 0
	// Iterate over elements.
	m.IterCb(func(key string, v []byte) {
		counter++
	})
	if counter != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestItems(t *testing.T) {

	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal(strconv.Itoa(i)))
	}

	items, err := m.Items()
	if err != nil {
		t.Error(err)
	}

	if len(items) != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestConcurrent(t *testing.T) {
	ch := make(chan int)
	const iterations = 1000
	var a [iterations]int

	// Using go routines insert 1000 ints into our map.
	go func() {
		for i := 0; i < iterations/2; i++ {
			// Add item to map.
			m.Set(strconv.Itoa(i), Animal(strconv.Itoa(i)))

			// Retrieve item from map.
			val, _, err := m.Get(strconv.Itoa(i))
			if err != nil {
				t.Error(err)
			}
			// Parse int
			k, err := strconv.Atoi(string(val))
			if err != nil {
				t.Error("Err parse int")
			}
			// Write to channel inserted value.
			ch <- k
		} // Call go routine with current index.
	}()

	go func() {
		for i := iterations / 2; i < iterations; i++ {
			// Add item to map.
			m.Set(strconv.Itoa(i), Animal(strconv.Itoa(i)))

			// Retrieve item from map.
			val, _, err := m.Get(strconv.Itoa(i))
			if err != nil {
				t.Error(err)
			}
			// Parse int
			k, err := strconv.Atoi(string(val))
			if err != nil {
				t.Error("Err parse int")
			}
			// Write to channel inserted value.
			ch <- k
		} // Call go routine with current index.
	}()

	// Wait for all go routines to finish.
	counter := 0
	for elem := range ch {
		a[counter] = elem
		counter++
		if counter == iterations {
			break
		}
	}

	// Sorts array, will make is simpler to verify all inserted values we're returned.
	sort.Ints(a[0:iterations])
	// Make sure map contains 1000 elements.
	i, err := m.Count()
	if err != nil {
		t.Error(err)
	}

	if i != iterations {
		t.Error("Expecting 1000 elements.")
	}

	// Make sure all inserted values we're fetched from map.
	for i := 0; i < iterations; i++ {
		if i != a[i] {
			t.Error("missing value", i)
		}
	}
}

func TestKeys(t *testing.T) {
	removeAllKey()
	// Insert 100 elements.
	for i := 0; i < 100; i++ {
		m.Set(strconv.Itoa(i), Animal(strconv.Itoa(i)))
	}

	keys, err := m.Keys()
	if err != nil {
		t.Error(err)
	}
	if len(keys) != 100 {
		t.Error("We should have counted 100 elements.")
	}
}

func TestMInsert(t *testing.T) {
	removeAllKey()
	animals := map[string][]byte{
		"elephant": Animal("elephant"),
		"monkey":   Animal("monkey"),
	}

	m.MSet(animals)
	i, err := m.Count()
	if err != nil {
		t.Error(err)
	}
	if i != 2 {
		t.Error("map should contain exactly two elements.")
	}
}

func TestFnv32(t *testing.T) {
	key := []byte("ABC")

	hasher := fnv.New32()
	hasher.Write(key)
	if fnv32(string(key)) != hasher.Sum32() {
		t.Errorf("Bundled fnv32 produced %d, expected result from hash/fnv32 is %d", fnv32(string(key)), hasher.Sum32())
	}
}

func TestUpsert(t *testing.T) {
	removeAllKey()
	dolphin := Animal("dolphin")
	whale := Animal("whale")
	tiger := Animal("tiger")
	lion := Animal("lion")

	cb := func(exists bool, valueInMap []byte, newValue []byte) []byte {
		if !exists {
			return []byte("cb_")
		}
		return append(valueInMap, newValue...)
	}

	m.Set("marine", []byte(dolphin))
	m.Upsert("marine", whale, cb)
	m.Upsert("predator", tiger, cb)
	m.Upsert("predator", lion, cb)
	i, err := m.Count()
	if err != nil {
		t.Error(err)
	}
	if i != 2 {
		t.Error("map should contain exactly two elements.")
	}

	marineAnimals, ok, err := m.Get("marine")
	if err != nil {
		t.Error(err)
	}

	if !ok || string(marineAnimals) != "dolphinwhale" {
		t.Error("Set, then Upsert failed")
	}

	predators, ok, err := m.Get("predator")
	if err != nil {
		t.Error(err)
	}
	if !ok || string(predators) != "cb_lion" {
		t.Error("Upsert, then Upsert failed")
	}

}

func TestKeysWhenRemoving(t *testing.T) {

	// Insert 100 elements.
	Total := 100
	for i := 0; i < Total; i++ {
		m.Set(strconv.Itoa(i), Animal(strconv.Itoa(i)))
	}

	// Remove 10 elements concurrently.
	Num := 10
	for i := 0; i < Num; i++ {
		go func(c *Fenech, n int) {
			c.Remove(strconv.Itoa(n))
		}(m, i)
	}
	keys, err := m.Keys()
	if err != nil {
		t.Error(err)
	}
	for _, k := range keys {
		if k == "" {
			t.Error("Empty keys returned")
		}
	}
}

func removeAllKey() {
	keys, _ := m.Keys()
	for _, key := range keys {
		m.Remove(key)
	}
}
