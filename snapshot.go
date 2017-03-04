package fenech

import (
	"bytes"
	"os"
)

func (f *Fenech) restoreSnapshot(dir string) error {
	file, err := os.OpenFile(dir+"/snapshot.binlog", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
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
				if err := data.UnmarshalJSON(val); err != nil {
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

func (f *Fenech) upSnapshot(dir string) error {
	file, err := os.OpenFile(dir+"/snapshot.binlog", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
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
		b, err := item.MarshalJSON()
		if err != nil {
			return err
		}
		if _, err := file.Write(append(b, []byte("\n")...)); err != nil {
			return err
		}
	}
	return nil
}
