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
var state string

const (
	tries_to_send   = 5
	time_to_respond = 50 * time.Millisecond

	doors_open_for = 1 * time.Second
)

func Fsm(id int, ordersTx chan<- network.Orders, ordersRx <-chan network.Orders, updateTx chan<- network.Update, updateRx <-chan network.Update, messageTx chan<- network.Message, messageRx <-chan network.Message) {
	doors_open_since := time.Now().Add(-time.Hour)
	state = "idle"
	peers_internal_orders = make(map[int][config.NUMFLOORS]int)
	score_responseRx := make(chan [2]int)
	orders_responseRx := make(chan int)
	motor.Go(config.DOWN)
	for sensors.Get() == 0 {
	}
	motor.Stop()
	floor = sensors.Get()
	fmt.Printf("Started at floor %d\n", floor)
	order_handling.New_floor_reached(floor)
	go check_network(id, ordersTx, ordersRx, updateTx, updateRx, messageTx, messageRx, score_responseRx, orders_responseRx)
	for {
		check_buttons_and_update_orders(id, updateTx)
		order_handling.State = state
		//current_order_floor := order_handling.Get_next(state)
		switch state {
		case "idle":
			current_order_floor := order_handling.Get_next(state)
			if current_order_floor == 0 {
				motor.Stop()
				state = "idle"
			} else if current_order_floor < floor {
				motor.Go(config.DOWN)
				order_handling.Set_direction(config.DOWN)
				state = "running"
			} else if current_order_floor > floor {
				motor.Go(config.UP)
				order_handling.Set_direction(config.UP)
				state = "running"
			} else if current_order_floor == floor {
				motor.Stop()
				doors_open_since = time.Now()
				state = "door_open"
				order_handling.Clear_order(floor)
			}

		case "running":
			sensor := sensors.Get()
			if sensor != 0 && sensor != floor {
				floor = sensor
				motor.Stop()
				lights.Set(config.INDICATE, floor)
				lights.Clear_floor(floor)
				fmt.Printf("Arrived at %d \n", floor)
				state = "idle"
				order_handling.New_floor_reached(floor)

			}
		case "door_open":
			//fmt.Println(time.Since(doors_open_since))
			if time.Since(doors_open_since) < doors_open_for {
				lights.Set(config.DOOR, 0)
			} else {
				lights.Clear(config.DOOR, 0)
				state = "idle"
			}

		case "":
			panic("No state in fsm")
		}
	}
}

func check_buttons_and_update_orders(id int, updateTx chan<- network.Update) {
	for button_type_i := 0; button_type_i <= config.DOWN; button_type_i++ {
		for floor_i := 1; floor_i <= config.NUMFLOORS; floor_i++ {
			if buttons.Get(button_type_i, floor_i) {
				lights.Set(button_type_i, floor_i)
				if !order_handling.Already_exists(floor_i, button_type_i) {

					if button_type_i == config.INTERNAL {
						if order_handling.Insert(floor_i, button_type_i, id) {
							fmt.Println("orderinserted")
						}
					} else {
						order_handling.Insert(floor_i, button_type_i, order_handling.NO_EXECUTER)
					}
					order_handling.Print_order_matrix()
					updateTx <- network.Update{floor_i, button_type_i, order_handling.NO_EXECUTER, id}
				}
			}
		}
	}
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
			if b.From_id != id {
				order_handling.Print_external_order_matrix(b.Orders)
				order_handling.Merge_order_matrices(b.Orders)
				var temp_internal_orders [4]int
				for i := 0; i < config.NUMFLOORS; i++ {
					temp_internal_orders[i] = b.Orders[i][0]
				}
				peers_internal_orders[b.From_id] = temp_internal_orders
				messageTx <- network.Message{b.From_id, id, network.ORDERS_RESPONSE_T, 0}
			}
		case d := <-updateRx:

			str := fmt.Sprintf("Order update at floor %d of type ", d.Floor)
			if d.From_id != id {
				order_handling.Insert(d.Floor, d.Button_type, d.Executer)
				if d.Button_type == config.INTERNAL {
					str += "INT  "
					temp_internal_order := peers_internal_orders[d.Executer]
					temp_internal_order[d.Floor-1] = d.Executer
					peers_internal_orders[d.Executer] = temp_internal_order
				} else if d.Button_type == config.UP {
					str += "UP   "
					order_handling.Insert(d.Floor, d.Button_type, d.Executer)
				} else {
					str += "DOWN "
					order_handling.Insert(d.Floor, d.Button_type, d.Executer)
				}
				if d.Executer == order_handling.NO_EXECUTER {
					str += "Order is without executer\n"
				} else if d.Executer == order_handling.NO_ORDER {
					str += "Order cleared\n"
				} else {
					str += fmt.Sprintf("Order handled by: %d\n", d.Executer)
				}
				if d.From_id != id {
					if order_handling.Insert(d.Floor, d.Button_type, d.Executer) {
						messageTx <- network.Message{d.From_id, id, network.ORDERS_RESPONSE_T, order_handling.Get_cost(d.Floor, d.Button_type)}
					}
				}
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

func send_orders(id int, to_id int, orders [config.NUMFLOORS][config.NUMBUTTON_TYPES]int, ordersTx chan<- network.Orders, orders_responseRx <-chan int) {
	for i, v := range peers_internal_orders[to_id] {
		orders[i][0] = v
	}

	for i := 0; i < tries_to_send; i++ {
		ordersTx <- network.Orders{orders, id}
		select {
		case a := <-orders_responseRx:
			fmt.Printf("Received ack of orders from %d\n", a)
		case <-time.After(time_to_respond):
			break
		}
	}
}
