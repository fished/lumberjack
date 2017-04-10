package backend

import "time"

type Message struct {
	Timestamp time.Time         `json:timestamp`
	Data      string            `json:data`
	Source    string            `json:source`
	Host      string            `json:host`
	Keys      map[string]string `json:keys`
}
