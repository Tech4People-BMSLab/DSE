package debug

import "github.com/davecgh/go-spew/spew"

func Dump(args ...interface{}) {
	spew.Dump(args...)
}

func Print(args ...interface{}) {
	spew.Print(args...)
}

func Breakpoint() {
	return
}
