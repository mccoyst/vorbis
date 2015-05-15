// © 2015 Steve McCoy under the MIT license. See LICENSE for details.

package vorbis

import (
	"bytes"
	"testing"
)

func TestBitReader_ReadUint(t *testing.T) {
	bits := []byte{ 1 }
	r := NewBitReader(bytes.NewReader(bits))
	n, err := r.ReadUint(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("wanted 1: %v", n)
	}

	bits = []byte{ 2 }
	r = NewBitReader(bytes.NewReader(bits))
	n, err = r.ReadUint(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("wanted 0: %v", n)
	}
	n, err = r.ReadUint(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("wanted 1: %v", n)
	}

	bits = []byte{ 3 }
	r = NewBitReader(bytes.NewReader(bits))
	n, err = r.ReadUint(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("wanted 1: %v", n)
	}
	n, err = r.ReadUint(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != 1 {
		t.Errorf("wanted 1: %v", n)
	}
}

func TestBitReader_ReadInt(t *testing.T) {
	bits := []byte{ 1 }
	r := NewBitReader(bytes.NewReader(bits))
	n, err := r.ReadInt(1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != -1 {
		t.Errorf("wanted -1: %v", n)
	}

	bits = []byte{ 2 }
	r = NewBitReader(bytes.NewReader(bits))
	n, err = r.ReadInt(2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != -2 {
		t.Errorf("wanted -2: %d", n)
	}

	bits = []byte{ 3 }
	r = NewBitReader(bytes.NewReader(bits))
	n, err = r.ReadInt(2)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if n != -1 {
		t.Errorf("wanted 1: %v", n)
	}
}
