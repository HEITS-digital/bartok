package client

import (
	"bartok/internal/dto"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func GetAnswer() (apiResponse dto.TruthApiResponse, err error) {
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

	err = json.Unmarshal(responseData, &apiResponse)

	return apiResponse, err
}
