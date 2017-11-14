package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	mu    = sync.RWMutex{}
	store = make(map[string]interface{})
	port  string
)

type Proto struct {
	ID   string
	Data interface{}
}

func main() {
	port = os.Getenv("PORT")
	if port == "" {
		flag.StringVar(&port, "p", "8080", "listen port")
		flag.Parse()
	}

	http.HandleFunc("/write", writeJSON)
	http.HandleFunc("/read", readJSON)

	log.Println("mockstoragego started on Port: ", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("HTTP Server Failed: ", err)
	}
}

func writeJSON(w http.ResponseWriter, req *http.Request) {
	log.Println("writeJSON Endpoint: ", req.RemoteAddr)
	if req.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "we only speak json", http.StatusUnsupportedMediaType)
		return
	}
	if req.Body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}

	var proto *Proto
	err = json.Unmarshal(body, &proto)
	if err != nil || proto.ID == "" || proto.Data == "" {
		log.Println("Request Error")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	mu.Lock()
	store[proto.ID] = proto.Data
	mu.Unlock()
	fmt.Println(proto.Data)
}

func readJSON(w http.ResponseWriter, req *http.Request) {
	log.Println("readJSON Endpoint: ", req.RemoteAddr)
	if req.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	if req.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		http.Error(w, "We Only speak URL Encoded in GET responses", http.StatusUnsupportedMediaType)
		return
	}

	id := req.FormValue("id")
	if id == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}

	mu.RLock()
	source, found := store[id]
	mu.RUnlock()
	if !found {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	b, err := json.Marshal(source)
	if err != nil {
		log.Println("Serialize Error")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	fmt.Println("response: ", b)
	fmt.Fprint(w, string(b))
}
