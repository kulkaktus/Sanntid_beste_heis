package fsm

import (
	"../config"
	"../io/buttons"
	//"../io/io"
	"../io/lights"
	"../io/motor"
	"../io/sensors"
	"../network"
	//"../peers"
	//"../order_handling"
	"fmt"
	//"os"
	//"time"
)

var floor int

func Fsm(id int, tx chan<- network.Message, rx <-chan network.Message) {
	floor = 0
	fmt.Println("Started")
	state := "idle"
	check_network(id, tx, rx)
	for {
		//fmt.Println("Loop\n")
		switch state {
		case "idle":
			check_io(id, tx)
			//check_network(id, tx, rx)
		case "running":
			check_io(id, tx)
			//check_network(id, tx, rx)
		}
	}
}

func check_io(id int, tx chan<- network.Message) (next_state string) {
	for button_type_i := 0; button_type_i <= config.DOWN; button_type_i++ {
		for floor_i := 1; floor_i <= config.NUMFLOORS; floor_i++ {
			if buttons.Get(button_type_i, floor_i) {
				lights.Set(button_type_i, floor_i)
				motor.Go(button_type_i)
				next_state = "running"
				/*go func() {
					order := order_handling.Order{floor_i, button_type_i, ""}
					if order_handling.Insert(order) {
						//network.Send_order(order_handling.Get_cost(order), order, tx)
					}
				}()*/
			} else {
				lights.Clear(button_type_i, floor_i)
			}
		}
	}
	if floor != sensors.Get() && sensors.Get() != 0 {
		motor.Stop()
		next_state = "idle"
		floor = sensors.Get()
		lights.Set(config.INDICATE, floor)
		fmt.Printf("Arrived at %d \n", floor)
	}
	return
}

func check_network(id int, tx chan<- network.Message, rx <-chan network.Message) int {
	for {
		select {
		case p := <-network.PeerUpdateCh:
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %q\n", p.Peers)
			fmt.Printf("  New:      %q\n", p.New)
			fmt.Printf("  Lost:     %q\n", p.Lost)

		case a := <-rx:
			if a.Id != id {
				fmt.Printf("Received: %#v\n", a.Msg)
			}
		}
		return 0
	}
}
