package logger

import (
	"log"
	"os"
)

// Logger simple wrapper untuk logging
var (
	Info  *log.Logger
	Error *log.Logger
)

// Init inisialisasi logger
func Init() {
	Info = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
