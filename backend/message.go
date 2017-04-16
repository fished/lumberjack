package backend

import (
	"regexp"

	"github.com/fished/lumberjack/util"
	"time"
)

type Message struct {
	Timestamp time.Time         `json:"timestamp"`
	Data      string            `json:"data"`
	Keys      map[string]string `json:"keys"`
}

func MessageKey(key string, value string) func(*Message) error {
	return func(msg *Message) error {
		if value != msg.Keys[key] {
			msg.Keys[key] = value
		}
		return nil
	}
}

type MessageParser struct {
	Parse func(data string, options ...func(*Message) error) (*Message, error)
	IndexedKeys []string
}

func NewStringMessageParser(
	parseRegexp *regexp.Regexp,
	tsParser TimestampParser,
	defaultSource, defaultInstance string) *MessageParser {

	mp := &MessageParser{
		Parse: func(data string, options ...func(*Message) error) (*Message, error) {
			re := parseRegexp.Copy() // Ensure that we avoid a locking situation in the regex library.

			msg := &Message{
				Timestamp: time.Now(),
				Data: data,
				Keys: map[string]string{
					"source":    defaultSource,
					"instance":  defaultInstance,
				},
			}

			/* Find indexable patters */
			matches := re.FindAllStringSubmatch(msg.Data, -1)
			indexable := util.DecodeNamedSubmatches(matches, re)

			// Any values picked up by regex in data will override defaults.
			for _, i := range indexable {
				if i == nil || i[0] == "" {
					continue
				}

				switch i[0] {
				case "timestamp":
					if timestamp, err := tsParser.Parse(i[1]); err != nil {
						return nil, err
					} else {
						msg.Timestamp = timestamp
					}
				default:
					msg.Keys[i[0]] = i[1]
				}
			}

			// Run any options on the parser.
			for _, opt := range options {
				if err := opt(msg); err != nil {
					return nil, err
				}
			}

			return msg, nil
		},
		IndexedKeys: parseRegexp.SubexpNames(),
	}
	return mp
}

//AsMap returns the message as a map[string]interface{}, for purposes of saving to tiedot
func (msg *Message) Map() map[string]interface{} {
	return map[string]interface{}{
		"timestamp": msg.Timestamp.UnixNano(),
		"data": msg.Data,
		"keys": msg.Keys,
	}
}
