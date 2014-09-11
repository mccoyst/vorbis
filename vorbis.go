// © 2014 Steve McCoy under the MIT license. See LICENSE for details.

package vorbis

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"github.com/mccoyst/ogg"
)

var byteOrder = binary.LittleEndian

// Decode decodes r as an ogg-vorbis file, returning the interleaved channel data and
// number of channels if the data was not ogg-vorbis.
func Decode(r io.Reader) ([]int16, int, error) {
	var data []int16
	o := ogg.NewDecoder(r)

	p, err := o.Decode()
	if err == io.EOF {
		return nil, 0, errors.New("missing identification header")
	}
	if err != nil {
		return nil, 0, err
	}
	_, err = decodeIdentHeader(bytes.NewReader(p.Packet))
	if err != nil {
		return nil, 0, err
	}

	p, err = o.Decode()
	if err == io.EOF {
		return nil, 0, errors.New("missing identification header")
	}
	if err != nil {
		return nil, 0, err
	}
	_, err = decodeCommentHeader(bytes.NewReader(p.Packet))
	if err != nil {
		return nil, 0, err
	}

	for {
		p, err = o.Decode()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, 0, err
		}

		b := bytes.NewReader(p.Packet)
		var n int16
		for {
			err = binary.Read(b, byteOrder, &n)
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, 0, err
			}
			data = append(data, n)
		}
	}
	return data, 1, nil
}
