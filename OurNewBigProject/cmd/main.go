package main

import (
	"OurNewBigProject/internal"
	"OurNewBigProject/internal/ui"
	"fmt"

	"flag"
)

const (
	Port = "25042"
)

func main() {
	name := flag.String("name", "krulsaidme0w", "peer's name")
	flag.Parse()

	p := internal.NewProto(*name, Port)

	runNetworkManager(p)

	runUI(p)
	var a string
	fmt.Scan(&a)
}

func runNetworkManager(p *internal.Proto) {
	networkManager := internal.NewManager(p)
	networkManager.Start()
}

func runUI(p *internal.Proto) error {
	return ui.NewApp(p).Run()
}
