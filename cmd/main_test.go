package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteWithoutPostMethod(t *testing.T) {
	req, err := http.NewRequest("GET", "/writeJSON", nil)
	if err != nil {
		t.Error("Request Creation Failed: ", err)
	}

	reqr := httptest.NewRecorder()
	http.HandlerFunc(writeJSON).ServeHTTP(reqr, req)

	if status := reqr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Status code differs. Expected %d \n Got %d", http.StatusMethodNotAllowed, status)
	}
}

func TestWriteWithoutJSONHeader(t *testing.T) {
	req, err := http.NewRequest("POST", "/writeJSON", nil)
	if err != nil {
		t.Error("Request Creation Failed: ", err)
	}

	reqr := httptest.NewRecorder()
	http.HandlerFunc(writeJSON).ServeHTTP(reqr, req)

	if status := reqr.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf("Status code differs. Expected %d \n Got %d", http.StatusUnsupportedMediaType, status)
	}
}

func TestWriteWithEmptyBody(t *testing.T) {
	req, err := http.NewRequest("POST", "/writeJSON", nil)
	if err != nil {
		t.Error("Request Creation Failed: ", err)
	}
	req.Header.Set("Content-Type", "application/json")

	reqr := httptest.NewRecorder()
	http.HandlerFunc(writeJSON).ServeHTTP(reqr, req)

	if status := reqr.Code; status != http.StatusBadRequest {
		t.Errorf("Status code differs. Expected %d \n Got %d", http.StatusBadRequest, status)
	}
}

func TestReadWithoutGet(t *testing.T) {
	req, err := http.NewRequest("POST", "/readJSON", nil)
	if err != nil {
		t.Error("Request Creation Failed: ", err)
	}

	reqr := httptest.NewRecorder()
	http.HandlerFunc(readJSON).ServeHTTP(reqr, req)

	if status := reqr.Code; status != http.StatusMethodNotAllowed {
		t.Errorf("Status code differs. Expected %d \n Got %d", http.StatusMethodNotAllowed, status)
	}
}

func TestReadWithoutFormHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "/readJSON", nil)
	if err != nil {
		t.Error("Request Creation Failed: ", err)
	}

	reqr := httptest.NewRecorder()
	http.HandlerFunc(readJSON).ServeHTTP(reqr, req)

	if status := reqr.Code; status != http.StatusUnsupportedMediaType {
		t.Errorf("Status code differs. Expected %d \n Got %d", http.StatusUnsupportedMediaType, status)
	}
}

func TestReadWithoutID(t *testing.T) {
	req, err := http.NewRequest("GET", "/readJSON", nil)
	if err != nil {
		t.Error("Request Creation Failed: ", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	reqr := httptest.NewRecorder()
	http.HandlerFunc(readJSON).ServeHTTP(reqr, req)

	if status := reqr.Code; status != http.StatusBadRequest {
		t.Errorf("Status code differs. Expected %d \n Got %d", http.StatusBadRequest, status)
	}
}

func TestWriteAndReadSuccess(t *testing.T) {
	type Envelope struct {
		Name string
		Age  int
	}
	type Transport struct {
		ID   string
		Data Envelope
	}

	startData := &Transport{}
	startData.ID = "1234abc"
	startData.Data.Name = "Johnny"
	startData.Data.Age = 9999
	jsonStr, err := json.Marshal(startData)
	if err != nil {
		t.Error("Json Serialization Failed")
	}
	req, err := http.NewRequest("POST", "/writeJSON", bytes.NewBuffer(jsonStr))
	if err != nil {
		t.Error("Request Creation Failed: ", err)
	}
	req.Header.Set("Content-Type", "application/json")

	reqr := httptest.NewRecorder()
	http.HandlerFunc(writeJSON).ServeHTTP(reqr, req)

	if status := reqr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d \n Got %d", http.StatusOK, status)
	}

	//Begin Read Test
	req, err = http.NewRequest("GET", "/readJSON", nil)
	if err != nil {
		t.Error("Request Creation Failed: ", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	q := req.URL.Query()
	q.Add("id", startData.ID)
	req.URL.RawQuery = q.Encode()

	reqr = httptest.NewRecorder()
	http.HandlerFunc(readJSON).ServeHTTP(reqr, req)

	if status := reqr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d \n Got %d", http.StatusOK, status)
	}

	body, err := ioutil.ReadAll(reqr.Body)
	if err != nil {
		t.Errorf("Request Body Read Failed: ", err)
	}

	var readData *Envelope
	err = json.Unmarshal(body, &readData)
	if err != nil {
		t.Errorf("JSON deserialize failed: ", err)
	}

	if startData.Data.Name != readData.Name {
		t.Errorf("Response Differs. Expected %s \n Got %s", startData.Data.Name, readData.Name)
	}
	if startData.Data.Age != readData.Age {
		t.Errorf("Response Differs. Expected %s \n Got %s", startData.Data.Age, readData.Age)
	}
}
