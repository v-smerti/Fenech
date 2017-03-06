package fenech

import "testing"
import "strconv"

var DB *Fenech

func init() {
	M, err := New("./test")
	if err != nil {
		panic(err)
	}
	DB = M

}
func BenchmarkItems(b *testing.B) {

	// Insert 100 elements.
	for i := 0; i < 10000; i++ {
		DB.Set(strconv.Itoa(i), []byte(strconv.Itoa(i)))
	}
	for i := 0; i < b.N; i++ {
		DB.Items()
	}
}

func BenchmarkStrconv(b *testing.B) {
	for i := 0; i < b.N; i++ {
		strconv.Itoa(i)
	}
}

func BenchmarkSingleInsertAbsent(b *testing.B) {

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DB.Set(strconv.Itoa(i), []byte("value"))
	}
}

func BenchmarkSingleInsertPresent(b *testing.B) {

	DB.Set("key", []byte("value"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DB.Set("key", []byte("value"))
	}
}

func BenchmarkMultiInsertDifferent(b *testing.B) {

	finished := make(chan struct{}, b.N)
	_, set := GetSet(DB, finished)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set(strconv.Itoa(i), "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiInsertSame(b *testing.B) {

	finished := make(chan struct{}, b.N)
	_, set := GetSet(DB, finished)
	DB.Set("key", []byte("value"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set("key", "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSame(b *testing.B) {

	finished := make(chan struct{}, b.N)
	get, _ := GetSet(DB, finished)
	DB.Set("key", []byte("value"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		get("key", "value")
	}
	for i := 0; i < b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetDifferent(b *testing.B) {

	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(DB, finished)
	DB.Set("-1", []byte("value"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set(strconv.Itoa(i-1), "value")
		get(strconv.Itoa(i), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func BenchmarkMultiGetSetBlock(b *testing.B) {

	finished := make(chan struct{}, 2*b.N)
	get, set := GetSet(DB, finished)
	for i := 0; i < b.N; i++ {
		DB.Set(strconv.Itoa(i%100), []byte("value"))
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		set(strconv.Itoa(i%100), "value")
		get(strconv.Itoa(i%100), "value")
	}
	for i := 0; i < 2*b.N; i++ {
		<-finished
	}
}

func GetSet(m *Fenech, finished chan struct{}) (set func(key, value string), get func(key, value string)) {
	return func(key, value string) {
			for i := 0; i < 10; i++ {
				m.Get(key)
			}
			finished <- struct{}{}
		}, func(key, value string) {
			for i := 0; i < 10; i++ {
				m.Set(key, []byte(value))
			}
			finished <- struct{}{}
		}
}

func BenchmarkKeys(b *testing.B) {

	// Insert 100 elements.
	for i := 0; i < 10000; i++ {
		DB.Set(strconv.Itoa(i), []byte(strconv.Itoa(i)))
	}
	for i := 0; i < b.N; i++ {
		DB.Keys()
	}
}
