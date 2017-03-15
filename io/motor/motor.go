package motor

import (
	"./../../config"
	"./../channels"
	"./../io"
)

func Init() {
	config.Init()
	io.Write_analog(channels.Motor_value, 0)
}

func Move(direction int) {
	if direction == config.UP {
		io.Clear_bit(channels.Motor_dir)
		io.Write_analog(channels.Motor_value, config.Motor_speed)
	} else if direction == config.DOWN {
		io.Set_bit(channels.Motor_dir)
		io.Write_analog(channels.Motor_value, config.Motor_speed)
	}
}

func Stop() {
	io.Write_analog(channels.Motor_value, 0)
}
