package order_handling

import (
	"../config"
	"fmt"
	"math"
)

const (
	NO_EXECUTER = 0
	NO_ORDER    = -1
)

var underlying_order_matrix [config.NUMFLOORS][config.NUMBUTTON_TYPES]int
var order_matrix [][config.NUMBUTTON_TYPES]int

var last_floor int
var direction int

var State string
var self int

func Init(self_id int) {
	order_matrix = underlying_order_matrix[:][:]
	for floor_i := 0; floor_i < config.NUMFLOORS; floor_i++ {
		for button_type_i := 0; button_type_i < config.NUMBUTTON_TYPES; button_type_i++ {
			order_matrix[floor_i][button_type_i] = -1
		}
	}

	self = self_id
	last_floor = 2
	direction = config.DOWN
}

func Insert(destination int, button_type int, executer_id int) bool {
	if order_is_in_bounds(destination, button_type) {
		order_matrix[destination-1][button_type] = executer_id
		return true
	}
	return false
}

func Clear_orders_in_floor(destination int) {
	for button_type_i := 0; button_type_i < config.NUMBUTTON_TYPES; button_type_i++ {
		if order_is_in_bounds(destination, button_type_i) {
			order_matrix[destination-1][button_type_i] = NO_ORDER
		}
	}
}

func Clear_order_matrix() {
	for floor_i := 1; floor_i <= config.NUMFLOORS; floor_i++ {
		Clear_orders_in_floor(floor_i)
	}
}

func Set_direction(new_direction int) {
	direction = new_direction
}

func Set_floor(floor int) {
	last_floor = floor
}

func Get_cost(destination int, button_type int) (cost int) {
	var next_floor int

	if State == "running" && direction == config.UP {
		next_floor = last_floor + 1
		if destination < next_floor {
			cost += config.DIRECTION_CHANGE_COST
		}
	} else if State == "running" && direction == config.DOWN {
		next_floor = last_floor - 1
		if destination > next_floor {
			cost += config.DIRECTION_CHANGE_COST
		}
	} else if State == "idle" {
		next_floor = last_floor
		cost += config.STARTUP_FROM_IDLE_COST
	} else if State == "stuck" {
		next_floor = last_floor
		cost += 10000
	}
	cost += config.DISTANCE_COST * int(math.Abs(float64(destination-last_floor)))

	return cost
}

func Get_next() (next_order_at_floor int) {
	var iterator_dir int
	var button_type_i int

	if direction == config.UP {
		iterator_dir = 1
	} else {
		iterator_dir = -1
	}
	var endpoints [2]int
	if iterator_dir == 1 {
		endpoints = [2]int{1, config.NUMFLOORS}
	} else {
		endpoints = [2]int{config.NUMFLOORS, 1}
	}

	//Iterates from nextfloor to end in last direction, then from end to end in opposite direction, then back to, but not including, start
	for floor_i := last_floor; floor_i != endpoints[1]+iterator_dir; floor_i += iterator_dir {
		button_type_i = config.INTERNAL
		if order_matrix[floor_i-1][button_type_i] == self || order_matrix[floor_i-1][button_type_i] == NO_EXECUTER {
			return floor_i
		}
		button_type_i = direction
		if order_matrix[floor_i-1][button_type_i] == self || order_matrix[floor_i-1][button_type_i] == NO_EXECUTER {
			return floor_i
		}
	}
	button_type_i = direction
	button_type_i += iterator_dir //Swap direction of button type

	for floor_i := endpoints[1]; floor_i != endpoints[0]-iterator_dir; floor_i -= iterator_dir {
		if order_matrix[floor_i-1][button_type_i] == self {
			return floor_i
		}
	}
	for floor_i := endpoints[0]; floor_i != last_floor; floor_i += iterator_dir {
		button_type_i = config.INTERNAL
		if order_matrix[floor_i-1][button_type_i] == self {
			return floor_i
		}
		button_type_i = direction
		if order_matrix[floor_i-1][button_type_i] == self {
			return floor_i
		}
	}
	button_type_i = direction
	//Iterates from end to end in opposite direction, then back to, but not including, start
	button_type_i += iterator_dir //Swap direction of button type

	for floor_i := endpoints[1]; floor_i != endpoints[0]-iterator_dir; floor_i -= iterator_dir {
		if order_matrix[floor_i-1][button_type_i] == NO_EXECUTER {
			return floor_i
		}
	}
	for floor_i := endpoints[0]; floor_i != last_floor; floor_i += iterator_dir {
		button_type_i = direction
		if order_matrix[floor_i-1][button_type_i] == NO_EXECUTER {
			return floor_i
		}
	}
	return NO_ORDER
}

func Get_order_matrix() [config.NUMFLOORS][config.NUMBUTTON_TYPES]int {
	return underlying_order_matrix
}

func Merge_external_order_matrix_with_current(new_order_matrix [config.NUMFLOORS][config.NUMBUTTON_TYPES]int) {
	for floor_i := 0; floor_i < config.NUMFLOORS; floor_i++ {
		for button_type_i := 0; button_type_i < config.NUMBUTTON_TYPES; button_type_i++ {
			if order_matrix[floor_i][button_type_i] < new_order_matrix[floor_i][button_type_i] {
				order_matrix[floor_i][button_type_i] = new_order_matrix[floor_i][button_type_i]
			}
		}
	}
}

func Unassign_orders_handled_by(id_of_stuck_elevator int) {
	for i := 0; i < config.NUMFLOORS; i++ {
		for j := 1; j < config.NUMBUTTON_TYPES; j++ {
			if order_matrix[i][j] == id_of_stuck_elevator {
				order_matrix[i][j] = NO_EXECUTER
			}
		}
	}
}

func order_is_in_bounds(destination int, button_type int) bool {
	return (destination <= config.NUMFLOORS) && (destination > 0) && (button_type >= 0) && (button_type < config.NUMBUTTON_TYPES)
}

func Already_exists(destination int, button_type int) bool {
	if order_matrix[destination-1][button_type] == NO_ORDER {
		return false
	}
	return true
}
