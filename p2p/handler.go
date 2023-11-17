package p2p

import (
	"fmt"
	"io"
)

type Handler interface {
	HandleMessage(*Message) error
}

type DefaultHandler struct{
	Version string
}

func (h *DefaultHandler) HandleMessage(msg *Message) error {
	b, err := io.ReadAll(msg.Payload)
	if err != nil {
		return err
	}
	fmt.Printf("handler msg from %s: %s\n", msg.From.String(), string(b))
	return nil
}


// func NewDefaultHandler() *Handler {
// 	return &Handler{}
// }
