package main

import (
	"time"

	"github.com/dattrantien00/p2p-pocker/p2p"
)

func makeServerAndStart(addr string) *p2p.Server {
	server := p2p.NewServer(
		p2p.ServerConfig{
			ListenAddr:  addr,
			Version:     "Version1",
			GameVarient: p2p.TexasHoldem,
		})
	go server.Start()
	return server
}

func main() {
	playerA := makeServerAndStart(":3000")
	playerB := makeServerAndStart(":4000")
	playerC := makeServerAndStart(":5000")
	playerD := makeServerAndStart(":6000")
	time.Sleep(1 * time.Second)
	playerB.Connect(playerA.ListenAddr)
	playerC.Connect(playerB.ListenAddr)
	playerD.Connect(playerC.ListenAddr)
	// playerB.Connect(playerA.ListenAddr)
	_ = playerB

	select {}
}
