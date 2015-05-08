// © 2014 Steve McCoy under the MIT license. See LICENSE for details.

package vorbis

import (
	"encoding/binary"
	"errors"
)

type reader interface {
	Read([]byte) (int, error)
	ReadByte() (byte, error)
}

type commonHeader struct {
	Type byte
	Vorbis [6]byte
}

type identHeader struct {
	commonHeader
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

type ErrBadBlocksize struct {
	Value uint
	Block int
}

func (e ErrBadBlocksize) Error() string {
	return "blocksize is out of range"
}

func decodeIdentHeader(r reader) (identHeader, error) {
	var h identHeader
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
		return h, ErrBadBlocksize{b0, 0}
	}
	if b1 < 64 || b1 > 8192 {
		return h, ErrBadBlocksize{b1, 1}
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

func decodeCommentHeader(r reader) (commentHeader, error) {
	var h commentHeader
	var junk commonHeader
	err := binary.Read(r, byteOrder, &junk)
	if err != nil {
		return h, err
	}

	var vlength uint32
	err = binary.Read(r, byteOrder, &vlength)
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

type setupHeader struct {
	CodebookCount byte
	Codebooks []codebook
	// TimeCount 6 bits
	// uint16 * TimeCount of zeroes
	FloorCount byte
	Floors []floor
	ResidueCount byte
	Residues []residue
	MapCount byte
	Maps []vmap
	ModeCount byte
	Modes []mode
}

type codebook int
type floor int
type residue int
type mode int

type vmap struct {
	Submaps byte
	CouplingSteps byte
	MappingMagnitudes []byte
	MappingAngles []byte
	// 2 bits of zeroes
	MappingMux []byte
	SubmapFloors []byte
	SubmapResidue []byte

}


