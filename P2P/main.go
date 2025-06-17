package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	ui()
}

func Writer(name string, text string) {
	conn, err := net.Dial("tcp", "192.168.1.156:4545")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	message := fmt.Sprintf("%s:%s\n", name, text)
	_, err = conn.Write([]byte(message))
	if err != nil {
		fmt.Printf("Write error: %v\n", err)
		return
	}
	fmt.Println("Message sent successfully")
}

func Listener() (string, string) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		fmt.Println("Ошибка при запуске сервера:", err)
		return "", ""
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Ошибка при принятии соединения:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) (string, string) {
	defer conn.Close()

	message, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка при чтении:", err)
		return "", ""
	}
	parts := strings.SplitN(string(message), ":", 2)
	fmt.Println(parts)
	return parts[0], parts[1]
}
