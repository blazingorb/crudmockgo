package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	mockstorage "github.com/blazingorb/mockstoragego"
	"github.com/rs/cors"
)

var (
	port        string
	crossdomain string
	store       = mockstorage.NewMockStorage()
)

type Proto struct {
	ID   string
	Data interface{}
}

func main() {
	port = os.Getenv("PORT")
	crossdomain = os.Getenv("CORS")
	if port == "" {
		flag.StringVar(&port, "p", "8080", "listen port")
		flag.Parse()
	}
	if crossdomain == "" {
		flag.StringVar(&crossdomain, "c", "", "Allowed CORS Origin")
	}

	access := cors.AllowAll().Handler
	if crossdomain != "" {
		c := cors.Options{}
		c.AllowedOrigins = strings.Split(crossdomain, ",")
		c.AllowedMethods = []string{"POST", "GET", "PUT", "PATCH"}
		c.AllowedHeaders = []string{"Origin", "Content-Type", "Authorization"} //, "Accept", "Content-Type", "hello"}
		c.AllowCredentials = true
		access = cors.New(c).Handler
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/write", writeJSON)
	mux.HandleFunc("/read", readJSON)
	mux.HandleFunc("/list", listJSON)
	mux.HandleFunc("/clear", clear)

	log.Println("mockstoragego started on Port: ", port)
	err := http.ListenAndServe(":"+port, access(mux))
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
		log.Println("Request Error: ", string(body), proto, err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	store.Store(proto.ID, proto.Data)
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

	source := store.Load(id)
	if source == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	b, err := json.Marshal(source)
	if err != nil {
		log.Println("Serialize Error")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, string(b))
}

func listJSON(w http.ResponseWriter, req *http.Request) {
	log.Println("listJSON Endpoint: ", req.RemoteAddr)
	if req.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	list := store.List()
	b, err := json.Marshal(list)
	if err != nil {
		log.Println("Serialize Error")
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(b))
}

func clear(w http.ResponseWriter, req *http.Request) {
	log.Println("clear Endpoint: ", req.RemoteAddr)
	if req.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	store.Clear()
}
