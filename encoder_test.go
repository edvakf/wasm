// Copyright 2016 The wasm Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasm_test

import (
	"bytes"
	"testing"

	"github.com/sbinet/wasm"
)

func TestSimple(t *testing.T) {
	mod, err := wasm.NewModule(1)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte{'\x00', 'a', 's', 'm', '\x01', '\x00', '\x00', '\x00'}

	CompareEncodeed(t, *mod, expected)
}

func TestTypeSection(t *testing.T) {
	mod, err := wasm.NewModule(1)
	if err != nil {
		t.Fatal(err)
	}

	expected := []byte{'\x00', 'a', 's', 'm', '\x01', '\x00', '\x00', '\x00'}

	CompareEncodeed(t, *mod, expected)
}

func CompareEncodeed(t *testing.T, mod wasm.Module, expected []byte) {
	var b bytes.Buffer
	enc := wasm.NewEncoder(&b)
	err := enc.Encode(mod)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(expected, b.Next(b.Len())) {
		t.Errorf("invalid")
	}
}
