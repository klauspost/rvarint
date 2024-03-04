# rvarint

Package rvarint implements reverse varint encoding and decoding.

The encoding is similar to varint, but the bytes are written in reverse order.

[![Go Reference](https://pkg.go.dev/badge/klauspost/rvarint.svg)](https://pkg.go.dev/github.com/klauspost/rvarint?tab=subdirectories)

## Usage

This allows you to read varints from the end of a buffer, which can be useful when reading from a buffer that is being filled from the end.

See [GoDoc](https://pkg.go.dev/github.com/klauspost/rvarint?tab=doc) for examples.

The main difference to [enconding/binary](https://pkg.go.dev/encoding/binary) is that `Uvarint` and `Varint` will read from the *end* of the buffer.

This also means that appending several values will result in the values being read back in *reverse* order.

