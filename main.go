package main

import (
	"./config"
	"./fsm"
	"./io/buttons"
	"./io/io"
	"./io/lights"
	"./io/motor"
	"./io/sensors"
	"./network"
	"./network/peers"
	"./order_handling"
	"fmt"
	"os"
	"strconv"
	"time"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.

func main() {
	lights.Init()
	peers.Init()
	config.Init()
	buttons.Init()
	motor.Init()
	sensors.Init()
	io.Init()
	var id_in string
	var id int
	var err error
	err = nil
	time.Sleep(1 * time.Millisecond)
	for {
		if len(os.Args) <= 1 || err != nil {
			fmt.Println("Enter id as a unique positive integer: ")
			id_in = ""
			for id_in == "" {
				fmt.Scanln(&id_in)
			}
		} else {
			id_in = os.Args[1]
		}
		id, err = strconv.Atoi(id_in)
		if err == nil {
			if id > 0{
				break
			}
		}

	}
	fmt.Print("My id is: ", id, "\n")
	order_handling.Init(id)
	ordersTx, ordersRx, updateTx, updateRx, messageTx, messageRx := network.Init(id) //get transmit and receive channels
	// The example message. We just send one of these every second.
	/*go func() {
		helloMsg := network.Message{"This is PATRICK!", id, 0}
		for {
			helloMsg.Iter++
			tx <- helloMsg
			time.Sleep(1 * time.Second)
		}
	}()*/
	fsm.Fsm(id, ordersTx, ordersRx, updateTx, updateRx, messageTx, messageRx)
}
