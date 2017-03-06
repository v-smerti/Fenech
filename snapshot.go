package fenech

import (
	"bytes"
	"encoding/base64"
	"errors"
	"os"
	"strconv"
	"sync"
)

var snap = new(sync.Mutex)

//Encode Tuple encode to []byte
func (t *Tuple) Encode() []byte {
	key := base64.URLEncoding.EncodeToString([]byte(t.Key))
	val := base64.URLEncoding.EncodeToString(t.Val)
	return []byte(key + " " + val + "\n")
}

//Decode []byte to Tuple
func (t *Tuple) Decode(b []byte) error {
	d := bytes.Split(b, []byte(" "))
	if len(d) == 2 {
		key, err := base64.URLEncoding.DecodeString(string(d[0]))
		if err != nil {
			return errors.New("*Tuple.Decode: DecodeKey: base64.DecodeString: " + err.Error())
		}
		val, err := base64.URLEncoding.DecodeString(string(d[1]))
		if err != nil {
			return errors.New("*Tuple.Decode: DecodeVal: base64.DecodeString: " + err.Error())
		}
		t.Key = string(key)
		t.Val = val
		return nil
	} else {
		return errors.New("*Tuple.Decode: len(slice) != 2; len(slice) == " + strconv.Itoa(len(d)))
	}

}

func (f *Fenech) restoreSnapshot() error {
	file, err := os.OpenFile(f.dir+"/snapshot.binlog", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()
	if err != nil {
		return err
	}

	s, err := file.Stat()
	if err != nil {
		return err
	}
	if s.Size() > 0 {
		m := make([]byte, s.Size())

		if _, err := file.Read(m); err != nil {
			return err
		}
		data := new(Tuple)
		for _, val := range bytes.Split(m, []byte("\n")) {
			if len(val) > 0 {

				if err := data.Decode(val); err != nil {
					return err
				}
				shard := f.getShard(data.Key)
				shard.Lock()
				shard.items[data.Key] = data.Val
				shard.Unlock()
			}
		}
	}
	return nil

}

func (f *Fenech) upSnapshot() error {
	snap.Lock()
	defer snap.Unlock()
	file, err := os.OpenFile(f.dir+"/snapshot.binlog", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	defer file.Close()
	if err != nil {
		return err
	}
	if err := file.Truncate(0); err != nil {
		return err
	}
	items, err := f.IterBuffered()
	if err != nil {
		return err
	}
	for item := range items {
		if _, err := file.Write(item.Encode()); err != nil {
			return err
		}
	}
	return nil
}

func (f *Fenech) clearSnapshot() error {
	snap.Lock()

	file, err := os.OpenFile(f.dir+"/snapshot.binlog", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	defer func() {
		snap.Unlock()
		file.Close()
	}()

	if err != nil {
		return err
	}
	return file.Truncate(0)
}
