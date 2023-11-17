package p2p

import (
	"bytes"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

type Peer struct {
	conn net.Conn
}

func (p *Peer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

func (p *Peer) readLoop(msgChan chan *Message) {
	buf := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			break
		}
		msgChan <- &Message{
			Payload: bytes.NewReader(buf[:n]),
			From:    p.conn.RemoteAddr(),
		}
	}
	p.conn.Close()
}

type TcpTransport struct {
	listenAddr string
	listener   net.Listener
	AddPeer    chan *Peer
	DelPeer    chan *Peer
}

func NewTcpTransport(addr string) *TcpTransport {
	return &TcpTransport{
		listenAddr: addr,
		// addPeer: addPeer,
		// delPeer: delPeer,
	}
}

func (t *TcpTransport) ListenAndAccept() error {
	ls, err := net.Listen("tcp", t.listenAddr)
	if err != nil {
		return err
	}
	t.listener = ls
	// t.listener.Close()
	for {
		conn, err := ls.Accept()
		if err != nil {
			logrus.Error(err)
			continue
		}
		peer := Peer{
			conn: conn,
		}
		
		t.AddPeer <- &peer
	}
	return fmt.Errorf("TCP transport stopped reason: ?")
}
