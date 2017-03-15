package fsm

import (
	"../config"
	"../io/buttons"
	"../io/lights"
	"../io/motor"
	"../io/sensors"
	"../network"
	"../order_handling"
	"fmt"
	"time"
)

var floor int
var peers []int
var peers_internal_orders map[int][config.NUMFLOORS]int
var state string

const (
	tries_to_send              = 10
	time_to_respond            = 200 * time.Millisecond
	time_threshold_motor_stuck = 5 * time.Second
	doors_open_for             = 1 * time.Second
)

func Fsm(id int, ordersTx chan<- network.Orders, ordersRx <-chan network.Orders, updateTx chan<- network.Update, updateRx <-chan network.Update, messageTx chan<- network.Message, messageRx <-chan network.Message) {
	peers_internal_orders = make(map[int][config.NUMFLOORS]int)
	score_responseRx := make(chan [2]int)
	orders_responseRx := make(chan int)

	motor.Go(config.DOWN)
	for sensors.Get() == 0 {
	}

	motor.Stop()
	floor = sensors.Get()
	fmt.Printf("Started at floor %d\n", floor)
	order_handling.Set_floor(floor)
	go message_manager(id, ordersTx, ordersRx, updateTx, updateRx, messageTx, messageRx, score_responseRx, orders_responseRx, orders_responseRx)

	doors_open_since := time.Now().Add(-time.Hour)
	motor_on_since := time.Now().Add(-time.Hour)
	status_time_since := time.Now()

	wait_for_order_assignment := true
	state = "idle"

	for {
		order_handling.State = state
		check_buttons_and_update_orders(id, updateTx, score_responseRx)
		update_lights()

		if time.Second < time.Since(status_time_since) {
			status_time_since = time.Now()
			fmt.Println(state)
		}

		switch state {
		case "idle":
			current_order_destination := order_handling.Get_next()
			if current_order_destination == 0 {
				motor.Stop()
				state = "idle"
			} else if wait_for_order_assignment == true {
				time.Sleep(2 * time_to_respond)
				wait_for_order_assignment = false
			} else if current_order_destination < floor {
				motor.Go(config.DOWN)
				motor_on_since = time.Now()
				order_handling.Set_direction(config.DOWN)
				state = "running"

				if order_handling.Already_exists(current_order_destination, config.UP) {
					update_order(id, config.UP, current_order_destination, id, updateTx, score_responseRx)
				}
				if order_handling.Already_exists(current_order_destination, config.DOWN) {
					update_order(id, config.DOWN, current_order_destination, id, updateTx, score_responseRx)
				}
			} else if current_order_destination > floor {
				motor.Go(config.UP)
				motor_on_since = time.Now()
				order_handling.Set_direction(config.UP)
				state = "running"

				if order_handling.Already_exists(current_order_destination, config.UP) {
					update_order(id, config.UP, current_order_destination, id, updateTx, score_responseRx)
				}
				if order_handling.Already_exists(current_order_destination, config.DOWN) {
					update_order(id, config.DOWN, current_order_destination, id, updateTx, score_responseRx)
				}
			} else if current_order_destination == floor {
				motor.Stop()
				doors_open_since = time.Now()
				order_handling.Clear_orders_in_floor(floor)
				state = "door_open"

				if order_handling.Order_is_valid(floor, config.UP) {
					update_order(id, config.UP, floor, order_handling.NO_ORDER, updateTx, score_responseRx)
				}
				if order_handling.Order_is_valid(floor, config.DOWN) {
					update_order(id, config.DOWN, floor, order_handling.NO_ORDER, updateTx, score_responseRx)
				}
				wait_for_order_assignment = true
			}

		case "running":
			sensor := sensors.Get()
			if sensor != 0 && sensor != floor {
				floor = sensor
				motor.Stop()
				lights.Set(config.INDICATE, floor)
				order_handling.Set_floor(floor)
				state = "idle"

				fmt.Printf("Arrived at %d \n", floor)
			} else if sensor == floor && order_handling.Get_next() == 0 {
				motor.Stop()
			} else if time.Since(motor_on_since) > time_threshold_motor_stuck {
				send_stuck_message(id, messageTx)
				order_handling.Unassign_orders_handled_by(id)
				state = "stuck"
				fmt.Println("I am stuck")
			}

		case "door_open":
			if time.Since(doors_open_since) < doors_open_for {
				lights.Set(config.DOOR, 0)
			} else {
				lights.Clear(config.DOOR, 0)
				state = "idle"
			}

		case "stuck":
			sensor := sensors.Get()
			if (sensor != 0) && (sensor != floor) {
				floor = sensor
				motor.Stop()
				lights.Set(config.INDICATE, floor)
				order_handling.Set_floor(floor)
				state = "idle"

				fmt.Println("I am no longer stuck")
				fmt.Printf("Arrived at %d \n", floor)
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
								if update_order(id, button_type_i, floor_i, id, updateTx, score_responseRx) {
									order_handling.Insert(floor_i, button_type_i, id)
								}
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
	status_time_since := time.Now()
	for {
		if time.Second < time.Since(status_time_since) {
			status_time_since = time.Now()
			fmt.Println("network online")
		}
		select {
		case p := <-network.PeerUpdateCh:
			peers = p.Peers
			fmt.Printf("Peer update:\n")
			fmt.Printf("  Peers:    %d\n", p.Peers)
			fmt.Printf("  New:      %d\n", p.New)
			fmt.Printf("  Lost:     %d\n", p.Lost)

			if _, exists := peers_internal_orders[p.New]; !exists {
				peers_internal_orders[p.New] = [4]int{order_handling.NO_ORDER, order_handling.NO_ORDER, order_handling.NO_ORDER, order_handling.NO_ORDER}
			}
			go send_orders(id, p.New, order_handling.Get_order_matrix(), ordersTx, orders_responseRx_in, score_responseRx)
		case b := <-ordersRx:
			if b.From_id != id {
				order_handling.Print_external_order_matrix(b.Orders)
				order_handling.Merge_external_order_matrix_with_current(b.Orders)
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
					fmt.Println("sendt my cost")
					messageTx <- network.Message{d.From_id, id, network.SCORE_RESPONSE_T, order_handling.Get_cost(d.Floor, d.Button_type)}
				}
			}
			fmt.Printf(str)
		case f := <-messageRx:
			if f.To_id == id {
				fmt.Printf("Received: %d of type %d from %d\n", f.Content, f.Type, f.From_id)
				if f.Type == network.SCORE_RESPONSE_T {
					score_responseRx <- [2]int{f.From_id, f.Content}
				} else if f.Type == network.ORDERS_RESPONSE_T {
					orders_responseRx_out <- f.From_id
				} else if f.Type == network.STUCK_SEND_T {
					order_handling.Unassign_orders_handled_by(f.From_id)
				}
			}
		}
	}
	panic("no network")
}

func update_order(id int, button_type int, floor int, handler int, updateTx chan<- network.Update, score_responseRx <-chan [2]int) bool {
	for i := 0; i < tries_to_send; i++ {
		updateTx <- network.Update{floor, button_type, handler, id}
	}
	fmt.Println("sending order update")
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
			break
		}
	}
	if len(pending_peers) != 0 {

		for {
			_, not_empty := <-score_responseRx
			if !not_empty {
				break
			}
		}
		return true
	}
	return false
}

