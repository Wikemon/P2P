package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {
	ui()
}

func Writer() {
	conn, err := net.Dial("tcp", "192.168.1.156:4545")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	var s string
	fmt.Scan(&s)
	conn.Write([]byte(s))
	io.Copy(os.Stdout, conn)
	fmt.Println("\nDone")
}

func Listener() {
	listener, err := net.Listen("tcp", ":4545")

	if err != nil {
		fmt.Println(err)
		return
	}

	defer listener.Close()
	fmt.Println("Server is listening...")
	conn, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}
	input := make([]byte, (1024 * 4))
	n, err := conn.Read(input)
	if n == 0 || err != nil {
		fmt.Println("Read error:", err)
	}
	fmt.Println(string(input[:n]))
}
