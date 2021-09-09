package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func getYesNoAnswer() (apiResponse YesNoApiResponse, err error) {
	response, err := http.Get("https://yesno.wtf/api/")

	if err != nil {
		fmt.Println("Error when calling yes/no API: ", err.Error())
		return apiResponse, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error when reading yes/no output: ", err.Error())
		return apiResponse, nil
	}

	json.Unmarshal(responseData, &apiResponse)

	return apiResponse, nil
}

type YesNoApiResponse struct {
	Answer string `json:"answer"`
	Forced bool   `json:"forced"`
	Image  string `json:"image"`
}
