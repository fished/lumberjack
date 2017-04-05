package message

import (
	"encoding/json"
	"time"
)

const (
	Text = iota
	JSON
	XML
)

type Message struct {
	Timestamp time.Time `json:timestamp`
	Data      []byte    `json:data`
	Format    int       `json:format`
}

func NewMessage(data []byte, format int) *Message {
	msg := Message{
		Timestamp: time.Now(),
		Data:      data,
		Format:    format,
	}
	return &msg
}

func (m *Message) JSON() ([]byte, error) {
	return json.Marshal(m)
}
