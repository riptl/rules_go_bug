package golib

//#include "clib/lib.h"
import "C"

func hello() {
	if C.hello() != 42 {
		panic("hello() failed")
	}
}
