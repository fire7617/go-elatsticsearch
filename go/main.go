package main

import (
	"fmt"
	"net/http"
	"io"
	"os"
	"time"
    "context"
    "bytes"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v8"
    "github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/sirupsen/logrus"
)


type ElasticDoc struct {
    Level       string  `json:"level"`
    Message     string  `json:"message"`
    CreateAt    string  `json:"create_at"`
    Url         string  `json:"url"`
}

// formatting and printing values to the console.

// Used for build HTTP servers and clients.

const port string = ":8080"

var log = logrus.New()
var esClient *elasticsearch.Client

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
	log.SetLevel(logrus.DebugLevel)           // 设置输出警告级别
	// Output to stdout instead of the default stderr

	log.WithFields(logrus.Fields{
		"message": "logWithField",
	}).Info("a walrus appears")


    cfg := elasticsearch.Config{
    	Addresses: []string{
    		"http://localhost:9200",
    	},
    }

    esClient, err = elasticsearch.NewClient(cfg)
	if err != nil {
		panic("elasticsearch.NewTypedClient failed")
	}

}

func writeLog(text string) {
	log.WithFields(logrus.Fields{
		"message": text,
	}).Info("a walrus appears")




    doc := ElasticDoc{}
    doc.Level = "Test1"
    doc.Message = text
    doc.Url = "/home"
    doc.CreateAt = time.Now().Format("2006-1-2 15:4:5")

    jsonString, err := json.Marshal(&doc)
    if err != nil {
        panic("json Marshal eror")
    }


	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

// Set up the request object.
	req := esapi.IndexRequest{
		Index:      "courses",
        DocumentID: "33",
		Body:       bytes.NewReader(jsonString),
		Refresh:    "true",
	}


    // Perform the request with the client.
    res, err := req.Do(context.Background(), es)
    if err != nil {
    	log.Fatalf("Error getting response: %s", err)
    } else {
        fmt.Println("getting response: %s", res)
    }

    defer res.Body.Close()
}

// Handler functions.
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Homepage")
	writeLog("enter home")

	time.Sleep(20 * 1000 * time.Millisecond)
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
