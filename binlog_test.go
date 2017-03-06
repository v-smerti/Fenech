package fenech

import (
	"testing"
)

func TestEncode(t *testing.T) {
	b := new(binlog)
	b.Type = 0
	b.Key = "kev"
	b.Value = []byte("value")
	s := string(b.Encode())
	t.Log(s)
	// output:  0 a2V2 dmFsdWU=
}

func TestDecode(t *testing.T) {
	bN := new(binlog)
	bN.Type = 0
	bN.Key = "key"
	bN.Value = []byte("value")
	s := bN.Encode()

	b := new(binlog)
	if err := b.Decode(s); err != nil {
		t.Error("Failed decode")
	}
	if b.Type != bN.Type {
		t.Error("Incorrect decode Type")
	}
	if b.Key != bN.Key {
		t.Error("Incorrect decode Key")
	}
	if string(b.Value) != string(bN.Value) {
		t.Error("Incorrect decode Value")
	}

	if err := b.Decode([]byte("q a2V2 dmFsdWU=")); err == nil {
		t.Error("decode Type: missed error")
	} else {
		if err.Error() != `*binlog.Decode: ParseUint: strconv.ParseUint: parsing "q": invalid syntax` {
			t.Error("wrong to recognize an error")
		}
	}

	if err := b.Decode([]byte("0 __ dmFsdWU=")); err == nil {
		t.Error("decode Type: missed error")
	} else {
		if err.Error() != `*binlog.Decode: DecodeKey: base64.DecodeString: illegal base64 data at input byte 0` {
			t.Error("wrong to recognize an error")
		}
	}

	if err := b.Decode([]byte("0 a2V2 qq")); err == nil {
		t.Error("decode Type: missed error")
	} else {
		if err.Error() != `*binlog.Decode: DecodeVal: base64.DecodeString: illegal base64 data at input byte 0` {
			t.Error("wrong to recognize an error")
		}
	}

	if err := b.Decode([]byte("qweeeeee")); err == nil {
		t.Error("decode Type: missed error")
	} else {
		if err.Error() != `*binlog.Decode: len(slice) != 3; len(slice) == 1` {
			t.Error("wrong to recognize an error")
		}
	}

}
