package sensors

import (
	"./../../config"
	"./../channels"
	"./../elevio"
)

func Sensors_init() error {
	return nil
}

func Get() int {
	floor := 0
	for floor = config.NUMFLOORS; floor >= 1; floor-- {
		if elevio.Get_bit(channels.Sensor(floor)) {
			break
		}
	}
	return floor
}
