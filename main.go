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
	// playerE := makeServerAndStart(":7000")
	// playerF := makeServerAndStart(":8000")


	time.Sleep(1000 * time.Millisecond)
	playerB.Connect(playerA.ListenAddr)
	time.Sleep(10000 * time.Millisecond)
	playerC.Connect(playerB.ListenAddr)
	time.Sleep(1000 * time.Millisecond)
	playerD.Connect(playerC.ListenAddr)
	time.Sleep(1000 * time.Millisecond)
	// playerE.Connect(playerD.ListenAddr)
	// time.Sleep(1000 * time.Millisecond)
	// playerF.Connect(playerE.ListenAddr)
	// playerB.Connect(playerA.ListenAddr)
	// _ = playerB

	select {}
}
