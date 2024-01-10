package main

import (
	"bytes"
	"slices"
	"testing"
)

func TestEncodeStream(t *testing.T) {
	tests := []struct {
		spec []byte
		give []byte
		want []byte
	}{
		{
			spec: []byte{0x02, 0x01, 0x01, 0x00},
			give: []byte("hello"),
			want: []byte{0x96, 0x26, 0xb6, 0xb6, 0x76},
		},
		{
			spec: []byte{0x05, 0x05, 0x00},
			give: []byte("hello"),
			want: []byte{0x68, 0x67, 0x70, 0x72, 0x77},
		},
		{
			spec: []byte{0x03, 0x05, 0x02, 155, 0x00},
			give: []byte("hello"),
			want: []byte{0xf3, 0xfe, 0xeb, 0xe9, 0xf4},
		},
	}

	for _, test := range tests {
		var buf bytes.Buffer
		w := NewObfuscatedWriter(&buf, test.spec)
		w.Write(test.give)
		value := buf.Bytes()
		if !slices.Equal(value, test.want) {
			t.Errorf("want %v have %v", test.want, value)
		}
	}
}

func TestReverseSpec(t *testing.T) {
	tests := []struct {
		give []byte
		want []byte
	}{
		{
			give: []byte{0x02, 0x01, 0x01, 0x00},
			want: []byte{0x01, 0x02, 0x01, 0x00},
		},
	}

	for _, test := range tests {
		reversed := reverseSpec(test.give)
		if !slices.Equal(reversed, test.want) {
			t.Errorf("want %v have %v", test.want, reversed)
		}
	}
}

func TestDecodeStream(t *testing.T) {
	tests := []struct {
		spec []byte
		give []byte
		want []byte
	}{
		{
			spec: []byte{0x02, 0x01, 0x01, 0x00},
			give: []byte{0x96, 0x26, 0xb6, 0xb6, 0x76},
			want: []byte("hello"),
		},
		{
			spec: []byte{0x05, 0x05, 0x00},
			give: []byte{0x68, 0x67, 0x70, 0x72, 0x77},
			want: []byte("hello"),
		},
		{
			spec: []byte{0x03, 0x05, 0x02, 155, 0x00},
			give: []byte{0xf3, 0xfe, 0xeb, 0xe9, 0xf4},
			want: []byte("hello"),
		},
	}

	for _, test := range tests {
		r := NewObfuscatedReader(bytes.NewReader(test.give), test.spec)
		p := make([]byte, len(test.want))

		n, _ := r.Read(p)
		if !slices.Equal(test.want, p[:n]) {
			t.Errorf("want %v have %v", test.want, p[:n])
		}
	}
}

func TestIsNoop(t *testing.T) {
	specs := [][]byte{
		{0x00}, // Empty spec
		{0x02, 0x00, 0x00},
		{0x02, 0xab, 0x02, 0xab, 0x00},
		{0x01, 0x01, 0x00},
		{0x02, 0xa0, 0x02, 0x0b, 0x02, 0xab, 0x00},
	}

	for _, spec := range specs {
		if !isNoop(spec) {
			t.Errorf("%v is not a no-op spec", spec)
		}
	}
}

func TestIsNotNoop(t *testing.T) {
	specs := [][]byte{
		{0x02, 0x01, 0x00},
		{0x01, 0x05, 0x01, 0x00},
	}

	for _, spec := range specs {
		if isNoop(spec) {
			t.Errorf("%v is a no-op spec", spec)
		}
	}
}
