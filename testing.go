package main

import (
	"./config"
	"./order_handling"
	"fmt"
)

func main() {

	order := order_handling.Order{2, config.UP, ""}
	order_handling.New_floor_reached(3)
	cost := order_handling.Get_cost(order)

	fmt.Printf("Cost: %d\n", cost)

	order_handling.Insert(order)
	external_order_list := order_handling.Get_list()

	for i, v := range external_order_list {
		fmt.Printf("Order %d \t Destination: %d \t Button_type: %d\n",
			i+1, v.Destination, v.Button_type)
	}

}
