package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// UnpackJSONBody returns the JSON body posted
// in an HTTP request
func UnpackJSONBody(req *http.Request, in interface{}) error {

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, in)
	if err != nil {
		return err
	}

	return nil
}
