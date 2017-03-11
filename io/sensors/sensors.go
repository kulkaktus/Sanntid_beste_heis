package sensors

import (
	"./../../config"
	"./../channels"
	"./../io"
)

func Init() error {
	return nil
}

func Get() int {
	floor := 0
	for floor = config.NUMFLOORS; floor >= 1; floor-- {
		if io.Get_bit(channels.Sensor(floor)) {
			break
		}
	}
	return floor
}
