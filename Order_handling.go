package main

import (
	"fmt"
)

const Num_floors = 4
const Num_elevators = 3

//const Ext_order_list_size = 2 * Num_floors * Num_elevators
//const Int_order_list_size = Num_elevators

const (
	Internal = iota
	Up       = iota
	Down     = iota
)

const (
	no_executor = -1
	self = 0

type Order struct {
	Floor       int //Rename to Destination for clarity
	Button_type int
	Executer    int
}

var external_order_list []Order
var internal_order_list []Order

var last_floor int
var direction int

func main() {

	for i := 0; i < 10; i++ {
		_ = Order_insert(Order{i, Up, -1})

	}
	

	for i, v := range external_order_list {
		fmt.Printf("Order %d \t Floor: %d \t Button_type: %d\n",
			i+1, v.Floor, v.Button_type)
	}

}

func Order_init() {

}

func Order_insert(order Order) bool {
	if order.Button_type == Internal {
		order_list := internal_order_list[:]

		if Order_already_exists(order_list, order) {
			return false
		} else {
			internal_order_list = append(order_list, order)
			return true
		}
	} else if order.Button_type == Up || order.Button_type == Down {
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

func Order_get_score(order Order) int {
	return 1
}

func Order_assign_elevator() {

}

func Order_get_next() next_order Order {

	next_order := Order{0,0, no_executor}
	temp := Order{1000, Up, no_executor}
	

	if (current_floor == Num_floors){

	}

	if (direction == Up) {

		for i := last_floor + 1; i < Num_floors; i++ { //will not go out of bounds here, since direction == down if at top floor, see Order_new_floor_reached

			for j, order := range internal_order_list{
				if order.Floor == i {
					return order
				}
			}

			for k, order := range external_order_list{
				if order.Executor == self{
					if order.Button_type == Up {
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
		direction = Up
	}
	else if floor == Num_floors{
		direction = Down
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
		if v.Floor == order.Floor && v.Button_type == order.Button_type {
			return true
		}

	}
	return false
}
