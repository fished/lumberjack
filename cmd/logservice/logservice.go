package main

import (
	"fmt"
	"log"
	"net/http"

	"io/ioutil"
	"regexp"

	"github.com/fished/lumberjack/backend"
	"github.com/julienschmidt/httprouter"
)

var recorder *backend.Recorder
var parsers map[string]*backend.MessageParser

func postMessageHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()

	var body string
	defer r.Body.Close()
	if b, err := ioutil.ReadAll(r.Body); err == nil {
		body = string(b)
	} else {
		errorString := fmt.Sprintf("ERROR: Couldn't read request body: %s", err)
		log.Println(errorString)
		http.Error(w, errorString, http.StatusBadRequest)
	}

	var msg *backend.Message
	var err error
	parser := parsers["default"]
	if msg, err = parser.Parse(body,
		backend.MessageKey("source", r.FormValue("source")),
		backend.MessageKey("instance", r.FormValue("instance"))); err != nil {
		errorString := fmt.Sprintf("ERROR: Couldn't read parse message: %s", err)
		log.Println(errorString)
		http.Error(w, errorString, http.StatusBadRequest)
	}

	if err := recorder.Record(msg); err != nil {
		errorString := fmt.Sprintf("ERROR: Couldn't record msg: %s", err)
		log.Println(errorString)
		http.Error(w, errorString, http.StatusInternalServerError)
	}
}

func init() {

}

func main() {
	var err error

	log.Print("Setting up service...")
	// Setup the parsers.  TODO: Should be configurable.
	parsers = make(map[string]*backend.MessageParser)
	parsers["default"] = backend.NewStringMessageParser(
		regexp.MustCompile(`^(?P<timestamp>.{15})\s+((?P<instance>\S+)\s+(?P<source>.+?)\[(?P<pid>.*?)\]:)?`),
		backend.TimestampParserRegistry["Stamp"],
		"default",
		"default")

	// Generate the recorder as a global.  TODO: Should be configurable.
	var indexes []string
	for _, i := range parsers {
		indexes = append(indexes, i.IndexedKeys...)
	}

	//TODO: Should be configurable.
	recorder, err = backend.NewRecorder("/tmp/lumberjack", indexes)
	if err != nil {
		log.Fatalf("Failed to initializer recorder: %s", err)
	}
	log.Println("done.")

	log.Println("Serving.")
	router := httprouter.New()
	router.POST("/log/:source", postMessageHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
