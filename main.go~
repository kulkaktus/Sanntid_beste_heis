package main

import (
	"./config"
	"./io/buttons"
	"./io/elevio"
	"./io/lights"
	"./io/sensors"
	"./network"
	"./peers"
	"fmt"
	"os"
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
	elevio.Io_init()
	//floor := 0
	var id string
	if len(os.Args) <= 1 {
		fmt.Println("Enter id: ")
		for id == "" {
			fmt.Scanln(&id)
		}
	} else {
		id = os.Args[1]
	}
	if id == "" {
		fmt.Printf("error, no id\n")
	} else {
		fmt.Print("My id is: ", id, "\n")
	}
	/*for {
		for j := 1; j <= config.DOWN; j++ {
			for i := 1; i <= 4; i++ {
				if buttons.Get(j, i) {
					lights.Set(j, i)
				} else {
					lights.Clear(j, i)
				}
			}
		}
		if floor != sensors.Get() && sensors.Get() != 0 {
			floor = sensors.Get()
			fmt.Printf("Arrived at %d \n", floor)
		}
	}*/
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

	fmt.Println("Started")
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
	}
}
