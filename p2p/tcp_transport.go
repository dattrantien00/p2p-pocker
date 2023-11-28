package p2p

import (
	"encoding/gob"
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
)

type Peer struct {
	conn     net.Conn
	outbound bool //node1 ->node2 => node1 = true, node2= false
	listenAddr string
}

func (p *Peer) Send(b []byte) error {
	_, err := p.conn.Write(b)
	return err
}

func (p *Peer) readLoop(msgChan chan *Message) {
	for {
		msg := new(Message)
		if err := gob.NewDecoder(p.conn).Decode(msg); err != nil {
			logrus.Errorf("decode message error: %s", err)
			break
		}
		fmt.Printf("%+v, %s\n", msg, "receive")
		msgChan <- msg
	}
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
			conn:     conn,
			outbound: false,
		}

		t.AddPeer <- &peer
	}

}
