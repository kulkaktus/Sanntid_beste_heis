package order_handling

import (
	"../config"
	"../network"
	"fmt"
	//"math"
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

func Init() {
	order_list = underlying_order_array[:][:]
	/*for i := 0; i < config.NUMFLOORS; i++ {
		new_slice := make([]int, 3)
		order_list = append(order_list, new_slice)

	}*/
	fmt.Printf("Length: %d\n", len(order_list))
}

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

/*func Get_cost(floor int, button_type int, state string) int {

	order_cost := 0
	distance := int(math.Abs(float64(last_floor-order.Destination)))
	order_cost += config.DISTANCE_COST * distance

	stops_inbetween := 0
	if last_floor < order.Destination {

		for _, v := range external_order_list {
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
	order_cost += config.STOPS_INBETWEEN_COST * stops_inbetween


	direction_changes := 0

	if (state == "idle") && (order.Destination == last_floor){
		if (order.Button_type == config.INTERNAL) || (order.Button_type == last_direction){
			return 0
		} else{
			direction_changes += 1
		}
	} else {
		if last_direction == config.UP {
			if order.Button_type == config.INTERNAL  {
				if order.Destination < last_floor {
					direction_changes += 1
				}

			} else {

				//Checking orders in order of priority
				for i := last_floor + 1; i <= config.NUMFLOORS - 1; i++ {
					if (order.Destination > last_floor) && (order.Button_type == config.UP)


				}

			}
		}
	}



	return order_cost
}*/

func Assign_order_executer(floor int, button_type int, executer_id int) {
	if order_is_valid(floor, button_type) {
		order_list[floor-1][button_type] = executer_id
	}
}

/*
func Get_next(state string) (next_order Order) {

	temp_lowest_cost = 1000

	for _, ext_order_i := range external_order_list {
		if Get_cost(ext_order_i, state) < temp_lowest_cost {
			next_order = ext_order_i
		}
	}

	for _, int_order_i := range internal_order_list {
		if Get_cost(int_order_i, state) < temp_lowest_cost {
			next_order = int_order_i
		}
	}

	return next_order
}
*/
func New_floor_reached(floor int) {
	last_floor = floor
	if floor == 0 {
		last_direction = config.UP
	} else if floor == config.NUMFLOORS {
		last_direction = config.DOWN

	}
}

func Storage_test() {

}

func Get_list() [config.NUMFLOORS][NUMBUTTON_TYPES]int {
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
