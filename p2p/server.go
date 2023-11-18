package p2p

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

type GameVarient int

func (gv GameVarient) String() string {
	return []string{"Texas Holdem", "Other"}[gv]
}

const (
	TexasHoldem GameVarient = iota
	Other
)

type TCPTransport struct{}

type ServerConfig struct {
	ListenAddr  string
	Version     string
	GameVarient GameVarient
}

type Server struct {
	ServerConfig

	transport *TcpTransport

	// listener  net.Listener
	// mu        sync.RWMutex
	peers     map[net.Addr]*Peer
	addPeer   chan *Peer
	delPeer   chan *Peer
	msgCh     chan *Message
	gameState *GameState
}

func NewServer(cfg ServerConfig) *Server {
	s := &Server{
		ServerConfig: cfg,

		peers:     make(map[net.Addr]*Peer),
		addPeer:   make(chan *Peer,10),
		delPeer:   make(chan *Peer),
		msgCh:     make(chan *Message),
		gameState: NewGameState(),
	}
	tr := NewTcpTransport(cfg.ListenAddr)
	tr.AddPeer = s.addPeer
	tr.DelPeer = s.delPeer
	s.transport = tr
	return s
}

func (s *Server) Start() {

	go s.loop()

	logrus.WithFields(logrus.Fields{
		"port":        s.ListenAddr,
		"gameVariant": s.GameVarient,
		"gameStatus":  s.gameState.GameStatus,
	}).Infoln("started new game server")
	s.transport.ListenAndAccept()
}

func (s *Server) Connect(addr string) error {
	fmt.Printf("dialing from %s to %s\n",s.ListenAddr,addr)
	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		fmt.Println(err)
		return err
	}
	peer := Peer{
		conn:     conn,
		outbound: true,
	}
	
	s.addPeer <- &peer
	return s.SendHandShake(&peer)

}

func (s *Server) sendPeerList(p *Peer) error {

	peerList := MessagePeerList{
		Peers: make([]string, len(s.peers)),
	}
	it := 0
	for addr, _ := range s.peers {
		peerList.Peers[it] = addr.String()
		it++
	}
	msg := NewMessage(s.ListenAddr, peerList)
	buf := new(bytes.Buffer)
	if err := gob.NewEncoder(buf).Encode(msg); err != nil {
		return err
	}

	return p.Send(buf.Bytes())
}
func (s *Server) SendHandShake(p *Peer) error {
	hs := &HandShake{
		Version:     s.Version,
		GameVarient: s.GameVarient,
		GameStatus:  s.gameState.GameStatus,
	}
	buf := new(bytes.Buffer)

	if err := gob.NewEncoder(buf).Encode(hs); err != nil {
		fmt.Println(err)
		return err
	}

	if err := p.Send(buf.Bytes()); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
func (s *Server) loop() {
	for {
		select {
		case peer := <-s.addPeer:
			
			// s.SendHandShake(peer)
			if !peer.outbound {
				if err := s.SendHandShake(peer); err != nil {
					logrus.Errorln("fail to send handshake with peer ", err)
					peer.conn.Close()
					delete(s.peers, peer.conn.RemoteAddr())
					continue
				}

				if err := s.sendPeerList(peer); err != nil {
					logrus.Errorf("peerlist error: %s", err)
					continue
				}
			}
			if err := s.handShake(peer); err != nil {
				logrus.Errorln("handshake with incoming player failed", err)
				peer.conn.Close()
				delete(s.peers, peer.conn.RemoteAddr())
				continue
			}

			go peer.readLoop(s.msgCh)

			logrus.WithFields(logrus.Fields{
				"addr": peer.conn.RemoteAddr(),
			}).Infoln("handshake successfull: new player connected")

			s.peers[peer.conn.RemoteAddr()] = peer

		case delPeer := <-s.delPeer:
			logrus.WithFields(logrus.Fields{
				"addr": delPeer.conn.RemoteAddr(),
			}).Infoln("delete player connection")
			delete(s.peers, delPeer.conn.RemoteAddr())

		case msg := <-s.msgCh:
			if err := s.handleMessge(msg); err != nil {
				panic(err)
			}
		}
	}
}

func (s *Server) handleMessge(msg *Message) error {
	logrus.WithFields(logrus.Fields{
		"from": msg.From,
	}).Infoln("receive message")

	switch v := msg.Payload.(type) {
	case MessagePeerList:
		fmt.Printf("%+v\n", v)
		return s.handlePeerList(v)
	}
	return nil
}

func (s *Server) handShake(peer *Peer) error {

	hs := &HandShake{}

	if err := gob.NewDecoder(peer.conn).Decode(hs); err != nil {

		return err
	}

	if s.GameVarient != hs.GameVarient {
		return fmt.Errorf("gamevariant does not match %s", hs.GameVarient)
	}
	if s.Version != hs.Version {

		return fmt.Errorf(" version does not match %s ", hs.Version)
	}
	logrus.WithFields(logrus.Fields{
		"peer":        peer.conn.RemoteAddr(),
		"version":     hs.Version,
		"gameVariant": hs.GameVarient.String(),
		"gameStatus":  hs.GameStatus,
	}).Infoln("receive handshake")
	return nil
}

func (s *Server) handlePeerList(l MessagePeerList) error {
	
	for i := 0; i < len(l.Peers); i++ {
	
		if err := s.Connect(l.Peers[i]); err != nil {
			logrus.Errorf("fail to dail peer: %s", err)
			continue
		}
	}
	return nil
}
func init() {
	gob.Register(MessagePeerList{})
}
