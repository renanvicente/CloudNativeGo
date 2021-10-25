package main

type fox struct {
}

func (f fox) Says() string {
	return "ring-ding-ding-ding-dingeringeding!"
}

var Animal fox
