package internal

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

func HandleConnection(conn net.Conn, messages []Message) []Message {
	conn.SetWriteDeadline(time.Time{})
	messag, _ := bufio.NewReader(conn).ReadString('\n')
	fmt.Printf("Получено: %s", messag)
	if messag == "Read\n" {
		fmt.Println("Start of sending")
		//conn.Write([]byte(fmt.Sprintf("%v\n", len(messages))))
		for i := 0; i < len(messages); i++ {
			//fmt.Print(messages[i].Name, " ", messages[i].Text)
			_, err := conn.Write([]byte(fmt.Sprintf("%v:%v|", messages[i].Name, messages[i].Text)))
			fmt.Println(fmt.Sprintf("%v:%v", messages[i].Name, messages[i].Text))
			if err != nil {
				fmt.Println(err)
			}
		}
		conn.Write([]byte("&"))
		for u := range AllUsersNames {
			_, err := conn.Write([]byte(fmt.Sprintf("%v|", u)))
			fmt.Println(fmt.Sprintf("%v", u))
			if err != nil {
				fmt.Println(err)
			}
		}
		conn.Write([]byte("\n"))
		fmt.Println("End of sending")
	} else {
		parts := strings.SplitN(string(messag), ":", 2)
		fmt.Println(parts)
		if parts[1] == "EXIT" {
			delete(AllUsersNames, parts[0])
			return messages
		}
		messages = append(messages, Message{parts[0], parts[1]})
		AllUsersNames[parts[0]] = struct{}{}
	}
	for i := 0; i < len(messages); i++ {
		fmt.Println(messages[i].Name, " ", messages[i].Text)
	}
	for u := range AllUsersNames {
		fmt.Println(u)
	}
	conn.Close()
	return messages
}
