package client

import (
	"bartok/internal/dto"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetGeekJoke() string {
	response, err := http.Get("https://geek-jokes.sameerkumar.website/api?format=json")

	const errorResponse string = "I ain't got time for this. No jokes for you today!"

	if err != nil {
		fmt.Println("Error when calling geek joke API: ", err.Error())
		return errorResponse
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error when reading geek joke API output: ", err.Error())
		return errorResponse
	}

	var responseObject dto.JokeResponse
	_ = json.Unmarshal(responseData, &responseObject)

	return responseObject.Joke
}
