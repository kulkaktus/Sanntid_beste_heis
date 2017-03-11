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

	fmt.Printf("Cost of 2 UP: %d\n", order_handling.Get_cost(2, config.UP, "idle"))
	fmt.Printf("Cost of 3 UP: %d\n", order_handling.Get_cost(3, config.UP, "idle"))

	order_handling.Print_order_array()

}
