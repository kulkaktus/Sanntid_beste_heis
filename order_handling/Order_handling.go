package order_handling

import (
	"../config"
	"../network"
	"fmt"
	"math"
)

/*Her har jeg initialisert last_direction = config.DOWN, dvs at den ved initialisering
vil kjøre NEDOVER til nærmeste etasje. Når den er blitt initialisert, er det viktig å
kjøre funksjonen arrived at floor, slik at man ikke ender opp med siste kjøreretning
nedover i første etasje, det kan by på problemer. En annen måte å gjøre det på er å
ta inn første etasjen man når i init-funksjonen som parameter*/

const (
	ORDER_WITHOUT_EXECUTER = -1
	NO_ORDER               = 0
)
const NUMBUTTON_TYPES = 3

var underlying_order_array [config.NUMFLOORS][NUMBUTTON_TYPES]int
var order_list [][NUMBUTTON_TYPES]int // Matrix on the form [floor][button_type]

var last_floor int
var last_direction int

var self int

func Init(self_id int) {
	order_list = underlying_order_array[:][:]
	self = self_id
	last_direction = config.DOWN
}

/*func Merge_order_lists(new_list [config.NUMFLOORS][NUMBUTTON_TYPES]int){
	for i:=0; i<config.NUMFLOORS; i++ {
		for j:=0; j<NUMBUTTON_TYPES; j++ {
			if order_list[i][j] != new_list[i][j] {
				if order_list [i][j] == NO_ORDER {
					order_list[i][j] == new_list[i][j]
				} if (order_list[i][j] == ORDER_WITHOUT_EXECUTER) && (order_list[i][j] != NO_ORDER){
					order_list[i][j] == new_list[i][j]
				}
			}
		}
	}
}*/

func Insert(floor int, button_type int) bool {

	if order_is_valid(floor, button_type) {
		order_list[floor-1][button_type] = ORDER_WITHOUT_EXECUTER
		return true
	}
	return false
}

func Print_order_array() {

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
			if underlying_order_array[i][j] == ORDER_WITHOUT_EXECUTER {
				str += "Order without executer\n"
			} else if underlying_order_array[i][j] == NO_ORDER {
				str += "No order\n"
			} else {
				str += fmt.Sprintf("%d\n", underlying_order_array[i][j])
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
			if order_struct.Orders[i][j] == ORDER_WITHOUT_EXECUTER {
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
	if state == "running" && last_direction == config.UP {
		next_floor = last_floor + 1
		if destination < next_floor {
			cost += config.DIRECTION_CHANGE_COST
		}
	} else if state == "running" && last_direction == config.DOWN {
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

func Assign_order_executer(floor int, button_type int, executer_id int) {
	if order_is_valid(floor, button_type) {
		order_list[floor-1][button_type] = executer_id
	}
}

func Get_next(state string) (destination int, button_type int) {
	var next_floor int
	if state == "running" && last_direction == config.UP {
		next_floor = last_floor + 1
	} else if state == "running" && last_direction == config.DOWN {
		next_floor = last_floor - 1
	} else {
		next_floor = last_floor
	}
	for floor_i := next_floor; floor_i < config.NUMFLOORS+next_floor; floor_i++ {
		for button_type_i := 0; button_type_i < NUMBUTTON_TYPES; button_type_i++ {
			if order_list[floor_i%config.NUMFLOORS][button_type_i] == self {
				return floor_i + 1, button_type_i

			}
		}
	}
	for floor_i := next_floor; floor_i < config.NUMFLOORS+next_floor; floor_i++ {
		for button_type_i := 1; button_type_i < NUMBUTTON_TYPES; button_type_i++ {
			if order_list[floor_i%config.NUMFLOORS][button_type_i] == ORDER_WITHOUT_EXECUTER {
				order_list[floor_i%config.NUMFLOORS][button_type_i] = self
				return floor_i + 1, button_type_i
			}
		}
	}
	return 0, 0
}

func New_floor_reached(floor int) bool {
	last_floor = floor
	if floor == 0 {
		last_direction = config.UP
	} else if floor == config.NUMFLOORS {
		last_direction = config.DOWN
	}

	if (order_list[floor][last_direction] == self) || (order_list[floor][config.INTERNAL] == self) {
		return true
	} else if order_list[floor][last_direction] == ORDER_WITHOUT_EXECUTER {
		return true
	} else {
		return false
	}
}

func Storage_test() {

}

func Get_order_array() [config.NUMFLOORS][NUMBUTTON_TYPES]int {
	return underlying_order_array
}

func Clear_order(destination int, button_type int) {
	order_list[destination][button_type] = 0
}

func Clear_all_orders() {
	for i := 0; i < config.NUMFLOORS; i++ {
		for j := 0; j < NUMBUTTON_TYPES; j++ {
			Clear_order(i, j)
		}
	}
}

func Already_exists(floor int, button_type int) bool {
	if order_list[floor-1][button_type] == 0 {
		return false
	}
	return true
}

func order_is_valid(floor int, button_type int) bool {

	if (floor > config.NUMFLOORS) || (floor < 1) {
		fmt.Printf("order_handling.FLOOR_ERROR, selected floor is %d, out of range\n", floor)
		return false
	}

	if (button_type > NUMBUTTON_TYPES-1) || (button_type < 0) {
		fmt.Printf("order_handling.BUTTON_TYPE_ERROR:\nselected button type is %d, out of range\n", button_type)
		return false
	}

	if (floor == config.NUMFLOORS) && (button_type == config.UP) {
		fmt.Printf("order_handling.ORDER_ERROR\nInvalid order, requested floor: NUMFLOORS, UP\n")
		return false
	}

	if (floor == 1) && (button_type == config.DOWN) {
		fmt.Printf("order_handling.ORDER_ERROR\nInvalid order, requested floor: 1, DOWN\n")
		return false
	}

	return true
}
