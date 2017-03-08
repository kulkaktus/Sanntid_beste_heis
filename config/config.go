package config

import (
	"runtime"
)

const (
	INSIDE   = iota
	UP       = iota
	DOWN     = iota
	INDICATE = iota
	STOP     = iota
)

const (
	NUMFLOORS    = 4
	NUMELEVATORS = 3
)

const (
	Motor_speed = 2800
)

func Config_init() error {
	runtime.GOMAXPROCS(4)
	return nil
}
