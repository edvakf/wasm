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
	return e.writeModule(m)
}

func (e *Encoder) writeModule(m Module) error {
	if e.err != nil {
		return e.err
	}

	e.err = e.writeHeader(m.Header)
	for _, s := range m.Sections {
		e.err = e.writeSection(s)
	}
	// if err != nil {
	// 	e.err = err
	// 	return err
	// }
	return e.err
}

func (e *Encoder) writeHeader(hdr ModuleHeader) error {
	if e.err != nil {
		return e.err
	}

	_, err := e.w.Write(hdr.Magic[:])
	if err != nil {
		return err
	}

	err = binary.Write(e.w, order, hdr.Version)
	if err != nil {
		return err
	}

	return nil
}

func (e *Encoder) writeSection(s Section) error {
	if e.err != nil {
		return e.err
	}

	switch s.ID() {
	case TypeID:
		e.err = e.writeTypeSection(TypeSection(s))
	}

	return nil
}

func (e *Encoder) writeTypeSection(s TypeSection) error {
	if e.err != nil {
		return e.err
	}

	//

	return nil
}
