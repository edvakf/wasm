package wasm

import (
	"bytes"
	"testing"
)

func TestVaruint7(t *testing.T) {
	inputs := [][]byte{
		[]byte{'\x00'},
		[]byte{'\x01'},
		[]byte{'\x7F'},
	}

	for i, in := range inputs {
		t.Logf("%d-th example", i)
		readWriteVaruint7(t, in)
	}
}

func readWriteVaruint7(t *testing.T, in []byte) {
	r := bytes.NewBuffer(in)
	var v varuint7
	(&v).read(r)

	w := new(bytes.Buffer)
	v.write(w)
	out := w.Bytes()
	if !bytes.Equal(in, out) {
		t.Errorf("in: %x, out: %x", in, out)
	}
}

func TestVaruint32(t *testing.T) {
	inputs := [][]byte{
		[]byte{'\x00'},
		[]byte{'\x01'},
		[]byte{'\x7F'},
		[]byte{'\x80', '\x01'},
		[]byte{'\x88', '\x7F'},
		[]byte{'\xFF', '\xFF', '\xFF', '\xFF', '\x0F'}, // max of varuint32
		// []byte{'\xFF', '\xFF', '\xFF', '\xFF', '\x10'}, // out of range
	}

	for i, in := range inputs {
		t.Logf("%d-th example", i)
		readWriteVaruint32(t, in)
	}
}

func readWriteVaruint32(t *testing.T, in []byte) {
	r := bytes.NewBuffer(in)
	var v varuint32
	(&v).read(r)

	w := new(bytes.Buffer)
	n, err := v.write(w)
	if err != nil {
		t.Error(err)
	}
	t.Logf("%d bytes written", n)
	out := w.Bytes()
	if !bytes.Equal(in, out) {
		t.Errorf("in: %x, out: %x", in, out)
	}
}
