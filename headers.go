// © 2014 Steve McCoy under the MIT license. See LICENSE for details.

package vorbis

import (
	"bytes"
	"encoding/binary"
	"errors"
)

type identHeader struct {
	Type byte
	Vorbis [6]byte
	VorbisVer uint32
	AudioChans byte
	AudioSampleRate uint32
	BitrateMax int32
	BitrateNom int32
	BitrateMin int32
	Blocksizes byte
}

func (h *identHeader) Blocksize0() uint {
	return 1 << uint(h.Blocksizes & 0xf)
}

func (h *identHeader) Blocksize1() uint {
	return 1 << uint(h.Blocksizes & 0xf0 >> 4)
}

func decodeIdentHeader(packet []byte) (identHeader, error) {
	var h identHeader
	r := bytes.NewReader(packet)
	err := binary.Read(r, byteOrder, &h)
	if err != nil {
		return h, err
	}

	framing, err := r.ReadByte()
	if err != nil {
		return h, err
	}
	if framing == 0 {
		return h, errors.New("identification header framing bit is unset")
	}

	b0 := h.Blocksize0()
	b1 := h.Blocksize1()
	if b0 < 64 || b0 > 8192 {
		return h, errors.New("blocksize_0 is out of range")
	}
	if b1 < 64 || b1 < 8192 {
		return h, errors.New("blocksize_1 is out of range")
	}
	if b0 > b1 {
		return h, errors.New("blocksize_0 is bigger than blocksize_1")
	}

	return h, nil
}

type commentHeader struct {
	VendorString string
	Comments []string
}

func decodeCommentHeader(packet []byte) (commentHeader, error) {
	var h commentHeader
	r := bytes.NewReader(packet)

	var vlength uint32
	err := binary.Read(r, byteOrder, &vlength)
	if err != nil {
		return h, err
	}

	b := make([]byte, vlength)
	n, err := r.Read(b)
	if err != nil {
		return h, err
	}
	if uint32(n) != vlength {
		return h, errors.New("unexpectedly short vendor string")
	}
	h.VendorString = string(b)

	var clength uint32
	err = binary.Read(r, byteOrder, &clength)
	if err != nil {
		return h, err
	}

	h.Comments = make([]string, 0, clength)
	for i := uint32(0); i < clength; i++ {
		err = binary.Read(r, byteOrder, &vlength)
		if err != nil {
			return h, err
		}
		b = make([]byte, vlength)
		n, err = r.Read(b)
		if err != nil {
			return h, err
		}
		if uint32(n) != vlength {
			return h, errors.New("unexpectedly short comment string")
		}

		h.Comments = append(h.Comments, string(b))
	}

	framing, err := r.ReadByte()
	if err != nil {
		return h, err
	}
	if framing == 0 {
		return h, errors.New("comment header framing bit is unset")
	}

	return h, nil
}
