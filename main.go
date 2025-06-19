package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type pair struct {
	Name string
	Text string
}

func main() {
	ui()
}

func Writer(name string, text string) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	if name == "KICK" {
		fmt.Println(fmt.Sprintf("%s %s", name, text))
		message := fmt.Sprintf("%s:%s", name, text)
		_, err = conn.Write([]byte(message))
		return
	}
	message := fmt.Sprintf("%s:%s", name, text)
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("Write error: %v\n", err)
		return
	}
	fmt.Println("Message sent successfully")
}

func Writing() ([]pair, error) {
	conn, err := net.Dial("tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Ошибка подключения: %v\n", err)
		return nil, nil
	}
	defer conn.Close()
	fmt.Println("Подключено к серверу")
	_, err = conn.Write([]byte("Read\n"))
	if err != nil {
		fmt.Printf("%v\n", err)
	}
	messages := make([]pair, 0, 100000000)
	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(string(response))
	parts := strings.SplitN(string(response), "|", 4096)
	for i := 0; i < len(parts); i++ {
		if parts[i] == "\n" {
			break
		}
		part := strings.SplitN(string(parts[i]), ":", 2)
		if parts[0] == "KICK" {
			var p []pair
			p = append(p, pair{Name: part[1], Text: "KICK"})
			return p, nil
		}
		messages = append(messages, pair{Name: part[0], Text: part[1]})
	}
	fmt.Println(messages)
	return messages, nil
}
