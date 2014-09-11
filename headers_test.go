// © 2014 Steve McCoy under the MIT license. See LICENSE for details.

package vorbis

import (
	"bytes"
	"testing"
)

func TestDecodeIdentHeader(t *testing.T) {
	oh := []byte{
		0x01,
		0x76, 0x6f, 0x72, 0x62, 0x69, 0x73,
		0x00, 0x00, 0x00, 0x00,
		0x02,
		0x44, 0xac, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x03, 0xf4, 0x01, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0xb8,
		0x01,
	}

	r := bytes.NewReader(oh)
	h, err := decodeIdentHeader(r)
	if err != nil {
		if bb, ok := err.(ErrBadBlocksize); ok {
			t.Fatalf("bad blocksize%d: %d", bb.Block, bb.Value)
		} else {
			t.Fatal("unexpected error", err)
		}
	}

	if h.Type != 1 {
		t.Error("wrong header type:", h.Type)
	}
	if bytes.Compare(h.Vorbis[:], []byte("vorbis")) != 0 {
		t.Error("wrong vorbis string:", h.Vorbis)
	}
	if h.VorbisVer != 0 {
		t.Error("wrong vorbis version:", h.VorbisVer)
	}
	if h.AudioChans != 2 {
		t.Error("wrong audio chans:", h.AudioChans)
	}
	if h.AudioSampleRate != 44100 {
		t.Error("wrong audio sample rate:", h.AudioSampleRate)
	}
	if h.BitrateMax != 0 {
		t.Error("wrong bitrate max:", h.BitrateMax)
	}
	if h.BitrateNom != 128003 {
		t.Error("wrong bitrate nom:", h.BitrateNom)
	}
	if h.BitrateMin != 0 {
		t.Error("wrong bitrate min:", h.BitrateMin)
	}
	if h.Blocksize0() != 256 {
		t.Error("wrong blocksize0:", h.BitrateMin)
	}
	if h.Blocksize1() != 2048 {
		t.Error("wrong blocksize1:", h.BitrateMin)
	}
}
