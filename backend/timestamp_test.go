package backend

import (
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
	parser := TimestampParserRegistry["Stamp"]
	now := time.Now()
	t.Run("Can correctly parse simple timestamp", func(t *testing.T) {
		result, err := parser.ParseAtTime("Apr  9 02:54:38", time.Now())
		if err != nil {
			t.Logf("Failed with error: '%s'", err)
			t.Fail()
		}

		if correct := time.Date(now.Year(), time.April, 9, 2, 54, 38, 0, time.Local); correct != result {
			t.Logf("Times didn't match: %s <-> %s", result, correct)
		}

	},
	)
}
