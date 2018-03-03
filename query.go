package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type graphQLResponse struct {
	Data   *interface{}  `json:"data"`
	Errors []interface{} `json:"error"`
}

func GraqhQLRequest(endpoint string, body string, args map[string]interface{}) (string, error) {

	//
	// Prepare Query
	//
	query := make(map[string]interface{})
	query["query"] = body
	query["variables"] = args
	marshaled, e := json.Marshal(query)
	if e != nil {
		return "", e
	}
	r, e := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(marshaled))
	if e != nil {
		return "", e
	}
	r.Header.Set("Content-Type", "application/json")

	//
	// Executing
	//

	client := http.Client{}
	response, e := client.Do(r)
	if e != nil {
		return "", e
	}
	defer response.Body.Close()
	responseBody, e := ioutil.ReadAll(response.Body)
	if e != nil {
		return "", e
	}

	//
	// Parsing
	//

	var responseText graphQLResponse
	e = json.Unmarshal(responseBody, &responseText)
	if e != nil {
		return "", e
	}

	//
	// Handle Errors
	//
	if len(responseText.Errors) > 0 {
		return "", fmt.Errorf("Errors: %v", responseText.Errors)
	}

	result, _ := json.Marshal(responseText.Data)
	return string(result), nil
}
