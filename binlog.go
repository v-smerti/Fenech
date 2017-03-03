package fenech

import (
	"bytes"
	"sync"
)

const (
	Put    uint = 1
	Delete uint = 2
)

//easyjson:json
type binlog struct {
	Type     uint
	Key      string
	Value    []byte
	callback func(error)
}

func (b *binlog) Commit(f *Fenech) {
	shard := f.binlog[getShardId(b.Key)]
	go func() {
		shard.Lock()

		bytes, err := b.MarshalJSON()
		if err != nil {
			shard.Unlock()
			b.callback(err)
			return
		}

		_, err = shard.file.Write(append(bytes, []byte("\n")...))
		shard.Unlock()
		b.callback(err)
	}()
}

func (f *Fenech) restoreBinlog() error {
	wg := new(sync.WaitGroup)
	for key, blog := range f.binlog {
		st, err := blog.file.Stat()
		if err != nil {
			return err
		}
		data := make([]byte, st.Size())
		i, err := blog.file.Read(data)
		if err != nil {
			return err
		}
		if i != 0 {
			wg.Add(1)
			go decodeBinLog(f, key, data, wg)
		}

	}
	wg.Wait()
	return nil
}

func decodeBinLog(f *Fenech, key int, data []byte, wg *sync.WaitGroup) {

	a := bytes.Split(data, []byte("\n"))
	for _, value := range a {
		if len(value) != 0 {
			binlogDec := new(binlog)
			if err := binlogDec.UnmarshalJSON(value); err != nil {
				panic(err.Error())
			}

			shard := f.maps[key]

			switch binlogDec.Type {
			case Put:
				shard.Lock()
				shard.items[binlogDec.Key] = binlogDec.Value
				shard.Unlock()
			case Delete:
				shard.Lock()
				delete(shard.items, binlogDec.Key)
				shard.Unlock()
			default:
				panic("Неизвестный тип операции. ")
			}
		}
	}
	wg.Done()
}

func putBinlog(f *Fenech, key string, value []byte) error {
	b := new(binlog)
	ch := make(chan struct{})
	b.Type = Put
	b.Key = key
	b.Value = value
	var err error
	b.callback = func(e error) {
		err = e
		ch <- struct{}{}
	}
	b.Commit(f)
	<-ch
	return err
}

func delBinlog(f *Fenech, key string) error {
	b := new(binlog)
	ch := make(chan struct{})
	b.Type = Delete
	b.Key = key
	var err error
	b.callback = func(e error) {
		err = e
		ch <- struct{}{}
	}
	b.Commit(f)
	<-ch
	return err
}
