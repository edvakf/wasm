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

func TestVarint7(t *testing.T) {
	inputs := map[varint7][]byte{
		0:   []byte{0x00},
		1:   []byte{0x01},
		63:  []byte{0x3F},
		-1:  []byte{0x7F},
		-2:  []byte{0x7E},
		-3:  []byte{0x7D},
		-4:  []byte{0x7C},
		-16: []byte{0x70},
		-32: []byte{0x60},
		-64: []byte{0x40},
	}

	for n, in := range inputs {
		readWriteVarint7(t, n, in)
	}
}

func readWriteVarint7(t *testing.T, n varint7, in []byte) {
	r := bytes.NewBuffer(in)
	v, _, err := varint(r)
	if err != nil {
		t.Error(err)
	}

	if n != varint7(v) {
		t.Errorf("expected: %d, actual: %d", n, v)
	}
	w := new(bytes.Buffer)
	err = varint7(v).write(w)
	if err != nil {
		t.Error(err)
	}
	out := w.Bytes()
	if !bytes.Equal(in, out) {
		t.Errorf("in: %x, out: %x", in, out)
		t.Errorf("%#v", out)
	}
}

// func TestVarint32(t *testing.T) {
// 	inputs := map[varint32][]byte{
// 		0:  []byte{0x00}, // 0000_0000
// 		1:  []byte{0x01}, // 0000_0001
// 		63: []byte{0x3F}, // 0011_1111
// 		-1: []byte{0x7F}, // 1111_1111 =>
// 		-2: []byte{0x7E}, // 1111_1110
// 		// -64: []byte{0x7E},       // 1100_0000 => 0100_0000 (two's complement)
// 		64: []byte{0xC0, 0x00}, // 0100_0000 => 0000_0000 1100_0000
// 		65: []byte{0xC1, 0x00}, // 0100_0001 => 0000_0000 1100_0001
// 		// -65: []byte{},           //
// 	}
//
// 	for n, in := range inputs {
// 		readWriteVarint32(t, n, in)
// 	}
// }
//
// func readWriteVarint32(t *testing.T, n varint32, in []byte) {
// 	r := bytes.NewBuffer(in)
// 	v, _, err := varint(r)
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	if n != varint32(v) {
// 		t.Errorf("expected: %d, actual: %d", n, v)
// 	}
// 	w := new(bytes.Buffer)
// 	_, err = varint32(v).write(w)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	out := w.Bytes()
// 	if !bytes.Equal(in, out) {
// 		t.Errorf("in: %x, out: %x", in, out)
// 	}
// }
