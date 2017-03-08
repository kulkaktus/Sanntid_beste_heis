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

func Get_cost(order Order) int {

	order_cost := 0
	order_cost += int(math.Abs(float64(last_floor-order.Destination))) * config.DISTANCE_COST

	if order.Button_type != last_direction {
		order_cost += config.OPPOSITE_DIRECTION_COST
	}

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

	order_cost += stops_inbetween * config.STOPS_INBETWEEN_COST

	return order_cost
}

func Assign_elevator() {

}

func Get_next() (next_order Order) {

	next_order = Order{0, 0, no_executer}
	//temp := Order{1000, config.UP, no_executer}

	if last_floor == config.NUMFLOORS {

	}

	if last_direction == config.UP {

		for i := last_floor + 1; i < config.NUMFLOORS; i++ { //will not go out of bounds here, since direction == down if at top floor, see Order_new_floor_reached

			for _, order := range internal_order_list {
				if order.Destination == i {
					return order
				}
			}

			for _, order := range external_order_list {
				if order.Executer == self {
					if order.Button_type == config.UP {
						return order
					}
				}
			}

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
