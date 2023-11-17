package main

import (
	"fmt"
	"time"

	"github.com/dattrantien00/p2p-pocker/p2p"
)

func main() {
	server := p2p.NewServer(
		p2p.ServerConfig{
			ListenAddr:  ":3000",
			Version:     "Version1",
			GameVarient: p2p.TexasHoldem,
		})
	go server.Start()
	time.Sleep(1 * time.Second)

	remoteServer := p2p.NewServer(
		p2p.ServerConfig{
			ListenAddr: ":4000",
			Version:    "Version1",
			GameVarient: p2p.TexasHoldem,
		})

	go remoteServer.Start()
	if err := remoteServer.Connect(":3000"); err != nil {
		fmt.Println(err)
	}
	// time.Sleep(1*time.Second)
	// server.Connect("127.0.0.1:8081")
	select {}
}
