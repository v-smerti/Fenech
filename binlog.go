package fenech

import (
	"bytes"
	"encoding/base64"
	"errors"
	"strconv"
	"sync"
)

const (
	//Put - type operations
	Put uint = 1
	//Delete - type operations
	Delete uint = 2
)

type binlog struct {
	Type     uint
	Key      string
	Value    []byte
	callback func(error)
}

func (b *binlog) Encode() []byte {
	types := strconv.Itoa(int(b.Type))
	key := base64.URLEncoding.EncodeToString([]byte(b.Key))
	val := base64.URLEncoding.EncodeToString(b.Value)
	return []byte(types + " " + key + " " + val + "\n")
}

func (b *binlog) Decode(bs []byte) error {
	d := bytes.Split(bs, []byte(" "))
	if len(d) == 3 {
		types, err := strconv.ParseUint(string(d[0]), 0, 64)
		if err != nil {
			return errors.New("*binlog.Decode: ParseUint: " + err.Error())
		}

		key, err := base64.URLEncoding.DecodeString(string(d[1]))
		if err != nil {
			return errors.New("*binlog.Decode: DecodeKey: base64.DecodeString: " + err.Error())
		}
		val, err := base64.URLEncoding.DecodeString(string(d[2]))
		if err != nil {
			return errors.New("*binlog.Decode: DecodeVal: base64.DecodeString: " + err.Error())
		}

		b.Type = uint(types)
		b.Key = string(key)
		b.Value = val
		return nil
	}
	return errors.New("*binlog.Decode: len(slice) != 3; len(slice) == " + strconv.Itoa(len(d)))

}

func (b *binlog) Commit(f *Fenech) {
	shard := f.binlog[getShardId(b.Key)]
	f.done.Add(1)
	go func() {
		f.done.Done()
		shard.Lock()
		_, err := shard.file.Write(b.Encode())
		shard.Unlock()
		b.callback(err)
	}()
}

func (f *Fenech) clearBinlog() error {
	f.Lock()
	defer f.Unlock()
	for _, blog := range f.binlog {
		if err := blog.file.Truncate(0); err != nil {
			return err
		}
	}
	return nil
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
			if err := binlogDec.Decode(value); err != nil {
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
