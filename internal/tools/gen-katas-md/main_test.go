// SPDX-License-Identifier: MIT
// Copyright the ZShellCheck contributors.
package main

import "testing"

func TestIdNum(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"ZC1001", 1001},
		{"ZC2003", 2003},
		{"ZC0001", 1},
	}
	for _, tc := range cases {
		if got := idNum(tc.in); got != tc.want {
			t.Errorf("idNum(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}

func TestIdNumEdge(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"", 0},
		{"ZC", 0},
		{"X", 0},
		{"ZCabc", 0},
	}
	for _, tc := range cases {
		if got := idNum(tc.in); got != tc.want {
			t.Errorf("idNum(%q) = %d, want %d", tc.in, got, tc.want)
		}
	}
}
