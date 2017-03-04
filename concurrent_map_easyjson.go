// AUTOGENERATED FILE: easyjson marshaler/unmarshalers.

package fenech

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonFc247620DecodeBitbucketOrgTdmvFenech(in *jlexer.Lexer, out *Tuple) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "Key":
			out.Key = string(in.String())
		case "Val":
			if in.IsNull() {
				in.Skip()
				out.Val = nil
			} else {
				out.Val = in.Bytes()
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonFc247620EncodeBitbucketOrgTdmvFenech(out *jwriter.Writer, in Tuple) {
	out.RawByte('{')
	first := true
	_ = first
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"Key\":")
	out.String(string(in.Key))
	if !first {
		out.RawByte(',')
	}
	first = false
	out.RawString("\"Val\":")
	out.Base64Bytes(in.Val)
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Tuple) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonFc247620EncodeBitbucketOrgTdmvFenech(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Tuple) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonFc247620EncodeBitbucketOrgTdmvFenech(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Tuple) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonFc247620DecodeBitbucketOrgTdmvFenech(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Tuple) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonFc247620DecodeBitbucketOrgTdmvFenech(l, v)
}
func easyjsonFc247620DecodeBitbucketOrgTdmvFenech1(in *jlexer.Lexer, out *Fenech) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonFc247620EncodeBitbucketOrgTdmvFenech1(out *jwriter.Writer, in Fenech) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Fenech) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonFc247620EncodeBitbucketOrgTdmvFenech1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Fenech) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonFc247620EncodeBitbucketOrgTdmvFenech1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Fenech) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonFc247620DecodeBitbucketOrgTdmvFenech1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Fenech) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonFc247620DecodeBitbucketOrgTdmvFenech1(l, v)
}
func easyjsonFc247620DecodeBitbucketOrgTdmvFenech2(in *jlexer.Lexer, out *ConcurrentMapShared) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonFc247620EncodeBitbucketOrgTdmvFenech2(out *jwriter.Writer, in ConcurrentMapShared) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ConcurrentMapShared) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonFc247620EncodeBitbucketOrgTdmvFenech2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ConcurrentMapShared) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonFc247620EncodeBitbucketOrgTdmvFenech2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ConcurrentMapShared) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonFc247620DecodeBitbucketOrgTdmvFenech2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ConcurrentMapShared) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonFc247620DecodeBitbucketOrgTdmvFenech2(l, v)
}
func easyjsonFc247620DecodeBitbucketOrgTdmvFenech3(in *jlexer.Lexer, out *ConcurrentBinlogShared) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonFc247620EncodeBitbucketOrgTdmvFenech3(out *jwriter.Writer, in ConcurrentBinlogShared) {
	out.RawByte('{')
	first := true
	_ = first
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v ConcurrentBinlogShared) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonFc247620EncodeBitbucketOrgTdmvFenech3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v ConcurrentBinlogShared) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonFc247620EncodeBitbucketOrgTdmvFenech3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *ConcurrentBinlogShared) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonFc247620DecodeBitbucketOrgTdmvFenech3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *ConcurrentBinlogShared) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonFc247620DecodeBitbucketOrgTdmvFenech3(l, v)
}