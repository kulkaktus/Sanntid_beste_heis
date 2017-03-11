package main

import (
	"./config"
	"./io/buttons"
	"./io/io"
	"./io/lights"
	//"./io/motor"
	"./fsm"
	"./io/sensors"
	"./network"
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
	lights.Lights_init()
	peers.Peers_init()
	config.Config_init()
	buttons.Buttons_init()
	sensors.Sensors_init()
	io.Io_init()
	var id_in string
	var id int
	var err error
	err = nil
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

	tx, rx := network.Init() //get transmit and receive channels
	// The example message. We just send one of these every second.
	go func() {
		helloMsg := network.Message{"Hello!", id, 0}
		for {
			helloMsg.Iter++
			tx <- helloMsg
			time.Sleep(1 * time.Second)
		}
	}()
	fsm.Fsm(id, tx, rx)
}
