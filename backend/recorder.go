package backend

import (
	"github.com/HouzuoGuo/tiedot/db"
)

type Recorder struct {
	Record func(*Message) error
}

func NewRecorder(dbPath string, indexes []string) (*Recorder, error) {
	var err error

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

	for _, name := range indexes {
		if !hasIndex(logs, name) {
			if err := logs.Index([]string{"keys", name}); err != nil {
				return nil, err
			}
		}
	}

	return &Recorder{
		func(msg *Message) error {
			_, err := logs.Insert(msg.Map())
			if err != nil {
				return err
			}
			return nil
		},
	}, nil
}

func hasIndex(col *db.Col, name string) bool {
	for _, ii := range col.AllIndexes() {
		if ii[len(ii)-1] == name {
			return true
		}
	}
	return false
}
