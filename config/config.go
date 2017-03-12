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
	DOOR     = iota
)

const (
	NUMFLOORS       = 4
	NUMELEVATORS    = 3
	NUMBUTTON_TYPES = 3
)

const (
	DISTANCE_COST          = 1
	STARTUP_FROM_IDLE_COST = 1
	DIRECTION_CHANGE_COST  = 5
)

const (
	Motor_speed = 2800
)

func Init() error {
	runtime.GOMAXPROCS(4)
	return nil
}
