package order_handling

import (
	"./../config"
	"fmt"
	"math"
)

type Order struct {
	Destination int
	Button_type int
	Executer    string
}

var external_order_list []Order
var internal_order_list []Order

var last_floor int
var last_direction int

var self string
var no_executer string

func Init() {
	//
	self = "dummy"
	no_executer = ""

}

func Insert(order Order) bool {

	if order.Button_type == config.INTERNAL {
		order_list := internal_order_list[:]

		if Already_exists(order_list, order) {
			return false
		} else {
			internal_order_list = append(order_list, order)
			return true
		}

	} else if order.Button_type == config.UP || order.Button_type == config.DOWN {
		order_list := external_order_list[:]

		if Already_exists(order_list, order) {
			return false
		} else {
			external_order_list = append(order_list, order)
			return true
		}

	} else {
		fmt.Printf("Button type ERROR, value is %d", order.Button_type)
		return false
	}
}

func Get_cost(order Order, state string) int {

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
}

func Assign_elevator() {

}

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

func Get_list() []Order {

	deep_copy := make([]Order, len(external_order_list))

	for i, v := range external_order_list {
		deep_copy[i] = v
	}
	return deep_copy
}

func Already_exists(order_list []Order, order Order) bool {
	for _, v := range order_list {
		if v.Destination == order.Destination && v.Button_type == order.Button_type {
			return true
		}

	}
	return false
}
