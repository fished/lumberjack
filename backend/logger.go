package backend

import (
	"encoding/json"
	"log"
	"regexp"
	"time"

	"github.com/HouzuoGuo/tiedot/db"
	"github.com/fished/lumberjack/util"
)

func hasIndex(col *db.Col, name string) bool {
	for _, ii := range col.AllIndexes() {
		if ii[len(ii)-1] == name {
			return true
		}
	}
	return false
}

func newLogger(dbPath string, parser *regexp.Regexp, tsParser TimestampParser) (func(string) error, error) {
	// (Create if not exist) open a database
	DB, err := db.OpenDB(dbPath)
	if err != nil {
		return nil, err
	}

	var logs *db.Col
	if logs = DB.Use("logs"); logs == nil {
		if err := DB.Create("logs"); err != nil {
			return nil, err
		}
		if logs = DB.Use("logs"); logs == nil {
			return nil, err
		}
	}

	for _, name := range parser.SubexpNames() {
		if !hasIndex(logs, name) {
			if err := logs.Index([]string{"message", "keys", name}); err != nil {
				return nil, err
			}
		}
	}

	return func(msg string) error {
		re := parser.Copy() // Ensure that we avoid a locking situation

		message := Message{
			Timestamp: time.Now(),
			Data:      msg,
			Keys:      make(map[string]string),
		}

		/* Find indexable patters */
		matches := re.FindAllStringSubmatch(msg, -1)
		indexable := util.DecodeNamedSubmatches(matches, re)
		log.Printf("%#v", indexable)

		for _, i := range indexable {
			if i == nil {
				continue
			}

			switch i[0] {
			case "timestamp":
				ts, err := tsParser.Parse(i[1])
				if err == nil {
					message.Timestamp = ts
				} else {
					// If we can't parse the timestamp, we just use time.Now()
					log.Printf("Couldn't parse timestring '%s': %s", i[1], err)
				}

				log.Printf("Using timestamp %s", message.Timestamp)
			case "":
				continue
			default:
				log.Printf("Would have indexed %s=%s", i[0], i[1])
				message.Keys[i[0]] = i[1]
			}
		}

		line, err := json.MarshalIndent(message, "", "  ")
		if err != nil {
			return err
		}

		log.Println(string(line))

		// file.Write([]byte(line))
		_, err = logs.Insert(map[string]interface{}{
			"message": message,
		})
		return err
	}, nil
}
