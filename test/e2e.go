package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var domain = "http://localhost:8000"
var token1 = ""
var username = "Thomas"
var password = "Langtkodeord22"

func main() {

	register(username, password)

	login(username, password)

	oldPassword := password
	password = "vildtlangtpassword1234"
	updatePassword(password)

	login(username, password)

	login(username, oldPassword)

	updateName("NyThomas")

	login(username, password)

	login("Thomas", password)

	println("\nLogout")
	logout()

	println("\n")
}

func login(username string, password string) {

	println("\nLogin user", username, "with password", password)
	body := map[string]string{"username": username, "password": password}
	req := createRequestWithJSONBody("POST", "/login", body, false)

	resp := sendRequest(req)

	if resp.StatusCode == 200 {
		fmt.Println(resp.Status)
		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		var jsonBody struct {
			Token string `json:"token"`
		}

		json.Unmarshal(body, &jsonBody)
		token1 = jsonBody.Token
	} else {
		printResponse(resp)
	}
}

func register(username string, password string) {

	println("\nRegister user", username, "with password", password)
	body := map[string]string{"username": username, "password": password}
	req := createRequestWithJSONBody("POST", "/register", body, false)

	resp := sendRequest(req)
	printResponse(resp)
}

func updatePassword(newPassword string) {

	println("\nUpdate password for user", username, "to", newPassword)
	body := map[string]string{"password": newPassword}
	req := createRequestWithJSONBody("POST", "/update/password", body, true)

	resp := sendRequest(req)
	if resp.StatusCode == http.StatusOK {
		password = newPassword
	}

	printResponse(resp)
}

func updateName(newUsername string) {

	println("\nUpdate name for user", username, "to", newUsername)
	body := map[string]string{"username": newUsername}
	req := createRequestWithJSONBody("POST", "/update/name", body, true)

	resp := sendRequest(req)
	if resp.StatusCode == http.StatusOK {
		username = newUsername
	}

	printResponse(resp)
}

func logout() {

	req := createRequest("GET", "/logout", true)
	resp := sendRequest(req)
	printResponse(resp)
}

func createRequest(method string, path string, header bool) *http.Request {

	req := createRequestWithJSONBody(method, path, nil, header)
	req.Header.Add("Content-Type", "text/plain")

	return req
}

func createRequestWithJSONBody(method string, path string, body map[string]string, header bool) *http.Request {

	urlStr := domain + path
	jsonValue, _ := json.Marshal(body)

	req, err := http.NewRequest(method, urlStr, bytes.NewBuffer(jsonValue))

	if err != nil {
		panic("Failed to initialize request")
	}

	if header && token1 != "" {
		req.Header.Add("Authorization", "Bearer "+token1)
	}

	if body != nil {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", strconv.Itoa(len(jsonValue)))
	}

	return req
}

func sendRequest(req *http.Request) *http.Response {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	return resp
}

func printResponse(resp *http.Response) {
	fmt.Println(resp.Status)

	body, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	if len(body) > 0 {
		fmt.Println(string(body))
	}
}