func new_order(id int, button_type int, floor int, updateTx chan<- network.Update, score_responseRx <-chan [2]int) bool {
	pending_peers := make([]int, len(peers))
	copy(pending_peers, peers)
	lowest_cost := order_handling.Get_cost(floor, button_type)
	has_lowestcost := id
	for i := 0; i < tries_to_send && len(pending_peers) != 0; i++ {
		updateTx <- network.Update{floor, button_type, order_handling.NO_EXECUTER, id}
	}
	fmt.Println("sending new order")
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
			return false
		}
	}
	for {
		_, not_empty := <-score_responseRx
		if !not_empty {
			break
		}
	}
	if len(pending_peers) == 0 {
		if update_order(id, button_type, floor, has_lowestcost, updateTx, score_responseRx) {
			order_handling.Insert(floor, button_type, has_lowestcost)
		}
		return true //Skal returnere true uavhengig av update_order sin returverdi
	} else {
		return false
	}
}
func send_orders(id int, to_id int, orders [config.NUMFLOORS][config.NUMBUTTON_TYPES]int, ordersTx chan<- network.Orders, orders_responseRx <-chan int, score_responseRx <-chan [2]int) {
	for i, v := range peers_internal_orders[to_id] {
		orders[i][0] = v
	}
	fmt.Println("sending all orders")
	for i := 0; i < send_attempts; i++ {
		ordersTx <- network.Orders{orders, id}
	}
	select {
	case a := <-orders_responseRx:
		if a == to_id {
			fmt.Printf("Received ack of orders from %d\n", a)
			break
		}
	case <-time.After(time_to_respond):
		break
	}
	for {
		_, not_empty := <-score_responseRx
		if !not_empty {
			break
		}
	}
}

func update_lights() {
	order_matrix := order_handling.Get_order_matrix()

	for floor_i := 0; floor_i < config.NUMFLOORS; floor_i++ {
		for button_type_i := 0; button_type_i < config.NUMBUTTON_TYPES; button_type_i++ {
			if order_matrix[floor_i][button_type_i] == order_handling.NO_ORDER {
				lights.Clear(button_type_i, floor_i+1)
			} else {
				lights.Set(button_type_i, floor_i+1)

			}
		}
	}
}

func send_stuck_message(id int, messageTx chan<- network.Message) {
	pending_peers := make([]int, len(peers))
	copy(pending_peers, peers)
	for i := 0; i < tries_to_send; i++ {
		for j := 0; j < len(pending_peers); j++ {
			messageTx <- network.Message{pending_peers[j], id, network.STUCK_SEND_T, 0}
		}
	}
}
