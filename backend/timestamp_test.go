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
	})

	t.Run("Can parse simple timestamp when the year rolls over", func(t *testing.T) {
		now := time.Date(2017, time.January, 1, 0, 0, 1, 0, time.Local) // January 1, 12:01AM
		result, err := parser.ParseAtTime("Dec 31 23:59:59", now)
		if err != nil {
			t.Logf("Failed with error: '%s'", err)
			t.Fail()
		}
		if correct := time.Date(2016, time.December, 31, 23, 59, 59, 0, time.Local); correct != result {
			t.Logf("Times didn't match: %s <-> %s", result, correct)
		}
	})

	// TODO: More tests with other time formats.
}
