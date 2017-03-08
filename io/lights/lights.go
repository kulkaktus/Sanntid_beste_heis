package lights

import (
	"./../../config"
	"./../channels"
	"./../io"
)

func Lights_init() error {
	return nil
}

func Set(type_ int, floor int) error {
	if type_ == config.INSIDE {
		if floor > 0 && floor <= config.NUMFLOORS {
			io.Set_bit(channels.Inside_light(floor))
		} else {
			return nil
		}
	} else if type_ == config.UP {
		if floor > 0 && floor <= config.NUMFLOORS {
			io.Set_bit(channels.Up_light(floor))
		} else {
			return nil
		}
	} else if type_ == config.DOWN {
		if floor > 0 && floor <= config.NUMFLOORS {
			io.Set_bit(channels.Down_light(floor))
		} else {
			return nil
		}
	} else if type_ == config.INDICATE {
		if floor > 0 && floor <= config.NUMFLOORS {
			if ((floor - 1) % 2) != 0 {
				io.Set_bit(channels.Floor_light_1)
			} else {
				io.Clear_bit(channels.Floor_light_1)
			}
			if ((floor - 1) / 2) != 0 {
				io.Set_bit(channels.Floor_light_0)
			} else {
				io.Clear_bit(channels.Floor_light_0)
			}
		} else {
			return nil
		}
	} else if type_ == config.STOP {
		io.Set_bit(channels.Stop_light)
	}
	return nil
}

func Clear(type_ int, floor int) error {
	if type_ == config.INSIDE {
		if floor > 0 && floor <= config.NUMFLOORS {
			io.Clear_bit(channels.Inside_light(floor))
		} else {
			return nil
		}
	} else if type_ == config.UP {
		if floor > 0 && floor <= config.NUMFLOORS {
			io.Clear_bit(channels.Up_light(floor))
		} else {
			return nil
		}
	} else if type_ == config.DOWN {
		if floor > 0 && floor <= config.NUMFLOORS {
			io.Clear_bit(channels.Down_light(floor))
		} else {
			return nil
		}
	} else if type_ == config.STOP {
		io.Clear_bit(channels.Stop_light)
	}
	return nil
}
