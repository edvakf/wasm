package wasm

import (
	"bytes"
	"testing"
)

func TestVaruint7(t *testing.T) {
	inputs := map[varuint7][]byte{
		0:   []byte{0x00},
		1:   []byte{0x01},
		127: []byte{0x7F},
	}

	for n, in := range inputs {
		readWriteVaruint7(t, n, in)
	}
}

func readWriteVaruint7(t *testing.T, n varuint7, in []byte) {
	r := bytes.NewBuffer(in)
	var v varuint7
	(&v).read(r)

	if n != v {
		t.Errorf("expected: %d, actual: %d", n, v)
	}
	w := new(bytes.Buffer)
	err := v.write(w)
	if err != nil {
		t.Error(err)
	}
	out := w.Bytes()
	if !bytes.Equal(in, out) {
		t.Errorf("in: %x, out: %x", in, out)
	}
}

func TestVaruint32(t *testing.T) {
	inputs := map[varuint32][]byte{
		0:          []byte{0x00},
		1:          []byte{0x01},
		127:        []byte{0x7F},
		128:        []byte{0x80, 0x01},
		16264:      []byte{0x88, 0x7F},
		4294967295: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x0F}, // max of varuint32
		// []byte{0xFF, 0xFF, 0xFF, 0xFF, 0x10}, // out of range
	}

	for n, in := range inputs {
		readWriteVaruint32(t, n, in)
	}
}

func readWriteVaruint32(t *testing.T, n varuint32, in []byte) {
	r := bytes.NewBuffer(in)
	var v varuint32
	(&v).read(r)

	if n != v {
		t.Errorf("expected: %d, actual: %d", n, v)
	}
	w := new(bytes.Buffer)
	_, err := v.write(w)
	if err != nil {
		t.Error(err)
	}
	out := w.Bytes()
	if !bytes.Equal(in, out) {
		t.Errorf("in: %x, out: %x", in, out)
	}
}
