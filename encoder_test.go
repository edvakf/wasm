// Copyright 2016 The wasm Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasm_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	"github.com/sbinet/wasm"
)

func TestBinaries(t *testing.T) {
	binaries := [][]byte{
		// simplest possible wasm
		// (module)
		[]byte{'\x00', 'a', 's', 'm', '\x01', '\x00', '\x00', '\x00'},
		// simplest wasm with a function
		// (module (func))
		[]byte{
			'\x00', '\x61', '\x73', '\x6d', '\x01', '\x00', '\x00', '\x00',
			'\x01', '\x04', '\x01', '\x60', '\x00', '\x00', '\x03', '\x02',
			'\x01', '\x00', '\x0a', '\x04', '\x01', '\x02', '\x00', '\x0b',
		},
	}
	for i, binary := range binaries {
		t.Logf("%d-th binary", i)
		decodeAndEncodeed(t, binary)
	}
}

func decodeAndEncodeed(t *testing.T, in []byte) {
	r := bytes.NewBuffer(in)
	mod, err := wasm.Decode(r)
	if err != nil {
		t.Fatalf("failed to decode: %s", err)
	}
	t.Logf("%#v", mod)

	w := new(bytes.Buffer)
	err = wasm.Encode(*mod, w)
	if err != nil {
		t.Error(err)
	}
	out := w.Bytes()
	t.Log(hex.Dump(out))

	if !bytes.Equal(in, out) {
		t.Errorf("re-encoded binary does not match the original bytes\nin :%x\nout:%x", in, out)
	}
}
