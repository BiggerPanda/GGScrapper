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
	}{{"test1", args{"<span class=\"price\">$ 1,99</span>"}, []float64{1.99}}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := utility.ParsePrice(tt.args.price); got[0] != tt.want[0] {
				t.Errorf("ParsePrice() = %v, want %v", got, tt.want)
			}
		})
	}

}
