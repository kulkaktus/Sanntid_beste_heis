package buttons

import (
	"./../../config"
	"./../channels"
	"./../elevio"
)

func Buttons_init() error {
	return nil
}

func Get(type_ int, floor int) bool {
	if type_ == config.INSIDE {
		if floor > 0 && floor <= config.NUMFLOORS {
			return elevio.Get_bit(channels.Inside_button(floor))
		} else {
			return false
		}
	} else if type_ == config.UP {
		if floor > 0 && floor <= config.NUMFLOORS {
			return elevio.Get_bit(channels.Up_button(floor))
		} else {
			return false
		}
	} else if type_ == config.DOWN {
		if floor > 0 && floor <= config.NUMFLOORS {
			return elevio.Get_bit(channels.Down_button(floor))
		} else {
			return false
		}
	} else if type_ == config.STOP {
		return false
	}
	return false
}
