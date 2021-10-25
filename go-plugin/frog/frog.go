package main

type frog struct {
}

func (f frog) Says() string {
	return "ribbit!"
}

var Animal frog
