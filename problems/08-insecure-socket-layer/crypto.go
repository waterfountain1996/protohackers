package main

import (
	"bytes"
	"io"
	"math/bits"
	"slices"
)

type cipherFunc func(uint64, byte) byte

func xorN(n byte) func(byte) byte {
	return func(b byte) byte {
		return b ^ n
	}
}

func addN(n byte) func(byte) byte {
	return func(b byte) byte {
		return b + n
	}
}

func subN(n byte) func(byte) byte {
	return func(b byte) byte {
		return b - n
	}
}

func fromCipherSpec(spec []byte) cipherFunc {
	return func(pos uint64, b byte) byte {
		for i := 0; i < len(spec)-1; i++ {
			var f func(byte) byte

			switch op := spec[i]; op {
			case 0x01:
				f = bits.Reverse8
			case 0x02:
				i++
				arg := spec[i]
				f = xorN(arg)
			case 0x03:
				f = xorN(byte(pos % 256))
			case 0x04:
				i++
				arg := spec[i]
				f = addN(arg)
			case 0x05:
				f = addN(byte(pos % 256))
			case 0x06: // Inverse of add(N)
				i++
				arg := spec[i]
				f = subN(arg)
			case 0x07: // Inverse of addpos
				f = subN(byte(pos % 256))
			}

			b = f(b)
		}

		return b
	}
}

func reverseSpec(spec []byte) []byte {
	reversed := make([]byte, len(spec))
	for i := 0; i < len(spec)-1; i++ {
		op := spec[i]

		if op == 0x02 || op == 0x04 {
			arg := spec[i+1]
			reversed[len(spec)-2-i] = arg
			i++
		}

		if op == 0x04 || op == 0x05 {
			op += 2
		}

		reversed[len(spec)-2-i] = op
	}
	reversed[len(reversed)-1] = 0x00
	return reversed
}

type ObfuscatedWriter struct {
	io.Writer

	pos       uint64
	obfuscate cipherFunc
}

func NewObfuscatedWriter(w io.Writer, spec []byte) *ObfuscatedWriter {
	return &ObfuscatedWriter{
		Writer:    w,
		obfuscate: fromCipherSpec(spec),
	}
}

func (w *ObfuscatedWriter) Write(p []byte) (int, error) {
	for i, b := range p {
		p[i] = w.obfuscate(w.pos, b)
		w.pos++
	}
	return w.Writer.Write(p)
}

type ObfuscatedReader struct {
	io.Reader

	pos         uint64
	deobfuscate cipherFunc
}

func NewObfuscatedReader(r io.Reader, spec []byte) *ObfuscatedReader {
	return &ObfuscatedReader{
		Reader:      r,
		deobfuscate: fromCipherSpec(reverseSpec(spec)),
	}
}

func (r *ObfuscatedReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	if err != nil {
		return n, err
	}

	for i, b := range p[:n] {
		p[i] = r.deobfuscate(r.pos, b)
		r.pos++
	}

	return n, nil
}

func isNoop(spec []byte) bool {
	var buf bytes.Buffer

	w := NewObfuscatedWriter(&buf, spec)
	w.Write([]byte("Hello, world!"))

	return slices.Equal(buf.Bytes(), []byte("Hello, world!"))
}
