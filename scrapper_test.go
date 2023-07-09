package main

import (
	"testing"
	"web-scrapper/utility"
)

func TestParsePrice(t *testing.T) {
	type args struct {
		price string
	}
	tests := []struct {
		name string
		args args
		want []float64
	}{{"test1", args{"<span class=\"price\">$ 1,99</span>"}, []float64{1.99}},
		{"test2", args{"<span class=\"price\">$ 1,99</span><span class=\"price\">$ 2,99</span>"}, []float64{1.99, 2.99}},
		{"test3", args{"<span class=\"price\">$ 1,99</span><span class=\"price\">$ 2,99</span><span class=\"price\">$ 3,99</span>"}, []float64{1.99, 2.99, 3.99}}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utility.ParsePrice(tt.args.price); got[0] != tt.want[0] {
				t.Errorf("ParsePrice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBrowserOpen(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{{"test1", args{"https://www.google.com/"}, true}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utility.OpenBrowser(tt.args.url); got != tt.want {
				t.Errorf("OpenBrowser() = %v, want %v", got, tt.want)
			}
		})
	}
}
