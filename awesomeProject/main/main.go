package main

import (
	"awesomeProject/internal"
	"fmt"
	"net"
)

func main() {
	internal.AllUsersNames = make(map[string]struct{})
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer ln.Close()
	fmt.Println("Сервер запущен на :8080")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		internal.Messages = internal.HandleConnection(conn, internal.Messages)
	}
}
