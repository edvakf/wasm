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
		[]byte{0x00, 'a', 's', 'm', 0x01, 0x00, 0x00, 0x00},
		// simplest wasm with a function
		// (module (func))
		[]byte{
			0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00,
			0x01, 0x04, 0x01, 0x60, 0x00, 0x00, 0x03, 0x02,
			0x01, 0x00, 0x0a, 0x04, 0x01, 0x02, 0x00, 0x0b,
		},
		// (module
		//  (func $add (export "add") (param $lhs i32) (param $rhs i32) (result i32)
		//   get_local $lhs
		//   get_local $rhs
		//   i32.add)
		// )
		[]byte{
			0x00, 0x61, 0x73, 0x6d, 0x01, 0x00, 0x00, 0x00,
			0x01, 0x07, 0x01, 0x60, 0x02, 0x7f, 0x7f, 0x01,
			0x7f, 0x03, 0x02, 0x01, 0x00, 0x07, 0x07, 0x01,
			0x03, 0x61, 0x64, 0x64, 0x00, 0x00, 0x0a, 0x09,
			0x01, 0x07, 0x00, 0x20, 0x00, 0x20, 0x01, 0x6a, 0x0b,
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
