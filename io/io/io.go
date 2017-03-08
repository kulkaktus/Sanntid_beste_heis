package io

/*
#cgo CFLAGS: -std=gnu11 -g
#cgo LDFLAGS: -lcomedi -lm
#include "io.h"
*/
import (
	"C"
)

func Io_init() bool {
	return (int(C.io_init()) != 0)
}

func Set_bit(channel int) {
	C.io_set_bit(C.int(channel))
}

func Clear_bit(channel int) {
	C.io_clear_bit(C.int(channel))
}

func Get_bit(channel int) bool {
	return (C.io_read_bit(C.int(channel)) != 0)
}

func Write_analog(channel int, value int) {
	C.io_write_analog(C.int(channel), C.int(value))
}
