// Copyright 2016 The wasm Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasm

import (
	"encoding/binary"
	"io"
)

type Encoder struct {
	w   io.Writer
	err error
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (e *Encoder) Encode(m Module) error {
	e.writeModule(m)

	if e.err != nil {
		return e.err
	}
	return nil
}

func (e *Encoder) writeVaruint7(v varuint7) {
	if e.err != nil {
		return
	}

	e.err = v.write(e.w)
}

func (e *Encoder) writeVaruint32(v varuint32) {
	if e.err != nil {
		return
	}

	_, e.err = v.write(e.w)
}

func (e *Encoder) writeModule(m Module) {
	if e.err != nil {
		return
	}

	e.writeHeader(m.Header)
	for _, s := range m.Sections {
		e.writeSection(s)
	}
}

func (e *Encoder) writeHeader(hdr ModuleHeader) {
	if e.err != nil {
		return
	}

	_, e.err = e.w.Write(hdr.Magic[:])
	if e.err != nil {
		return
	}

	e.err = binary.Write(e.w, order, hdr.Version)
}

func (e *Encoder) writeSection(sec Section) {
	if e.err != nil {
		return
	}

	switch s := sec.(type) {
	case TypeSection:
		e.writeTypeSection(s)
	default:
		panic("TODO")
	}
}

func (e *Encoder) writeTypeSection(s TypeSection) {
	if e.err != nil {
		return
	}

	e.writeVaruint32(varuint32(len(s.types)))
	for _, t := range s.types {
		e.writeFuncType(t)
	}
}

func (e *Encoder) writeFuncType(ft FuncType) {
	if e.err != nil {
		return
	}

	e.writeVaruint7(varuint7(ft.form))

	e.writeVaruint32(varuint32(len(ft.params)))
	for _, v := range ft.params {
		e.writeValueType(v)
	}

	e.writeVaruint32(varuint32(len(ft.results)))
	for _, v := range ft.results {
		e.writeValueType(v)
	}
}

func (e *Encoder) writeValueType(v ValueType) {
	if e.err != nil {
		return
	}

	e.writeVaruint7(varuint7(v))
}
