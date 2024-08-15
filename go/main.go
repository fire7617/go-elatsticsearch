package main

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"
    "os"
    "io"
    "time"
)

// formatting and printing values to the console.

// Used for build HTTP servers and clients.

const port string = ":8080"
var log = logrus.New()


func init() {
    //caution : format string is `2006-01-02 15:04:05.000000000`
    logFilename := time.Now().Format("2006-01-02.json.log")

    file, err := os.OpenFile(logFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0755)
    if err != nil {
        panic(err)
    }

    mw := io.MultiWriter(os.Stdout, file) // MultiWriter to log to stdout and file
	log.SetOutput(mw)

    log.SetFormatter(&logrus.JSONFormatter{}) // 设置 format json
	log.SetLevel(logrus.DebugLevel) // 设置输出警告级别
    // Output to stdout instead of the default stderr

	log.WithFields(logrus.Fields{
		"message": "logWithField",
	}).Info("a walrus appears")
}


func writeLog(text string) {
	log.WithFields(logrus.Fields{
		"message": text,
	}).Info("a walrus appears")
}

// Handler functions.
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage")
    writeLog("enter home")
}

func Info(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Info page")
	writeLog("enter info")
}

func main() {
	log.Println("Starting our simple http server.")

	// Registering our handler functions, and creating paths.
	http.HandleFunc("/", Home)
	http.HandleFunc("/info", Info)

	log.Println("Started on port", port)
	fmt.Println("To close connection CTRL+C :-)")

	// Spinning up the server.
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(err)
	}
}
