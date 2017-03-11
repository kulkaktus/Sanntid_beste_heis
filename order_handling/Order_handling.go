package order_handling

import (
	"../config"
	//"../network"
	"fmt"
	"math"
)

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
	fmt.Printf("Length: %d\n", len(order_list))
	self = self_id
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

/*
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
*/
func Get_cost(destination int, button_type int, state string) (order_cost int) {

	order_cost = 0
	distance := 0
	instance := 0
	if state == "idle" {
		if destination == last_floor {
			if (button_type == last_direction) || (button_type == config.INTERNAL) {
				return 0
			} else if button_type != last_direction {
				if button_type == config.UP {
					instance = 3
				} else if button_type == config.DOWN {
					instance = 2
				}
			}
		}
	} else {
		//Calculating scores for internal orders
		if button_type == config.INTERNAL {

			if last_direction == config.UP {
				if destination > last_floor {
					instance = 1

				} else if destination <= last_floor {
					instance = 3
				}

			} else if last_direction == config.DOWN {
				if destination < last_floor {
					instance = 1
				} else if destination > last_floor {
					instance = 2
				}
			}

		} else if last_direction == config.UP {
			if (destination > last_floor) && ((button_type == last_direction) || (destination == config.NUMFLOORS)) {
				instance = 1
			} else if button_type != last_direction {
				instance = 3
			} else if (destination <= last_floor) && (button_type == last_direction) {
				instance = 4
			}

		} else if last_direction == config.DOWN {
			if (destination < last_floor) && ((button_type == last_direction) || (destination == 1)) {
				instance = 1
			} else if button_type != last_direction {
				instance = 2
			} else if (destination >= last_floor) && (button_type == last_direction) {
				instance = 5
			}
		}
	}

	switch instance {
	case 1:
		distance += int(math.Abs(float64(destination - last_floor)))
	case 2:
		distance += last_floor + destination
	case 3:
		distance += 2*config.NUMFLOORS - last_floor - destination
	case 4:
		distance += 2*config.NUMFLOORS + destination - last_floor - 1
	case 5:
		distance += 2*config.NUMFLOORS + last_floor - destination - 1
	}

	order_cost += config.DISTANCE_COST * distance

	/*stops_inbetween := 0

	if button_type ==

	if last_floor < destination {

		for _, floor_i := range external_order_list {
			if v.Executer == self {
				if (v.Destination > last_floor) && (v.Destination < order.Destination) {
					stops_inbetween += 1
				}
				}
		}
	}
	if last_floor > order.Destination {

		for _, v := range external_order_list {
			if v.Executer == self {
					if (v.Destination < last_floor) && (v.Destination > order.Destination) {
						stops_inbetween += 1
					}

				}
			}

		}
	}
	order_cost += config.STOPS_INBETWEEN_COST * stops_inbetween

	*/
	return order_cost
}

func Assign_order_executer(floor int, button_type int, executer_id int) {
	if order_is_valid(floor, button_type) {
		order_list[floor-1][button_type] = executer_id
	}
}

func Get_next(state string) (destination int, button_type int) {

	cost := 1000
	temp_lowest_cost := cost
	destination, button_type = 0, 0

	for i := 0; i < config.NUMFLOORS; i++ {
		for j := 0; j < NUMBUTTON_TYPES; j++ {
			if order_list[i][j] == self {
				cost = Get_cost(i, j, state)
				if cost < temp_lowest_cost {
					temp_lowest_cost = cost
					destination, button_type = i, j
				}
			}
		}
	}

	if destination == 0 {
		for i := 0; i < config.NUMFLOORS; i++ {
			for j := 0; j < NUMBUTTON_TYPES; j++ {
				if order_list[i][j] == ORDER_WITHOUT_EXECUTER {
					cost = Get_cost(i, j, state)
					if cost < temp_lowest_cost {
						temp_lowest_cost = cost
						destination, button_type = i, j
					}
				}
			}
		}
		order_list[destination][button_type] = self
	}

	return destination, button_type
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
		fmt.Printf("order_handling.BUTTON_TYPE_ERROR, selected button type is %d, out of range\n", button_type)
		return false
	}

	return true
}
