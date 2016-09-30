package main

import "testing"

func TestGenerate(t *testing.T) {
	// Passing in no arguments returns no results
	output := generate()
	if len(output) > 0 {
		t.Fail()
	}

	// Passing in one argument returns one result
	output = generate("GDG")
	if len(output) != 1 {
		t.Fail()
	}
	if output[0] != "Hello GDG" {
		t.Fail()
	}

	// Passing in three arguments returns three results
	output = generate("GDG", "World", "Gophers")
	if len(output) != 3 {
		t.Fail()
	}
	if output[0] != "Hello GDG" {
		t.Fail()
	}
	if output[1] != "Hello World" {
		t.Fail()
	}
	if output[2] != "Hello Gophers" {
		t.Fail()
	}
}
