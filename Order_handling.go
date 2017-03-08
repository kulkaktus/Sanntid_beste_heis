package main

import (
	"fmt"
	"config"
	"math"
)


const (
	no_executor = -1
	self = 0

type Order struct {
	Destination       int //Rename to Destination for clarity
	Button_type int
	Executer    int
}

var external_order_list []Order
var internal_order_list []Order

var last_floor int
var last_direction int

func main() {

	for i := 0; i < 10; i++ {
		_ = Order_insert(Order{i, config.UP, -1})

	}


	for i, v := range external_order_list {
		fmt.Printf("Order %d \t Destination: %d \t Button_type: %d\n",
			i+1, v.Destination, v.Button_type)
	}

}

func Order_init() {

}

func Order_insert(order Order) bool {

	if order.Button_type == config.INSIDE {
		order_list := internal_order_list[:]

		if Order_already_exists(order_list, order) {
			return false
		} else {
			internal_order_list = append(order_list, order)
			return true
		}

	} else if order.Button_type == config.UP || order.Button_type == config.DOWN {
		order_list := external_order_list[:]

		if Order_already_exists(order_list, order) {
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

func Order_get_cost(order Order) int {

	order_cost := 0
	order_cost += config.DISTANCE_COST * abs(last_floor - order.Destination)

	if order.Button_type != last_direction {
		direction += config.OPPOSITE_DIRECTION_COST
	}

	orders_inbetween := 0
	for i, v := range external_order_list {

		if v.Executer == self {
			if v.Destination >

		}

	}



	return 1
}

func Order_assign_elevator() {

}

func Order_get_next() next_order Order {

	next_order := Order{0,0, no_executor}
	temp := Order{1000, config.UP, no_executor}
	
	if (current_floor == config.NUMFLOORS){

	}

	if (direction == config.UP) {

		for i := last_floor + 1; i < config.NUMFLOORS; i++ { //will not go out of bounds here, since direction == down if at top floor, see Order_new_floor_reached

			for j, order := range internal_order_list{
				if order.Destination == i {
					return order
				}
			}

			for k, order := range external_order_list{
				if order.Executor == self{
					if order.Button_type == config.UP {
						return order
					}
				}
			}

		}
	}



}

func Order_new_floor_reached(floor int) {
	last_floor = floor
	if floor == 0 {
		direction = config.UP
	}
	else if floor == config.NUMFLOORS{
		direction = config.DOWN 

	}
}

func Order_storage_test() {

}

func Order_list_get() []Order {

	deep_copy := make([]Order, len(external_order_list))

	for i, v := range external_order_list {
		deep_copy[i] = v
	}
	return deep_copy
}

func Order_already_exists(order_list []Order, order Order) bool {
	for _, v := range order_list {
		if v.Destination == order.Destination && v.Button_type == order.Button_type {
			return true
		}

	}
	return false
}
