package fenech

import (
	"log"
	"os"
)

func (f *Fenech) restoreSnapshot(dir string) error {
	file, err := os.OpenFile(dir+"/snapshot.binlog", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Println(0)
		return err
	}
	s, err := file.Stat()
	if err != nil {
		log.Println(1)
		return err
	}
	if s.Size() > 0 {

		m := make([]byte, s.Size())

		if _, err := file.Read(m); err != nil {
			return err
		}

		return f.unmarshalJSON(m)
	}
	return nil

}

func (f *Fenech) upSnapshot(dir string) error {
	file, err := os.OpenFile(dir+"/snapshot.binlog", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Println(3)
		return err
	}
	if err := file.Truncate(0); err != nil {
		log.Println(4)
		return err
	}
	if err := f.marshalJSON(file); err != nil {
		log.Println(5)
		return err
	}

	return nil
}
