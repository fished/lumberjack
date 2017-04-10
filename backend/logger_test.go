package backend

import (
	"os"
	"regexp"
	"testing"
)

const BASE = `/tmp/lumberjack_files/`
const TIME_LAYOUT = "Jan  2 15:04:05"

func TestLogger(t *testing.T) {
	re := regexp.MustCompile(`^(?P<timestamp>.{15})\s+((?P<host>\S+)\s+(?P<process>.+?)\[(?P<pid>.*?)\]:)?`)

	err := os.MkdirAll(BASE, 0700)
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(BASE)

	logger, err := newLogger(BASE, re, TimestampParserRegistry["Stamp"])
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Test a standard line", func(t *testing.T) {
		logger("Apr  7 12:03:17 avila Docker[32336]: SC database lists search domains:")
	})
}
