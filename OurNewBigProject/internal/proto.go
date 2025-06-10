package internal

import "OurNewBigProject/internal/MessageStruct"

type Proto struct {
	Name  string
	Peers *MessageStruct.PeerRepository
	Port  string
}

func NewProto(name string, port string) *Proto {
	return &Proto{
		Name:  name,
		Peers: MessageStruct.NewPeerRepository(),
		Port:  port,
	}
}
