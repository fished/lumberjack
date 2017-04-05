package message

import (
	"testing"
	"time"
)

func TestMessage(t *testing.T) {
	t.Run("can create a message", func(t *testing.T) {
		m := NewMessage([]byte("This is a test message"), Text)
		if m == nil {
			t.Fail()
		}
	})

	t.Run("new message has a current timestamp", func(t *testing.T) {
		basetime := time.Now()
		m := NewMessage([]byte("This is a test message"), Text)
		if m.Timestamp.Before(basetime) {
			t.Fail()
		}
	})
}
