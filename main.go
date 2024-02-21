package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

var startTime time.Time
var hostname string

func uptime() time.Duration {
	return time.Since(startTime)
}

func init() {
	startTime = time.Now()
	hostname, _ = os.Hostname()
}

func status(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Server", hostname)
	resp := make(map[string]string)
	resp["uptime"] = uptime().String()
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}

	w.Write(jsonResp)
}

func handleRequest(w http.ResponseWriter, req *http.Request) {
	log.Printf("Handling request : %s %s %s %s headers(%s)\n", hostname, req.Host, req.Method, req.URL.Path, req.Header)
	w.Header().Add("x-server", hostname)
	resp := make(map[string]string)
	resp["message"] = "success"
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func main() {
	address, ok := os.LookupEnv("LISTEN_ADDRESS")
	if !ok {
		log.Println("Listen address not defined.  Defaulting to 0.0.0.0")
		address = "0.0.0.0"
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		log.Println("Port not defined.  Defaulting to 8000")
		port = "8000"
	}

	http.HandleFunc("/status", status)
	http.HandleFunc("/", handleRequest)

	log.Println("Starting server : ", address+":"+port)
	log.Println("Version 0.1")

	http.ListenAndServe(address+":"+port, nil)
}
