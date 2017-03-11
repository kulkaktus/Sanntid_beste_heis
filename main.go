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
	"./order_handling"
	"./peers"
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
	order_handling.Init()
	io.Init()
	var id_in string
	var id int
	var err error
	err = nil
	time.Sleep(1 * time.Millisecond)
	for {
		if len(os.Args) <= 1 || err != nil {
			fmt.Println("Enter id: ")
			id_in = ""
			for id_in == "" {
				fmt.Scanln(&id_in)
			}
		} else {
			id_in = os.Args[1]
		}
		id, err = strconv.Atoi(id_in)
		if err == nil {
			break
		}

	}
	fmt.Print("My id is: ", id, "\n")

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
