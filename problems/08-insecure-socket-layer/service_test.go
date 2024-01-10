package main

import "testing"

func TestExtractQuantity(t *testing.T) {
	tests := []struct {
		give string
		want int
	}{
		{
			give: "10x toy car",
			want: 10,
		},
		{
			give: "15x dog on a string",
			want: 15,
		},
		{
			give: "4x inflatable motorcycle",
			want: 4,
		},
	}

	for _, test := range tests {
		value := extractQuantity(test.give)
		if value != test.want {
			t.Errorf("want %d have %d", test.want, value)
		}
	}
}

func TestFindToy(t *testing.T) {
	tests := []struct {
		request string
		want    string
	}{
		{
			request: "10x toy car,15x dog on a string,4x inflatable motorcycle",
			want:    "15x dog on a string",
		},
		{
			request: "4x dog,5x car",
			want:    "5x car",
		},
		{
			request: "3x rat,2x cat",
			want:    "3x rat",
		},
	}

	for _, test := range tests {
		value := FindToy(test.request)
		if value != test.want {
			t.Errorf("want %s have %s", test.want, value)
		}
	}
}
