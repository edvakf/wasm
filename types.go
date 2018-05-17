// Copyright 2016 The wasm Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wasm

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

var order = binary.LittleEndian

type (
	varuint32 uint32
	varuint7  uint32
	varuint1  uint32

	varint64 int64
	varint32 int32
	varint7  int32
)

func (v *varuint32) read(r io.Reader) (int, error) {
	vv, n, err := uvarint(r)
	if err != nil {
		return n, err
	}
	*v = varuint32(vv)
	return n, nil
}

func (v *varuint7) read(r io.Reader) (int, error) {
	vv, n, err := uvarint(r)
	if err != nil {
		return n, err
	}
	*v = varuint7(vv)
	return n, nil
}

func uvarint(r io.Reader) (uint32, int, error) {
	var x uint32
	var s uint
	var buf = make([]byte, 1)
	for i := 0; ; i++ {
		_, err := r.Read(buf)
		if err != nil {
			return 0, i, err
		}
		b := buf[0]
		if b < 0x80 {
			if i > 9 || i == 9 && b > 1 {
				return 0, i, errors.New("wasm: overflow")
			}
			return x | uint32(b)<<s, i, nil
		}
		x |= uint32(b&0x7f) << s
		s += 7
	}
	panic("unreachable")
}

func varint(r io.Reader) (int32, int, error) {
	var result int32 = 0
	var shift uint = 0
	var size uint = 32
	buf := make([]byte, 1)
	n := 0
	for {
		_, err := r.Read(buf)
		if err != nil {
			return 0, n, err
		}
		n++
		b := int32(buf[0])
		result |= (b << shift) & 0x7F
		shift += 7
		if b&0x80 == 0 {
			if shift < size && b&0x40 != 0 {
				result |= ^0 << shift
			}
			break
		}
	}
	return result, n, nil
}

func (v varuint32) write(w io.Writer) (int, error) {
	n := 0 // how many bytes written
	for {
		b := v & 0x7F
		v >>= 7
		if v != 0 {
			b |= 0x80
		}
		_, err := w.Write([]byte{uint8(b)})
		if err != nil {
			return n, err
		}
		n++
		if v == 0 {
			break
		}
	}
	return n, nil
}

func (v varuint7) write(w io.Writer) error {
	if v > 0x7F {
		return fmt.Errorf("too large for varuint7: %d", v)
	}
	_, err := w.Write([]byte{uint8(v)})
	if err != nil {
		return err
	}
	return nil
}

func (v varint64) write(w io.Writer) (int, error) {
	n := 0
	more := true
	for more {
		b := v & 0x7F
		v >>= 7 // arithmetic shift
		if (v == 0 && b&0x40 == 0) || (v == -1 && b&0x40 != 0) {
			more = false
		} else {
			b |= 0x80
		}
		_, err := w.Write([]byte{uint8(b)})
		if err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

// code is exactly the same as varint64.write
func (v varint32) write(w io.Writer) (int, error) {
	// varint64(v).write(w)
	n := 0
	more := true
	for more {
		b := v & 0x7F
		v >>= 7 // arithmetic shift
		if (v == 0 && b&0x40 == 0) || (v == -1 && b&0x40 != 0) {
			more = false
		} else {
			b |= 0x80
		}
		_, err := w.Write([]byte{uint8(b)})
		if err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

type ValueType varint7

type BlockType ValueType
type ElemType ValueType

type FuncType struct {
	form    ValueType   // value for the 'func' type constructor
	params  []ValueType // parameters of the function
	results []ValueType // results of the function
}

// GlobalType describes a global variable
type GlobalType struct {
	ContentType ValueType
	Mutability  varuint1 // 0:immutable, 1:mutable
}

// TableType describes a table
type TableType struct {
	ElemType ElemType // the type of elements
	Limits   ResizableLimits
}

// MemoryType describes a memory
type MemoryType struct {
	Limits ResizableLimits
}

// ExternalKind indicates the kind of definition being imported or defined:
// 0: indicates a Function import or definition
// 1: indicates a Table import or definition
// 2: indicates a Memory import or definition
// 3: indicates a Global import or definition
type ExternalKind byte

// 0: indicates a Function import or definition
// 1: indicates a Table import or definition
// 2: indicates a Memory import or definition
// 3: indicates a Global import or definition
const (
	FunctionKind ExternalKind = 0
	TableKind                 = 1
	MemoryKind                = 2
	GlobalKind                = 3
)

// ResizableLimits describes the limits of a table or memory
type ResizableLimits struct {
	Flags   uint32 // bit 0x1 is set if the maximum field is present
	Initial uint32 // initial length (in units of table elements or wasm pages)
	Maximum uint32 // only present if specified by Flags
}

// InitExpr encodes an initializer expression.
// FIXME(sbinet)
type InitExpr struct {
	Expr []byte
	End  byte
}
