package main

import (
	"os"
	"testing"
)

func TestReadFileFile(t *testing.T) {
	b, err := readFile("main.go")
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatalf("should read something: %v", err)
	}
}

func TestReadFileStdin(t *testing.T) {
	f, err := os.Open("main.go")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	stdin := os.Stdin
	os.Stdin = f
	defer func() {
		os.Stdin = stdin
	}()

	b, err := readFile("-")
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 {
		t.Fatalf("should read something: %v", err)
	}
}

func TestTextContent(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "", want: ""},
		{input: "<p>foo</p>", want: "foo"},
		{input: "<p>foo<span>\nbar\n</span>baz</p>", want: "foobarbaz"},
		{input: "<p>foo<span>\nbar<br></span>baz</p>", want: "foobar\nbaz"},
	}
	for _, test := range tests {
		got := textContent(test.input)
		if got != test.want {
			t.Fatalf("want %q but %q", test.want, got)
		}
	}
}
