package internal

import (
	"OurNewBigProject/internal/MessageStruct"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/big"
	"net"
	"net/http"
	"strings"
	"time"
)

var (
	upgrader = websocket.Upgrader{}
)

const (
	MulticastIP             = "224.0.0.1"
	ListenerIP              = "0.0.0.0"
	MulticastFrequency      = 1 * time.Second
	udpConnectionBufferSize = 1024
)

type Discoverer struct {
	Addr               *net.UDPAddr
	MulticastFrequency time.Duration
	Proto              *Proto
}

func NewDiscoverer(addr *net.UDPAddr, multicastFrequency time.Duration, proto *Proto) *Discoverer {
	return &Discoverer{
		Addr:               addr,
		MulticastFrequency: multicastFrequency,
		Proto:              proto,
	}
}

func (d *Discoverer) Start() {
	go d.startMulticasting()
	go d.listenMulticasting()
}

func (d *Discoverer) startMulticasting() {
	conn, err := net.DialUDP("udp", nil, d.Addr)
	if err != nil {
		log.Fatal(err)
	}

	ticker := time.NewTicker(d.MulticastFrequency)
	for {
		<-ticker.C
		_, err := conn.Write([]byte(fmt.Sprintf("%s:%s:%s:%s",
			"qwer",
			d.Proto.Name,
			d.Proto.Port)))
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (d *Discoverer) listenMulticasting() {
	conn, err := net.ListenMulticastUDP("udp", nil, d.Addr)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.SetReadBuffer(udpConnectionBufferSize)
	if err != nil {
		log.Fatal(err)
	}

	for {
		rawBytes, addr, err := ReadFromUDPConnection(conn, udpConnectionBufferSize)
		if err != nil {
			log.Fatal(err)
		}

		message, err := MessageStruct.UDPMulticastMessageToPeer(rawBytes)
		if err != nil {
			log.Fatal(err)
		}

		peer := &MessageStruct.Peer{
			Name:      message.Name,
			PubKey:    message.PubKey,
			PubKeyStr: message.PubKeyStr,
			Port:      message.Port,
			Messages:  make([]*MessageStruct.Message, 0),
			AddrIP:    addr.IP.String(),
		}
		d.Proto.Peers.Add(peer)
	}
}

type Listener struct {
	proto *Proto
	addr  string
}

func NewListener(addr string, proto *Proto) *Listener {
	return &Listener{
		proto: proto,
		addr:  addr,
	}
}

func (l *Listener) chat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		arr := strings.Split(string(message), ":")
		if len(arr) != 2 {
			continue
		}

		pubKeyStr := arr[0]
		messageText := arr[1]

		peer, found := l.proto.Peers.Get(pubKeyStr)
		if !found {
			continue
		}

		pubKey := new(big.Int)
		pubKey, ok := pubKey.SetString(pubKeyStr, 10)
		if !ok {
			continue
		}

		peer.AddMessage(messageText, peer.Name)
	}
}

func (l *Listener) meow(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	conn.Close()
}

func (l *Listener) Start() {
	http.HandleFunc("/chat", l.chat)
	http.HandleFunc("/word", l.meow)
	log.Fatal(http.ListenAndServe(l.addr, nil))
}

type Manager struct {
	Proto      *Proto
	Listener   *Listener
	Discoverer *Discoverer
}

func NewManager(proto *Proto) *Manager {
	multicastAddr, err := net.ResolveUDPAddr(
		"udp",
		fmt.Sprintf("%s:%s", MulticastIP, proto.Port))
	if err != nil {
		log.Fatal(err)
	}

	listenerAddr := fmt.Sprintf("%s:%s", ListenerIP, proto.Port)

	return &Manager{
		Proto:      proto,
		Listener:   NewListener(listenerAddr, proto),
		Discoverer: NewDiscoverer(multicastAddr, MulticastFrequency, proto),
	}
}

func (m *Manager) Start() {
	go m.Listener.Start()
	go m.Discoverer.Start()
}
