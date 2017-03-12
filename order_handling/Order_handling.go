package order_handling

import (
	"../config"
	"../network"
	"fmt"
	"math"
)

/*Her har jeg initialisert direction = config.DOWN, dvs at den ved initialisering
vil kjøre NEDOVER til nærmeste etasje. Når den er blitt initialisert, er det viktig å
kjøre funksjonen arrived at floor, slik at man ikke ender opp med siste kjøreretning
nedover i første etasje, det kan by på problemer. En annen måte å gjøre det på er å
ta inn første etasjen man når i init-funksjonen som parameter*/

const (
	NO_EXECUTER = 0
	NO_ORDER    = -1
)

var underlying_order_matrix [config.NUMFLOORS][config.NUMBUTTON_TYPES]int
var order_matrix [][config.NUMBUTTON_TYPES]int // Matrix on the form [floor][button_type]

var last_floor int
var direction int

var self int

func Init(self_id int) {
	order_matrix = underlying_order_matrix[:][:]
	for i := 0; i < config.NUMFLOORS; i++ {
		for j := 0; j < config.NUMBUTTON_TYPES; j++ {
			order_matrix[i][j] = -1
		}
	}
	self = self_id
	last_floor = 2
	direction = config.DOWN
}

func Merge_order_matrices(new_order_matrix [config.NUMFLOORS][config.NUMBUTTON_TYPES]int) {
	for i := 0; i < config.NUMFLOORS; i++ {
		for j := 1; j < config.NUMBUTTON_TYPES; j++ {
			if order_matrix[i][j] < new_order_matrix[i][j] {
				order_matrix[i][j] = new_order_matrix[i][j]
			}
		}
	}
}

func Insert(destination int, button_type int, executer_id int) bool {

	if order_is_valid(destination, button_type) {
		order_matrix[destination-1][button_type] = executer_id
		return true
	}
	return false
}

func Set_direction(new_direction int) {
	direction = new_direction
}

func Print_order_matrix() {

	for i := 0; i < config.NUMFLOORS; i++ {
		for j := 0; j < 3; j++ {
			str := ""
			str += fmt.Sprintf("Floor: %d ", i+1)

			if j == config.INTERNAL {
				str += "INT  "
			} else if j == config.UP {
				str += "UP   "
			} else {
				str += "DOWN "
			}
			if underlying_order_matrix[i][j] == NO_EXECUTER {
				str += "Order without executer\n"
			} else if underlying_order_matrix[i][j] == NO_ORDER {
				str += "No order\n"
			} else {
				str += fmt.Sprintf("%d\n", underlying_order_matrix[i][j])
			}
			fmt.Printf(str)
		}
	}
}

func Print_order_struct(order_struct network.Orders) {

	for i := 0; i < config.NUMFLOORS; i++ {
		for j := 0; j < 3; j++ {
			str := ""
			str += fmt.Sprintf("Floor: %d ", i+1)

			if j == config.INTERNAL {
				str += "INT  "
			} else if j == config.UP {
				str += "UP   "
			} else {
				str += "DOWN "
			}
			if order_struct.Orders[i][j] == NO_EXECUTER {
				str += "Order without executer\n"
			} else if order_struct.Orders[i][j] == NO_ORDER {
				str += "No order\n"
			} else {
				str += fmt.Sprintf("%d\n", order_struct.Orders[i][j])
			}
			fmt.Printf(str)
		}
	}
}

func Get_cost(destination int, button_type int, state string) (cost int) {
	var next_floor int
	if state == "running" && direction == config.UP {
		next_floor = last_floor + 1
		if destination < next_floor {
			cost += config.DIRECTION_CHANGE_COST
		}
	} else if state == "running" && direction == config.DOWN {
		next_floor = last_floor - 1
		if destination > next_floor {
			cost += config.DIRECTION_CHANGE_COST
		}
	} else {
		next_floor = last_floor
		cost += config.STARTUP_FROM_IDLE_COST
	}
	cost += int(math.Abs(float64(destination - last_floor)))
	return cost
}

func Get_next(state string) (next_order_at_floor int) {

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
	return 0
}

func New_floor_reached(floor int) bool {
	last_floor = floor
	if (order_matrix[floor-1][direction] == self) || (order_matrix[floor-1][config.INTERNAL] == self) {
		return true
	} else if order_matrix[floor-1][direction] == NO_EXECUTER {
		return true
	} else {
		return false
	}
}

func Storage_test() {

}

func Get_order_matrix() [config.NUMFLOORS][config.NUMBUTTON_TYPES]int {
	return underlying_order_matrix
}

func Clear_order(destination int) {
	for button_type_i := 0; button_type_i < config.NUMBUTTON_TYPES; button_type_i++ {
		order_matrix[destination-1][button_type_i] = NO_ORDER
	}
}

func Clear_order_matrix() {
	for i := 0; i < config.NUMFLOORS; i++ {
		Clear_order(i)
	}
}

func Already_exists(destination int, button_type int) bool {
	if order_matrix[destination-1][button_type] == NO_ORDER {
		return false
	}
	return true
}

func order_is_valid(destination int, button_type int) bool {

	if (destination > config.NUMFLOORS) || (destination < 1) {
		fmt.Printf("order_handling.FLOOR_ERROR, selected floor is %d, out of range\n", destination)
		return false
	}

	if (button_type > config.NUMBUTTON_TYPES-1) || (button_type < 0) {
		fmt.Printf("order_handling.BUTTON_TYPE_ERROR:\nselected button type is %d, out of range\n", button_type)
		return false
	}

	if (destination == config.NUMFLOORS) && (button_type == config.UP) {
		fmt.Printf("order_handling.ORDER_ERROR\nInvalid order, requested floor: NUMFLOORS, UP\n")
		return false
	}

	if (destination == 1) && (button_type == config.DOWN) {
		fmt.Printf("order_handling.ORDER_ERROR\nInvalid order, requested floor: 1, DOWN\n")
		return false
	}

	return true
}
