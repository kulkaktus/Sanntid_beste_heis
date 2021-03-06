package network

import (
	"../config"
	"./bcast"
	"./peers"
	"time"
)

const network_port string = ":40012"
const (
	SCORE_RESPONSE_T  = 123
	ORDERS_RESPONSE_T = 321
	STUCK_SEND_T      = 456
)

var broadcast_ip string

type Orders struct {
	Orders  [config.NUMFLOORS][config.NUMBUTTON_TYPES]int
	From_id int
}

type Update struct {
	Floor       int
	Button_type int
	Executer    int
	From_id     int
}

type Message struct {
	To_id   int
	From_id int
	Type    int
	Content int
}

const interval = 15 * time.Millisecond
const timeout = 50 * time.Millisecond

var PeerUpdateCh chan peers.PeerUpdate

func Init(id int) (chan<- Orders, <-chan Orders, chan<- Update, <-chan Update, chan<- Message, <-chan Message) {
	PeerUpdateCh = make(chan peers.PeerUpdate)
	peerTxEnable := make(chan bool)

	go peers.Transmitter(47412, id, peerTxEnable)
	go peers.Receiver(47412, id, PeerUpdateCh)

	ordersTx := make(chan Orders)
	ordersRx := make(chan Orders)

	updateTx := make(chan Update)
	updateRx := make(chan Update)

	messageTx := make(chan Message)
	messageRx := make(chan Message)

	go bcast.Transmitter(47512, ordersTx)
	go bcast.Receiver(47512, ordersRx)
	go bcast.Transmitter(47612, updateTx)
	go bcast.Receiver(47612, updateRx)
	go bcast.Transmitter(47712, messageTx)
	go bcast.Receiver(47712, messageRx)

	return ordersTx, ordersRx, updateTx, updateRx, messageTx, messageRx
}
