package debug

import "fmt"

type _Debug struct {
	On bool
}

func (d *_Debug) Printf(format string, a ...any) {
	if !d.On {
		return
	}
	fmt.Printf(format, a...)
}

func (d *_Debug) Println(a ...any) {
	if !d.On {
		return
	}
	fmt.Println(a...)
}

var Debug _Debug
