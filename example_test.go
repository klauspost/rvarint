// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rvarint_test

import (
	"encoding/hex"
	"fmt"

	"github.com/klauspost/rvarint"
)

func ExamplePutUvarint() {
	buf := make([]byte, rvarint.MaxVarintLen64)

	for _, x := range []uint64{1, 2, 127, 128, 255, 256} {
		n := rvarint.PutUvarint(buf, x)
		fmt.Printf("%x\n", buf[:n])
	}
	// Output:
	// 01
	// 02
	// 7f
	// 0180
	// 01ff
	// 0280
}

func ExamplePutVarint() {
	buf := make([]byte, rvarint.MaxVarintLen64)

	for _, x := range []int64{-65, -64, -2, -1, 0, 1, 2, 63, 64} {
		n := rvarint.PutVarint(buf, x)
		fmt.Printf("%x\n", buf[:n])
	}
	// Output:
	// 0181
	// 7f
	// 03
	// 01
	// 00
	// 02
	// 04
	// 7e
	// 0180
}

func ExampleAppendUvarint() {
	// Add some existing content to the buffer
	buf := []byte("prefix")

	// values to encode
	values := []uint64{1, 2, 127, 128, 255, 256}
	for _, x := range values {
		buf = rvarint.AppendUvarint(buf, x)
	}
	fmt.Printf("encoded: %s\n", hex.EncodeToString(buf))

	// When we read back we get the values in *reverse* order.
	for range values {
		y, n := rvarint.Uvarint(buf)
		if n <= 0 {
			fmt.Println("unable to read value")
			break
		}
		fmt.Println(y)
		buf = buf[:len(buf)-n]
	}
	fmt.Println("remaining buffer:", string(buf))
	// Output:
	// encoded: 70726566697801027f018001ff0280
	// 256
	// 255
	// 128
	// 127
	// 2
	// 1
	// remaining buffer: prefix
}

func ExampleAppendVarint() {
	// Add some existing content to the buffer
	buf := []byte("prefix")

	// values to encode
	values := []int64{-65, -64, -2, -1, 0, 1, 2, 63, 64}
	for _, x := range values {
		buf = rvarint.AppendVarint(buf, x)
	}
	fmt.Printf("encoded: %s\n", hex.EncodeToString(buf))

	// When we read back we get the values in *reverse* order.
	for range values {
		y, n := rvarint.Varint(buf)
		if n <= 0 {
			fmt.Println("unable to read value")
			break
		}
		fmt.Println(y)
		buf = buf[:len(buf)-n]
	}
	fmt.Println("remaining buffer:", string(buf))
	// Output:
	// encoded: 70726566697801817f03010002047e0180
	// 64
	// 63
	// 2
	// 1
	// 0
	// -1
	// -2
	// -64
	// -65
	// remaining buffer: prefix
}

func ExampleUvarint() {
	inputs := [][]byte{
		{0x01},
		{0x02},
		{0x7f},
		{0x01, 0x80},
		{0x01, 0xff},
		{0x02, 0x80},
	}
	for _, b := range inputs {
		x, n := rvarint.Uvarint(b)
		if n != len(b) {
			fmt.Println("Uvarint did not consume all of in")
		}
		fmt.Println(x)
	}
	// Output:
	// 1
	// 2
	// 127
	// 128
	// 255
	// 256
}

func ExampleVarint() {
	inputs := [][]byte{
		{0x01, 0x81},
		{0x7f},
		{0x03},
		{0x01},
		{0x00},
		{0x02},
		{0x04},
		{0x7e},
		{0x01, 0x80},
	}
	for _, b := range inputs {
		x, n := rvarint.Varint(b)
		if n != len(b) {
			fmt.Println("Varint did not consume all of in", len(b)-n, "left")
		}
		fmt.Println(x)
	}
	// Output:
	// -65
	// -64
	// -2
	// -1
	// 0
	// 1
	// 2
	// 63
	// 64
}
