package internal

type Message struct {
	Name string
	Text string
}

var Messages []Message

var AllUsersNames map[string]struct{}
