package fsm

import (
	"../config"
	"../io/buttons"
	//"../io/io"
	"../io/lights"
	"../io/motor"
	"../io/sensors"
	"../network"
	//"../network/peers"
	"../order_handling"
	"fmt"
	//"os"
	"time"
)

var floor int
var peers []int
var peers_internal_orders map[int][config.NUMFLOORS]int

const (
	tries_to_send   = 5
	time_to_respond = 50 * time.Millisecond

	doors_open_for = 3 * time.Second
)

func Fsm(id int, ordersTx chan<- network.Orders, ordersRx <-chan network.Orders, updateTx chan<- network.Update, updateRx <-chan network.Update, messageTx chan<- network.Message, messageRx <-chan network.Message) {
	doors_open_since := time.Hour
	state := "idle"
	peers_internal_orders = make(map[int][config.NUMFLOORS]int)
	score_responseRx := make(chan [2]int)
	orders_responseRx := make(chan int)
	motor.Go(config.DOWN)
	for sensors.Get() == 0 {
	}
	motor.Stop()
	floor = sensors.Get()
	fmt.Printf("Started at floor %d\n", floor)
	go check_network(id, ordersTx, ordersRx, updateTx, updateRx, messageTx, messageRx, score_responseRx, orders_responseRx)
	for {
		check_buttons_and_update_orders(id, updateTx)

		current_order_floor := order_handling.Get_next(state)
		if current_order_floor == sensors.Get() {
			doors_open_since = 0
			motor.Stop()
			state = "door_open"
		}
		switch state {
		case "idle":

		case "running":

		case "door_open":
			if doors_open_since < doors_open_for {
				lights.Set(config.DOOR, 0)
			}

		case "":

		}
	}
}

func check_buttons_and_update_orders(id int, updateTx chan<- network.Update) {
	for button_type_i := 0; button_type_i <= config.DOWN; button_type_i++ {
		for floor_i := 1; floor_i <= config.NUMFLOORS; floor_i++ {
			if buttons.Get(button_type_i, floor_i) {
				lights.Set(button_type_i, floor_i)
				if !order_handling.Already_exists(floor_i, button_type_i) {
					fmt.Println("orderinserted")
					order_handling.Print_order_array()
					order_handling.Insert(floor_i, button_type_i)
					updateTx <- network.Update{floor_i, button_type_i, 0, id}
				}
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
	/*if floor != sensors.Get() && sensors.Get() != 0 {
		motor.Stop()
		floor = sensors.Get()
		lights.Set(config.INDICATE, floor)
		fmt.Printf("Arrived at %d \n", floor)
	}*/
}

func check_network(id int, ordersTx chan<- network.Orders, ordersRx <-chan network.Orders, updateTx chan<- network.Update, updateRx <-chan network.Update, messageTx chan<- network.Message, messageRx <-chan network.Message, score_responseRx chan<- [2]int, orders_responseRx chan<- int) {
	for {
		select {
		case p := <-network.PeerUpdateCh:
			peers = p.Peers
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %d\n", p.Peers)
			fmt.Printf("  New:      %d\n", p.New)
			fmt.Printf("  Lost:     %d\n", p.Lost)
		case b := <-ordersRx:
			order_handling.Print_order_struct(b)
		case d := <-updateRx:
			order_handling.Assign_order_executer(d.Floor, 0, d.Executer)
			str := fmt.Sprintf("Order update at floor %d of type ", d.Floor)
			if d.Button_type == config.INTERNAL {
				str += "INT  "
				temp_internal_order := peers_internal_orders[d.Executer]
				temp_internal_order[d.Floor] = d.Executer
				peers_internal_orders[d.Executer] = temp_internal_order
			} else if d.Button_type == config.UP {
				str += "UP   "
				order_handling.Assign_order_executer(d.Floor, d.Button_type, d.Executer)
			} else {
				str += "DOWN "
			}
			if d.Executer == -1 {
				str += "Order is without executer\n"
			} else if d.Executer == 0 {
				str += "Order cleared\n"
			} else {
				str += fmt.Sprintf("Order handled by: %d\n", d.Executer)
			}
			fmt.Printf(str)
		case f := <-messageRx:
			if f.To_id == id {
				fmt.Printf("Received: %d\n", f.Content)
				if f.Type == network.SCORE_RESPONSE_T {
					score_responseRx <- [2]int{f.From_id, f.Content}
				} else if f.Type == network.SCORE_RESPONSE_T {
					orders_responseRx <- f.From_id
				}
			}
		}
	}
}

func update_order(id int, button_type int, floor int, handler int, updateTx chan<- network.Update, score_responseRx <-chan [2]int) {
	for i := 0; i < tries_to_send; i++ {
		updateTx <- network.Update{floor, button_type, handler, id}
		select {
		case a := <-score_responseRx:
			fmt.Printf("Received score of %d, from %d \n", a[0], a[1])
		case <-time.After(time_to_respond):
			break
		}
	}
}
