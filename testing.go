package main

import (
	"./config"
	"./order_handling"
	"fmt"
)

func main() {
	fmt.Println("hei\n")
	order_handling.Init()
	order_handling.Insert(3, config.UP)
	order_handling.Print_order_array()
	order_handling.Insert(0, config.DOWN)
	order_handling.Insert(2, -1)
	order_handling.Assign_order_executer(3, config.UP, 12345)
	order_handling.Print_order_array()

}
