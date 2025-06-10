package MessageStruct

import (
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"math/big"
	"net/url"
	"time"
)

var (
	ErrPeerIsDeleted = errors.New("peer disconnected")
)

type Peer struct {
	Name      string
	PubKey    *big.Int
	PubKeyStr string
	Port      string
	Messages  []*Message
	AddrIP    string
}

func (p *Peer) AddMessage(text, author string) {
	p.Messages = append(p.Messages, &Message{
		Time:   time.Now(),
		Text:   text,
		Author: author,
	})
}

func (p *Peer) SendMessage(pubKey, message string) error {
	u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%s", p.AddrIP, p.Port), Path: "/chat"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if c == nil {
		return ErrPeerIsDeleted
	}

	defer c.Close()
	if err != nil {
		return err
	}
	return c.WriteMessage(1, []byte(fmt.Sprintf("%s:%s", pubKey)))
}
