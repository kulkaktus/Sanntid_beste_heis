package peers

import (
	"../conn"
	"fmt"
	"net"
	"sort"
	"strconv"
	"time"
)

type PeerUpdate struct {
	Peers []int
	New   int
	Lost  []int
}

const interval = 15 * time.Millisecond
const timeout = 50 * time.Millisecond

func Init() error {
	return nil
}

func Transmitter(port int, id int, transmitEnable <-chan bool) {

	conn := conn.DialBroadcastUDP(port)
	addr, _ := net.ResolveUDPAddr("udp4", fmt.Sprintf("255.255.255.255:%d", port))

	enable := true
	for {
		select {
		case enable = <-transmitEnable:
		case <-time.After(interval):
		}
		if enable {
			conn.WriteTo([]byte(strconv.Itoa(id)), addr)
		}
	}
}

func Receiver(port int, id_self int, peerUpdateCh chan<- PeerUpdate) {

	var buf [1024]byte
	var p PeerUpdate
	lastSeen := make(map[int]time.Time)

	conn := conn.DialBroadcastUDP(port)

	for {
		updated := false

		conn.SetReadDeadline(time.Now().Add(interval))
		n, _, _ := conn.ReadFrom(buf[0:])
		id, err := strconv.Atoi(string(buf[:n]))
		if err != nil {
			id = 0
		}
		/*if id_ == "a" {
			fmt.Println(n)
		}*/
		// Adding new connection
		p.New = 0
		if id != 0 && id != id_self {
			if _, idExists := lastSeen[id]; !idExists {
				p.New = id
				updated = true
			}

			lastSeen[id] = time.Now()
		}
		// Removing dead connection
		p.Lost = make([]int, 0)
		for k, v := range lastSeen {
			if time.Now().Sub(v) > timeout {
				updated = true
				p.Lost = append(p.Lost, k)
				delete(lastSeen, k)
			}
		}

		// Sending update
		if updated {
			p.Peers = make([]int, 0, len(lastSeen))

			for k := range lastSeen {
				p.Peers = append(p.Peers, k)
			}

			sort.Ints(p.Peers)
			sort.Ints(p.Lost)
			peerUpdateCh <- p
		}
	}
}
