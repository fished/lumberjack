package backend

import (
	"os"
	"regexp"
	"testing"
)

const BASE = `/tmp/lumberjack_test/`

func TestRecorder(t *testing.T) {
	var err error
	re := regexp.MustCompile(`^(?P<timestamp>.{15})\s+((?P<instance>\S+)\s+(?P<source>.+?)\[(?P<pid>.*?)\]:)?`)

	msgParser := NewStringMessageParser(re, TimestampParserRegistry["Stamp"], "none", "none")
	msg, err := msgParser("Apr  7 12:03:17 avila Docker[32336]: SC database lists search domains:")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	err = os.MkdirAll(BASE, 0700)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(BASE)

	Recorder, err := NewRecorder(BASE, re.SubexpNames())
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Test a standard line", func(t *testing.T) {
		Recorder(msg)
	})
}
