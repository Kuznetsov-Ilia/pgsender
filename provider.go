package pgsender

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	login    = os.Getenv("sms_login")
	password = os.Getenv("sms_password")
	host     = os.Getenv("sms_host")
	timeout  = os.Getenv("sms_timeout")
	auth64   = Encode(login, password)
)

// Encode encodes login and password to basic auth header
func Encode(login string, password string) string {
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(login+":"+password))
}

// Post to sms gate endpoint
func Post(url string, body *bytes.Reader) (int, []byte) {
	request, err := http.NewRequest("POST", host+url, body)
	Handle(err, "Error posting "+host+url)
	return API(request)
}

// Get sms gate api endpoint
func Get(url string) (int, []byte) {
	request, err := http.NewRequest("GET", host+url, nil)
	Handle(err, "Error getting "+host+url)
	return API(request)
}

// API method to request all endpoints
func API(request *http.Request) (int, []byte) {
	timeout, err := strconv.Atoi(timeout)
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
	request.Header.Set("Accept", "application/json, text/plain, */*")
	request.Header.Set("Authorization", auth64)
	request.Header.Set("Content-Type", "application/json;charset=utf-8")
	responce, err := client.Do(request)
	Handle(err, "Error posting to sms gate")
	jsonData, err := ioutil.ReadAll(responce.Body)
	Handle(err, "Error reading from responce Body")
	status := responce.StatusCode
	defer responce.Body.Close()
	return status, jsonData
}

// Responce common structure
type Responce struct {
	BatchID      string `json:"batch_id"`
	ErrorMessage string `json:"error_message"`
}

// Unmarshal Responce
func Unmarshal(jsonData []byte) Responce {
	var responce Responce
	err := json.Unmarshal([]byte(jsonData), &responce)
	Handle(err, "Error unmarhalling")
	return responce
}
