package fenech

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

func (f *Fenech) restoreBinlog(file string) {

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
	b.Type = Put
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
