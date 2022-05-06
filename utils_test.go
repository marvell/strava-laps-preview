package main

import "testing"

func TestFormatDistance(t *testing.T) {
	tests := []struct {
		in  float64
		out string
	}{
		{100, "100m"},
		{800, "800m"},
		{995, "1km"},
		{1000, "1km"},
		{1005, "1.01km"},
		{1100, "1.1km"},
		{1150, "1.15km"},
		{1155, "1.16km"},
		{2997.54, "3km"},
		{3000.36, "3km"},
		{10000, "10km"},
		{10100, "10.1km"},
	}

	for _, test := range tests {
		got := FormatDistance(test.in)
		if got != test.out {
			t.Errorf("FormatDistance(%f): %s, but wants %s", test.in, got, test.out)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		in  int
		out string
	}{
		{5, "00:05"},
		{65, "01:05"},
		{605, "10:05"},
		{3605, "01:00:05"},
	}

	for _, test := range tests {
		got := FormatDuration(test.in)
		if got != test.out {
			t.Errorf("FormatDuration(%d): %s, but wants %s", test.in, got, test.out)
		}
	}
}
