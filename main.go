package main

import (
	"fmt"

	"./network"
	"./peers"
	"time"
)

// We define some custom struct to send over the network.
// Note that all members we want to transmit must be public. Any private members
//  will be received as zero-values.

func main() {

	tx, rx, id := network.Init() //get transmit and receive channels
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
