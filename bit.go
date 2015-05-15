// © 2015 Steve McCoy under the MIT license. See LICENSE for details.

package vorbis

import (
	"io"
)

type BitReader struct {
	r io.ByteReader
	buf byte
	p byte
}

func NewBitReader(r io.ByteReader) BitReader {
	return BitReader{r, 0, 8}
}

func (b *BitReader) ReadUint(n int) (byte, error) {
	if n > 8 {
		panic("too big n")
	}
	if n == 0 {
		return 0, nil
	}

	// super naive for the sake of having a bitreader
	v := byte(0)
	for i := 0; i < n; i++ {
		bit, err := b.nextBit()
		if err != nil {
			return 0, err
		}
		v |= bit << byte(i)
	}

	return v, nil
}

func (b *BitReader) nextBit() (byte, error) {
	if b.p < 8 {
		p := b.p
		b.p++
		return (b.buf >> p) & 1, nil
	}

	byte, err := b.r.ReadByte()
	if err != nil {
		return 0, err
	}
	b.buf = byte
	b.p = 1
	return b.buf & 1, nil
}
