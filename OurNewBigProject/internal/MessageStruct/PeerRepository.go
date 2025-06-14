package MessageStruct

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/url"
	"sort"
	"sync"
	"time"
)

const (
	peerValidationTimeOut = 1 * time.Second
)

type PeerRepository struct {
	rwMutex *sync.RWMutex
	peers   map[string]*Peer
}

func NewPeerRepository() *PeerRepository {
	peerRepository := &PeerRepository{
		rwMutex: &sync.RWMutex{},
		peers:   make(map[string]*Peer),
	}

	peerRepository.peersValidator()

	return peerRepository
}

func (p *PeerRepository) Add(peer *Peer) {
	p.rwMutex.Lock()
	defer p.rwMutex.Unlock()

	_, found := p.peers[peer.PubKeyStr]
	if !found {
		p.peers[peer.PubKeyStr] = peer
	}
}

func (p *PeerRepository) Delete(pubKey string) {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

	delete(p.peers, pubKey)
}

func (p *PeerRepository) Get(pubKey string) (*Peer, bool) {
	p.rwMutex.RLock()
	defer p.rwMutex.RUnlock()

	peer, found := p.peers[pubKey]
	return peer, found
}

func (p *PeerRepository) GetPeers() []*Peer {
	peersSlice := make([]*Peer, 0, len(p.peers))

	for _, peer := range p.peers {
		peersSlice = append(peersSlice, peer)
	}

	sort.Slice(peersSlice, func(i, j int) bool {
		return peersSlice[i].Name < peersSlice[j].Name
	})

	return peersSlice
}

func (p *PeerRepository) peersValidator() {
	ticker := time.NewTicker(peerValidationTimeOut)

	go func() {
		for {
			<-ticker.C
			for _, peer := range p.peers {
				u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%s", peer.AddrIP, peer.Port), Path: "/word"}

				c, _, _ := websocket.DefaultDialer.Dial(u.String(), nil)
				if c == nil {
					p.Delete(peer.PubKeyStr)
					continue
				}
				c.Close()
			}
		}
	}()
}
