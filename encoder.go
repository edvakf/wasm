// Copyright 2016 The wasm Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasm

import (
	"bytes"
	"encoding/binary"
	"io"
)

//Encode writes the wasm binary to the writer
func Encode(m Module, w io.Writer) error {
	e := encoder{w: w}
	e.writeModule(m)

	if e.err != nil {
		return e.err
	}
	return nil
}

type encoder struct {
	w   io.Writer
	err error
}

func (e *encoder) writeVaruint7(v varuint7) {
	if e.err != nil {
		return
	}

	e.err = v.write(e.w)
}

func (e *encoder) writeVaruint32(v varuint32) {
	if e.err != nil {
		return
	}

	_, e.err = v.write(e.w)
}

func (e *encoder) writeVarint7(v varint7) {
	if e.err != nil {
		return
	}

	e.err = v.write(e.w)
}

func (e *encoder) writeVarint32(v varint32) {
	if e.err != nil {
		return
	}

	_, e.err = v.write(e.w)
}

func (e *encoder) writeVarint64(v varint64) {
	if e.err != nil {
		return
	}

	_, e.err = v.write(e.w)
}

func (e *encoder) writeModule(m Module) {
	if e.err != nil {
		return
	}

	e.writeHeader(m.Header)
	for _, s := range m.Sections {
		e.writeSection(s)
	}
}

func (e *encoder) writeHeader(hdr ModuleHeader) {
	if e.err != nil {
		return
	}

	_, e.err = e.w.Write(hdr.Magic[:])
	if e.err != nil {
		return
	}

	e.err = binary.Write(e.w, order, hdr.Version)
}

func (e *encoder) writeSection(sec Section) {
	if e.err != nil {
		return
	}

	e.writeVaruint7(varuint7(sec.ID())) // id

	b := new(bytes.Buffer)
	encSec := &encoder{w: b}
	switch s := sec.(type) {
	case TypeSection:
		encSec.writeTypeSection(s)
	case FunctionSection:
		encSec.writeFunctionSection(s)
	case ExportSection:
		encSec.writeExportSection(s)
	case CodeSection:
		encSec.writeCodeSection(s)
	default:
		panic("TODO")
	}
	e.writeVaruint32(varuint32(b.Len())) // payload_len
	_, e.err = e.w.Write(b.Bytes())

	// TODO: name_len, name
}

func (e *encoder) writeTypeSection(s TypeSection) {
	if e.err != nil {
		return
	}

	e.writeVaruint32(varuint32(len(s.types)))
	for _, t := range s.types {
		e.writeFuncType(t)
	}
}

func (e *encoder) writeFuncType(ft FuncType) {
	if e.err != nil {
		return
	}

	e.writeVarint7(varint7(ft.form))

	e.writeVaruint32(varuint32(len(ft.params)))
	for _, v := range ft.params {
		e.writeValueType(v)
	}

	e.writeVaruint32(varuint32(len(ft.results)))
	for _, v := range ft.results {
		e.writeValueType(v)
	}
}

func (e *encoder) writeValueType(v ValueType) {
	if e.err != nil {
		return
	}

	e.writeVarint7(varint7(v))
}

func (e *encoder) writeFunctionSection(s FunctionSection) {
	if e.err != nil {
		return
	}

	e.writeVaruint32(varuint32(len(s.types)))
	for _, t := range s.types {
		e.writeVaruint32(varuint32(t))
	}
}

func (e *encoder) writeCodeSection(s CodeSection) {
	if e.err != nil {
		return
	}

	e.writeVaruint32(varuint32(len(s.Bodies)))
	for _, b := range s.Bodies {
		e.writeFunctionBody(b)
	}
}

func (e *encoder) writeFunctionBody(fb FunctionBody) {
	if e.err != nil {
		return
	}

	e.writeVaruint32(varuint32(fb.BodySize))
	e.writeVaruint32(fb.LocalCount)
	for _, l := range fb.Locals {
		e.writeLocalEntry(l)
	}
	e.writeCode(fb.Code)
}

func (e *encoder) writeLocalEntry(l LocalEntry) {
	if e.err != nil {
		return
	}

	e.writeVaruint32(varuint32(l.Count))
	e.writeValueType(l.Type)
}

func (e *encoder) writeExportSection(s ExportSection) {
	if e.err != nil {
		return
	}

	e.writeVaruint32(varuint32(len(s.exports)))
	for _, ex := range s.exports {
		e.writeExportEntry(ex)
	}
}

func (e *encoder) writeExportEntry(ex ExportEntry) {
	if e.err != nil {
		return
	}

	e.writeVaruint32(varuint32(len(ex.field)))
	_, e.err = e.w.Write([]byte(ex.field))
	if e.err != nil {
		return
	}
	e.writeExternalKind(ex.kind)
	e.writeVaruint32(varuint32(ex.index))
}

func (e *encoder) writeExternalKind(k ExternalKind) {
	if e.err != nil {
		return
	}

	_, e.err = e.w.Write([]byte{byte(k)})
}

func (e *encoder) writeCode(c Code) {
	if e.err != nil {
		return
	}

	_, e.err = e.w.Write(c.Code)
	if e.err != nil {
		return
	}
	_, e.err = e.w.Write([]byte{c.End})
}
