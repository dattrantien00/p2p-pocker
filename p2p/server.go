package p2p

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"

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

type Message struct {
	Payload io.Reader
	From    net.Addr
}

type Server struct {
	ServerConfig

	transport *TcpTransport
	handler   Handler
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

		handler:   &DefaultHandler{},
		peers:     make(map[net.Addr]*Peer),
		addPeer:   make(chan *Peer),
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
	}).Infoln("started new game server")
	s.transport.ListenAndAccept()
}

func (s *Server) Connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println(err)
		return err
	}
	peer := Peer{
		conn: conn,
	}
	s.addPeer <- &peer
	return peer.Send([]byte(s.Version))

}

func (s *Server) SendHandShake(p *Peer) error {
	hs := &HandShake{
		Version:     s.Version,
		GameVarient: s.GameVarient,
	}
	buf := new(bytes.Buffer)

	// if err := gob.NewEncoder(buf).Encode(hs); err != nil {
	// 	fmt.Println(err)
	// 	return err
	// }
	if err := hs.Encode(p.conn); err != nil {
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

			s.SendHandShake(peer)

			if err := s.handShake(peer); err != nil {
				logrus.Info("handshake with incoming player failed", err)
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

	return nil
}

type HandShake struct {
	Version     string
	GameVarient GameVarient
}

func (hs *HandShake) Encode(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, []byte(hs.Version)); err != nil {
		logrus.Errorln(err)
		return err
	}
	if err := binary.Write(w, binary.LittleEndian, []byte(hs.GameVarient.String())); err != nil {
		logrus.Errorln(err)
		return err
	}
	return nil
}

func (hs *HandShake) Decode(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, []byte(hs.Version)); err != nil {
		logrus.Errorln(err)
		return err
	}
	if err := binary.Read(r, binary.LittleEndian, []byte(hs.GameVarient.String())); err != nil {
		logrus.Errorln(err)
		return err
	}
	return nil
}
func (s *Server) handShake(peer *Peer) error {

	hs := &HandShake{}

	// if err := gob.NewDecoder(peer.conn).Decode(hs); err != nil {

	// 	return err
	// }
	if err := hs.Decode(peer.conn); err != nil {
		return err
	}

	if s.GameVarient != hs.GameVarient {
		return fmt.Errorf("invalid gamevariant %s", hs.GameVarient)
	}
	if s.Version != hs.Version {
		
		return fmt.Errorf("invalid version %s ", hs.Version)
	}
	logrus.WithFields(logrus.Fields{
		"peer":        peer.conn.RemoteAddr(),
		"version":     hs.Version,
		"gameVariant": hs.GameVarient.String(),
	}).Infoln("receive handshake")
	return nil
}
