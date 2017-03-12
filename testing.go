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
	fmt.Println("********TESTING CASE 1, UPWARDS\n*********")
	order_handling.New_floor_reached(1)
	order_handling.New_floor_reached(2)
	fmt.Printf("Start at floor 2, last_dir up \n")
	fmt.Printf("Cost of 2 UP: %d\n", order_handling.Get_cost(2, config.UP, "running"))
	fmt.Printf("Cost of 1 UP: %d\n", order_handling.Get_cost(1, config.UP, "running"))

	order_handling.Print_order_array()

}
