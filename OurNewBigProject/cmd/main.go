package main

import (
	"OurNewBigProject/internal"
	"OurNewBigProject/internal/ui"
	"flag"
	"log"
)

const (
	Port = "25042"
)

func main() {
	name := flag.String("name", "krulsaidme0w", "peer's name")
	flag.Parse()

	p := internal.NewProto(*name, Port)

	runNetworkManager(p)
	
	if err := runUI(p); err != nil {
		log.Fatal(err)
	}
}

func runNetworkManager(p *internal.Proto) {
	networkManager := internal.NewManager(p)
	networkManager.Start()
}

func runUI(p *internal.Proto) error {
	return ui.NewApp(p).Run()
}
