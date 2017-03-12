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
	time_to_respond = 200 * time.Millisecond

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
	go message_manager(id, ordersTx, ordersRx, updateTx, updateRx, messageTx, messageRx, score_responseRx, orders_responseRx, orders_responseRx)
	if id == 101 {
		time.Sleep(time.Hour)
	}
	for {
		check_buttons_and_update_orders(id, updateTx, score_responseRx)
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
				lights.Clear_floor(floor)
			}

		case "running":
			sensor := sensors.Get()
			if sensor != 0 && sensor != floor {
				floor = sensor
				motor.Stop()
				lights.Set(config.INDICATE, floor)

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

func check_buttons_and_update_orders(id int, updateTx chan<- network.Update, score_responseRx <-chan [2]int) {
	for button_type_i := 0; button_type_i <= config.DOWN; button_type_i++ {
		for floor_i := 1; floor_i <= config.NUMFLOORS; floor_i++ {
			if buttons.Get(button_type_i, floor_i) {
				if !order_handling.Already_exists(floor_i, button_type_i) {

					if button_type_i == config.INTERNAL {
						if !order_handling.Already_exists(floor_i, button_type_i) {
							if order_handling.Insert(floor_i, button_type_i, id) {
								update_order(id, button_type_i, floor_i, id, updateTx, score_responseRx)
							}
						}
					} else {
						if !order_handling.Already_exists(floor_i, button_type_i) {
							if order_handling.Insert(floor_i, button_type_i, id) {
								new_order(id, button_type_i, floor_i, updateTx, score_responseRx)
							}
						}
					}
					order_handling.Print_order_matrix()
				}
			}
		}
	}
}

func message_manager(id int, ordersTx chan<- network.Orders, ordersRx <-chan network.Orders, updateTx chan<- network.Update, updateRx <-chan network.Update, messageTx chan<- network.Message, messageRx <-chan network.Message, score_responseRx chan<- [2]int, orders_responseRx_out chan<- int, orders_responseRx_in <-chan int) {
	for {
		select {
		case p := <-network.PeerUpdateCh:
			peers = p.Peers
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %d\n", p.Peers)
			fmt.Printf("  New:      %d\n", p.New)
			fmt.Printf("  Lost:     %d\n", p.Lost)
			send_orders(id, p.New, order_handling.Get_order_matrix(), ordersTx, orders_responseRx_in)
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
			var str string
			if d.From_id != id {
				str = fmt.Sprintf("Order update at floor %d of type ", d.Floor)
				if d.Button_type == config.INTERNAL {
					str += "INT  "
					temp_internal_order := peers_internal_orders[d.Executer]
					temp_internal_order[d.Floor-1] = d.Executer
					peers_internal_orders[d.Executer] = temp_internal_order
				} else if d.Button_type == config.UP {
					str += "UP   "
				} else {
					str += "DOWN "
				}
				if d.Executer == order_handling.NO_EXECUTER {
					str += "Order is without executer\n"
				} else if d.Executer == order_handling.NO_ORDER {
					str += "Order cleared\n"
				} else {
					str += fmt.Sprintf("Order handled by: %d\n", d.Executer)
				}
				if order_handling.Insert(d.Floor, d.Button_type, d.Executer) {
					messageTx <- network.Message{d.From_id, id, network.SCORE_RESPONSE_T, order_handling.Get_cost(d.Floor, d.Button_type)}
				}
			}
			fmt.Printf(str)
		case f := <-messageRx:
			if f.To_id == id {
				fmt.Printf("Received: %d\n", f.Content)
				if f.Type == network.SCORE_RESPONSE_T {
					score_responseRx <- [2]int{f.From_id, f.Content}
				} else if f.Type == network.ORDERS_RESPONSE_T {
					orders_responseRx_out <- f.From_id
				}
			}
		}
	}
}

func update_order(id int, button_type int, floor int, handler int, updateTx chan<- network.Update, score_responseRx <-chan [2]int) {
	for i := 0; i < tries_to_send; i++ {
		updateTx <- network.Update{floor, button_type, handler, id}
		pending_peers := make([]int, len(peers))
		copy(pending_peers, peers)
		for len(pending_peers) != 0 {
			select {
			case a := <-score_responseRx:
				fmt.Printf("Received score of %d, from %d \n", a[1], a[0])
				for i := 0; i < len(pending_peers); i++ {
					if pending_peers[i] == a[1] {
						pending_peers = append(pending_peers[:i], pending_peers[i+1:]...)
					}
				}
			case <-time.After(time_to_respond):
				return
			}
		}
		return
	}
}

func new_order(id int, button_type int, floor int, updateTx chan<- network.Update, score_responseRx <-chan [2]int) {
	for i := 0; i < tries_to_send; i++ {
		updateTx <- network.Update{floor, button_type, order_handling.NO_EXECUTER, id}
		lowest_cost := order_handling.Get_cost(floor, button_type)
		has_lowestcost := id
		pending_peers := make([]int, len(peers))
		copy(pending_peers, peers)
		for len(pending_peers) != 0 {
			select {
			case a := <-score_responseRx:
				fmt.Printf("Received score for new order of %d, from %d \n", a[1], a[0])
				if a[1] < lowest_cost {
					lowest_cost = a[1]
					has_lowestcost = a[0]
				}
				for i := 0; i < len(pending_peers); i++ {
					if pending_peers[i] == a[0] {
						pending_peers = append(pending_peers[:i], pending_peers[i+1:]...)
					}
				}
			case <-time.After(time_to_respond):
				return
			}
		}
		order_handling.Insert(floor, button_type, has_lowestcost)
		update_order(id, button_type, floor, has_lowestcost, updateTx, score_responseRx)
		return
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
			if a == id {
				fmt.Printf("Received ack of orders from %d\n", a)
				break
			}
		case <-time.After(time_to_respond):
			break
		}
	}
}


func update_lights(){
	order_matrix := order_handling.Get_order_matrix()

	for floor_i := 0; floor_i < config.NUMFLOORS i++{
		for button_type_i := 0; button_type_i < config.NUMBUTTON_TYPES; button_type_i++ {
			if order_matrix[i][j] == order_handling.NO_ORDER {
				lights.Clear(button_type_i, floor_i)
			}else{
				lights.Set(button_type_i, floor_i)

			}
		}
	}
}