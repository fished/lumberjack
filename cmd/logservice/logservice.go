package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fished/lumberjack/backend"
	"github.com/julienschmidt/httprouter"
	"io/ioutil"
	"regexp"
)

var recorder *backend.Recorder
var parsers map[string]*backend.MessageParser

func postMessageHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	r.ParseForm()

	var body string
	bodyReader := r.Body()
	defer bodyReader.Close()
	if b, err := ioutil.ReadAll(bodyReader); err == nil {
		body = string(b)
	} else {
		errorString := fmt.Sprintf("ERROR: Couldn't read request body: %s", err)
		log.Println(errorString)
		http.Error(w, errorString, http.StatusBadRequest)
	}

	parser := parsers["default"]
	msg := parser(body,
		backend.MessageKey("source", r.FormValue("source")),
		backend.MessageKey("instance", r.FormValue("instance")))

	if err := recorder(msg); err != nil {
		errorString := fmt.Sprintf("ERROR: Couldn't record msg: %s", err)
		log.Println(errorString)
		http.Error(w, errorString, http.StatusInternalServerError)
	}
}

func init() {

}

func main() {
	var err error

	// Setup the parsers.  TODO: Should be configurable.
	parsers = make(map[string]*backend.MessageParser)
	parsers["default"] = backend.NewStringMessageParser(
		regexp.MustCompile(`^(?P<timestamp>.{15})\s+((?P<instance>\S+)\s+(?P<source>.+?)\[(?P<pid>.*?)\]:)?`),
		backend.TimestampParserRegistry["Stamp"],
		"default",
		"default")

	// Generate the recorder as a global.  TODO: Should be configurable.
	var indexes []string
	for _, i := range(parsers) {
		indexes = append(indexes, i...)
	}
	//TODO: Shoudl be configurable.
	recorder, err = backend.NewRecorder("/tmp/lumberjack", indexes)
	if err != nil {
		log.Fatalf("Failed to initializer recorder: %s", err)
	}

	router := httprouter.New()
	router.POST("/log/:source", postMessageHandler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
