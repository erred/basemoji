// Package basemoji implements emoji encoding
// as inspired by bas64 encoding
package basemoji

import (
	"bytes"
	"unicode/utf8"
	// "io"
)

// func NewDecoder(enc *Encoding, r io.Reader) io.Reader
// func NewEncoder(enc *Encoding, w io,Writer) io.Writer

// xxxxxx xx yyyy yyyy zz zzzzzz
// aaaaaa bb bbbb cccc cc dddddd

var (
	// StdEncoding is the default encoding
	// using the first 64 emojis
	StdEncoding = "😀😁😂😃😄😅😆😇😈😉😊😋😌😍😎😏😐😑😒😓😔😕😖😗😘😙😚😛😜😝😞😟😠😡😢😣😤😥😦😧😨😩😪😫😬😭😮😯😰😱😲😳😴😵😶😷😸😹😺😻😼😽😾😿"
)

// Encoding is a radix 64 encoding/decoding scheme,
// defined by a 64-rune alphabet
type Encoding struct {
	rs []rune
	m  map[rune]byte
}

// NewEncoding returns a new unpadded Encoding
// defined by the given alphabet,
// which must be a 64 rune string
func NewEncoding(encoding string) *Encoding {
	enc := &Encoding{
		rs: make([]rune, 64),
		m:  make(map[rune]byte),
	}
	var i int
	for _, r := range encoding {
		enc.rs[i] = r
		enc.m[r] = byte(i)
		i++
	}
	return enc
}

// Decode decodes src using the encoding enc
func (enc *Encoding) Decode(dst, src []byte) (n int, err error) {
	buf, err := enc.decodeBuffer(src)
	if err != nil {
		return 0, err
	}
	return copy(dst, buf.Bytes()), nil
}

// DecodeString returns the bytes represented by the basemoji string s
func (enc *Encoding) DecodeString(s string) ([]byte, error) {
	buf, err := enc.decodeBuffer([]byte(s))
	return buf.Bytes(), err
}

func (enc *Encoding) decodeBuffer(src []byte) (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	rs := make([]byte, 4)
	for len(src) > 3 {
		r, s := utf8.DecodeRune(src)
		src, rs[0] = src[s:], enc.m[r]
		r, s = utf8.DecodeRune(src)
		src, rs[1] = src[s:], enc.m[r]
		buf.WriteByte(rs[0]<<2 | rs[1]>>4)
		if len(src) == 0 {
			break
		}

		r, s = utf8.DecodeRune(src)
		src, rs[2] = src[s:], enc.m[r]
		buf.WriteByte(rs[1]<<4 | rs[2]>>2)
		if len(src) == 0 {
			break
		}

		r, s = utf8.DecodeRune(src)
		src, rs[3] = src[s:], enc.m[r]
		buf.WriteByte(rs[2]<<6 | rs[3])
	}
	return buf, nil
}

// Encode encodes src using the encoding enc
func (enc *Encoding) Encode(dst, src []byte) {
	buf := enc.encodeBuffer(src)
	copy(dst, buf.Bytes())
}

// EncodeToString returns the basemoji encoded src as a string
func (enc *Encoding) EncodeToString(src []byte) string {
	buf := enc.encodeBuffer(src)
	return buf.String()
}

func (enc *Encoding) encodeBuffer(src []byte) *bytes.Buffer {
	buf := &bytes.Buffer{}
	for len(src) > 2 {
		buf.WriteRune(enc.rs[src[0]>>2])
		buf.WriteRune(enc.rs[63&(src[0]<<4|src[1]>>4)])
		buf.WriteRune(enc.rs[63&(src[1]<<2|src[2]>>6)])
		buf.WriteRune(enc.rs[63&src[2]])
		src = src[3:]
	}
	switch len(src) {
	case 1:
		buf.WriteRune(enc.rs[src[0]>>2])
		buf.WriteRune(enc.rs[63&(src[0]<<4)])
	case 2:
		buf.WriteRune(enc.rs[src[0]>>2])
		buf.WriteRune(enc.rs[63&(src[0]<<4|src[1]>>4)])
		buf.WriteRune(enc.rs[63&(src[1]<<2)])
	}
	return buf
}
