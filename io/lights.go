package lights

import (
	"./elevio"
	"./channels"
	"errors"
)

const(
	INSIDE = 1
	UP = 2
	DOWN = 3
	INDICATE = 4
	STOP = 5
)


func set(type int, floor int)error{
	if(type == INSIDE){
		if(floor > 0 && floor <= 4/*config.num_floors*/){
			elevio.set_bit(channels.Inside_light(floor-1))
		}
		else{
			return
		}
	}
	else if(type == UP){
		if(floor > 0 && floor <= 4/*config.num_floors*/){
			elevio.set_bit(channels.Up_light(floor-1))
		}
		else{
			err = 1
			return
		}
	}
	else if(type == DOWN){
		if(floor > 0 && floor <= 4/*config.num_floors*/){
			elevio.set_bit(channels.Down_light(floor-1))
		}
		else{
			err = 1
			return
		}
	}
	else if(type == INDICATE){
		if(floor > 0 && floor <= 4/*config.num_floors*/){
			elevio.set_bit(channels.Floor_light_0(floor%2))
			elevio.set_bit(channels.Floor_light_1(floor/2))
		}
		else{
			err = 1
			return
		}
	}
	else if(type == STOP){
		elevio.set_bit(channels.Stop_light)
	}
}