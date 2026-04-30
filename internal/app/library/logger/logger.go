package logger

import (
	"io"
	"log"
	"os"

	"github.com/tidwall/pretty"
)

var Error *log.Logger

func Init() {
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(
		"logs/errors.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}

	mw := io.MultiWriter(os.Stdout, file)

	Error = log.New(mw, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func PrintBodyJSON(v []byte) {
	prettyJSON := pretty.Pretty(v)
	log.Printf("%s", prettyJSON)
}
