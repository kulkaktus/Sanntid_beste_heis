package main

import (
	"./config"
	"./order_handling"
	"fmt"
)

func main() {
	fmt.Println("hei\n")
	order_handling.Init(123123)
	order_handling.Print_order_array()

	order_handling.New_floor_reached(1)
	order_handling.New_floor_reached(2)
	fmt.Printf("Start at floor 2, last_dir up \n\n")

	order_handling.Insert(1, config.UP)

	order_handling.Assign_order_executer(4, config.DOWN, NO_EXECUTER)
	order_handling.Assign_order_executer(1, config.UP, NO_EXECUTER)

	order_handling.Print_order_array()
	fmt.Printf("Next floor to go to: %d\n\n", order_handling.Get_next("running"))
	order_handling.Clear_order(3, config.UP)
	order_handling.Print_order_array()
	fmt.Printf("Next floorto go to: %d\n", order_handling.Get_next("running"))

	order_handling.Print_order_array()

}
