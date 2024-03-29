// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package rvarint implements reverse varint encoding and decoding.
//
// The encoding is similar to varint, but the bytes are written in reverse order.
//
// This means that reading will read from the end of the buffer, and the
// encoded values will be read in reverse order.
package rvarint

import (
	"math/bits"
)

// MaxVarintLenN is the maximum length of a varint-encoded N-bit integer.
const (
	MaxVarintLen16 = 3
	MaxVarintLen32 = 5
	MaxVarintLen64 = 10
)

// AppendUvarint appends the reverse varint-encoded form of x
// to buf and returns the extended buffer.
func AppendUvarint(buf []byte, x uint64) []byte {
	n := (bits.Len64(x) - 1) / 7
	sh := uint(n * 7)
	buf = append(buf, byte(x>>sh)&0x7f)
	sh -= 7
	for ; n > 0; n-- {
		buf = append(buf, byte(x>>sh)|0x80)
		sh -= 7
	}
	return buf
}

// PutUvarint encodes a uint64 into buf and returns the number of bytes written.
// If the buffer is too small, PutUvarint will panic.
func PutUvarint(buf []byte, x uint64) int {
	n := (bits.Len64(x) - 1) / 7
	i := n
	for x >= 0x80 {
		buf[i] = byte(x) | 0x80
		x >>= 7
		i--
	}
	buf[i] = byte(x)
	return n + 1
}

// AppendVarint appends the varint-encoded form of x,
// as generated by [PutVarint], to buf and returns the extended buffer.
func AppendVarint(buf []byte, x int64) []byte {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return AppendUvarint(buf, ux)
}

// PutVarint encodes an int64 into buf and returns the number of bytes written.
// If the buffer is too small, PutVarint will panic.
func PutVarint(buf []byte, x int64) int {
	ux := uint64(x) << 1
	if x < 0 {
		ux = ^ux
	}
	return PutUvarint(buf, ux)
}

// Uvarint decodes a uint64 from the end of buf and returns that value and the
// number of bytes read (> 0). If an error occurred, the value is 0
// and the number of bytes n is <= 0 meaning:
//
//	n == 0: buf too small
//	n  < 0: value larger than 64 bits (overflow)
//	        and -n is the number of bytes read
//
// On success, the buffer can be truncated using buf = buf[:len(buf)-n]
func Uvarint(buf []byte) (uint64, int) {
	var x uint64
	var s uint
	for i := len(buf) - 1; i >= 0; i-- {
		if len(buf)-i > MaxVarintLen64 {
			// Catch byte reads past MaxVarintLen64.
			// See issue https://golang.org/issues/41185
			return 0, -(len(buf) - i) // overflow
		}
		b := buf[i]
		if b < 0x80 {
			if len(buf)-i > MaxVarintLen64-1 && b > 1 {
				return 0, -(len(buf) - i) // overflow
			}
			return x | uint64(b)<<s, len(buf) - i
		}
		x |= uint64(b&0x7f) << s
		s += 7
	}
	return 0, 0
}

// Varint decodes an int64 from the end of buf and returns that value and the
// number of bytes read (> 0). If an error occurred, the value is 0
// and the number of bytes n is <= 0 with the following meaning:
//
//	n == 0: buf too small
//	n  < 0: value larger than 64 bits (overflow)
//	        and -n is the number of bytes read
//
// On success, the buffer can be truncated using buf = buf[:len(buf)-n]
func Varint(buf []byte) (int64, int) {
	ux, n := Uvarint(buf) // ok to continue in presence of error
	x := int64(ux >> 1)
	if ux&1 != 0 {
		x = ^x
	}
	return x, n
}
