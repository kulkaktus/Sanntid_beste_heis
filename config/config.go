package config

import (
	"runtime"
)

const (
	INTERNAL = iota
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
	DISTANCE_COST           = 1
	STOPS_INBETWEEN_COST    = 1
	OPPOSITE_DIRECTION_COST = 5
)

func Config_init() error {
	runtime.GOMAXPROCS(4)
	return nil
}
