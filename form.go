package main

var activeForm *form

type form struct {
	Stage int
	Fn    func(*form, string)
}
