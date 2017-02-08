package elevio

/*
#cgo LDFLAGS: -lio
#include <io.h>
*/
import (
	"C"
)

func init() bool {
	return (int(C.io_init()) != 0)
}

func set_bit(channel int) {
	C.io_set_bit(C.int(channel))
}

func clear_bit(int channel) {
	C.io_clear_bit(C.int(channel))
}

func read_bit(int channel) int {
	return (C.io_read_bit(C.int(channel)) != 0)
}

func write_analog(int channel, int value) {
	C.io_write_analog(C.int(channel), C.int(value))
}
