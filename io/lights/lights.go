package lights

import (
	"./../../config"
	"./../channels"
	"./../io"
)

func Init() error {
	for i := 1; i < config.NUMFLOORS; i++ {
		for j := 0; j < config.NUMLIGHTS; j++ {
			Clear(j, i)
		}
	}
	return nil
}

func Clear_floor(floor int) {
	for i := 1; i < config.NUMBUTTON_TYPES; i++ {
		Clear(i, floor)
	}
}

func Set(type_ int, floor int) error {
	if type_ == config.INTERNAL {
		if floor > 0 && floor <= config.NUMFLOORS {
			io.Set_bit(channels.Internal_light(floor))
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
	} else if type_ == config.DOOR {
		io.Set_bit(channels.Door)
	}
	return nil
}

func Clear(type_ int, floor int) error {
	if type_ == config.INTERNAL {
		if floor > 0 && floor <= config.NUMFLOORS {
			io.Clear_bit(channels.Internal_light(floor))
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
	} else if type_ == config.DOOR {
		io.Clear_bit(channels.Door)
	}
	return nil
}
