package buffer

import "unicode/utf8"

type UTF8Encoding struct {
}

func (e *UTF8Encoding) Decode(in []byte) []rune {
	res := []rune{}
	for i := 0; i < len(in); {
		r, n := utf8.DecodeRune(in[i:])
		i += n
		res = append(res, r)
	}
	return res
}
