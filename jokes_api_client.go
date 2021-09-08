package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getJoke() string {
	response, err := http.Get("https://api.chucknorris.io/jokes/random")

	const errorResponse string = "I ain't got time for this. No jokes for you today!"

	if err != nil {
		fmt.Println("Error when calling jokes API: ", err.Error())
		return errorResponse
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error when reading jokes API output: ", err.Error())
		return errorResponse
	}

	var responseObject Response
	json.Unmarshal(responseData, &responseObject)

	return responseObject.Value
}

type Response struct {
	Value      string   `json:"value"`
	Categories []string `json:"categories"`
	CreatedAt  string   `json:"created_at"`
	UpdatedAt  string   `json:"updated_at"`
	Url        string   `json:"url"`
}
