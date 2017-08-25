package calc_util

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type ResultMsg struct {
	Result float64	`json:"result"`
	Error string	`json:"error"`
}

func SendOutput(w http.ResponseWriter, output ResultMsg) {
	jsonStr, err := json.MarshalIndent(output, "", "    ")
	if err != nil {
		panic("cannot marshal JSON: " + err.Error())
	}
	log.Printf("output: %s", jsonStr)

	if output.Error != "" {
		log.Printf("status code: %d", http.StatusBadRequest)
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		log.Printf("status code: %d", http.StatusOK)
	}
	w.Write(jsonStr)
}

func ParseFloatQueryParam(r *http.Request, paramName string) (num float64, err error) {
	param := r.FormValue(paramName)
	num, err = strconv.ParseFloat(param, 64)
	if err != nil {
		return 0, fmt.Errorf("cannot parse query param '%s' as a float64: '%s'", paramName, param)
	}
	log.Printf("%s = %f", paramName, num)

	return num, nil
}

func CallUrlAndReturnFloat(url string) (num float64, err error) {
	log.Println("GET '" + url + "'")
	res, err := http.Get(url)
	if err != nil {
		return 0, fmt.Errorf("cannot GET URL '%s': %s", url, err.Error())
	}

	jsonDecoder := json.NewDecoder(res.Body)

	var msg ResultMsg
	if err = jsonDecoder.Decode(&msg); err != io.EOF && err != nil {
		return 0, fmt.Errorf("cannot parse JSON document from URL '%s': %s", url, err.Error())
	}

	return msg.Result, nil
}
