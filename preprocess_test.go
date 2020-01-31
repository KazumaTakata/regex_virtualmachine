package main

import (
	"testing"
)

func TestPreprocess(t *testing.T) {

	input := "a+b+"
	expected := "a+,b+"

	output := Preprocess(input)

	if output != expected {
		t.Errorf("Preprocess expected:%s, got:%s", expected, output)
	}
}
