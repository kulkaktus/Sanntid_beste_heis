package motor

import (
	"./../../config"
	"./../channels"
	"./../io"
)

func Motor_init(){
	config.Config_init()
	io.Write_analog(channels.Motor_value, 0)
}

func Go(direction int) {
	if(direction == config.UP){
		io.Set_bit(channels.Motor_dir)
		io.Write_analog(channels.Motor_value, config.Motor_speed)
	}else if(direction == config.DOWN){
		io.Clear_bit(channels.Motor_dir)
		io.Write_analog(channels.Motor_value, config.Motor_speed)
	}
}

func Stop(){
	io.Write_analog(channels.Motor_value, 0)
}