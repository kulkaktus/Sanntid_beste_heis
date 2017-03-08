package network

import (
	"./../bcast"
	"./../conn"
	//"./../localip"
	"./../peers"
	"./../order_handling"
	"fmt"
	"net"
	"sort"
	"time"
)

const network_port string = ":40012"

var broadcast_ip string

type Order struct {
	Floor       int
	Button_type int
}

type Message struct {
	Msg  string
	Id   string
	Iter int
	//Message_type int
	//Order
	//Status_update, At_elevator int8
}

const interval = 15 * time.Millisecond
const timeout = 50 * time.Millisecond

var PeerUpdateCh chan peers.PeerUpdate

//returns transmit, receive channels
func Init() (chan<- Message, <-chan Message) {

	// We make a channel for receiving updates on the id's of the peers that are
	//  alive on the network
	PeerUpdateCh = make(chan peers.PeerUpdate)
	// We can disable/enable the transmitter after it has been started.
	// This could be used to signal that we are somehow "unavailable".
	peerTxEnable := make(chan bool)
	go peers.Transmitter(15647, "0", peerTxEnable)
	go peers.Receiver(15647, PeerUpdateCh)

	// We make channels for sending and receiving our custom data types
	helloTx := make(chan Message)
	helloRx := make(chan Message)
	// ... and start the transmitter/receiver pair on some port
	// These functions can take any number of channels! It is also possible to
	//  start multiple transmitters/receivers on the same port.
	go bcast.Transmitter(16569, helloTx)
	go bcast.Receiver(16569, helloRx)

	return helloTx, helloRx
	// The example message. We just send one of these every second.

}

func Transmitter(port int, id string, transmitEnable <-chan bool) {

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	enable := true
	for {
		select {
		case enable = <-transmitEnable:
		case <-time.After(interval):
		}
		if enable {
			conn.WriteTo([]byte(id), addr)
		}
	}
}

func Receiver(port int, peerUpdateChannel chan<- peers.PeerUpdate) {

	var buf [1024]byte
	var p peers.PeerUpdate
	lastSeen := make(map[string]time.Time)

	conn := conn.DialBroadcastUDP(port)

	for {
		updated := false

		conn.SetReadDeadline(time.Now().Add(interval))
		n, _, _ := conn.ReadFrom(buf[0:])

		id := string(buf[:n])

		// Adding new connection
		p.New = ""
		if id != "" {
			if _, idExists := lastSeen[id]; !idExists {
				p.New = id
				updated = true
			}

			lastSeen[id] = time.Now()
		}

		// Removing dead connection
		p.Lost = make([]string, 0)
		for k, v := range lastSeen {
			if time.Now().Sub(v) > timeout {
				updated = true
				p.Lost = append(p.Lost, k)
				delete(lastSeen, k)
			}
		}

		// Sending update
		if updated {
			p.Peers = make([]string, 0, len(lastSeen))

			for k, _ := range lastSeen {
				p.Peers = append(p.Peers, k)
			}

			sort.Strings(p.Peers)
			sort.Strings(p.Lost)
			peerUpdateChannel <- p
		}
	}
}

func SendMessage(message Message, recipient_id string) {

}

func Send_order(cost int, order order_handling.Order, tx chan, {
	
}


/*



func receive() (m Message, address int, error){
	_, addr, err := conn.ReadFromUDP(buffer)
	return m, address, err
}

func send(m Message, elevator Elevator)(returned_value int, error) {
	udp_addr, err := net.ResolveUDPAddr("udp", elevator.address+network_port)
	conn, err := net.DialUDP("udp", nil, udp_addr)
	if err != nil {
		fmt.Printf("Some error %v", err)
		return
	}
	conn.Write([]byte("This is Patrick"))

}

func broadcast(m Message) {

}


*/
