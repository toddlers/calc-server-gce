package main

import (
	"calc_util"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request, url string) {
	output := calc_util.ResultMsg{}

	a, err := parseFloatQueryParamAndTransform(r, "a", url)
	if err != nil {
		output.Error = err.Error()
		log.Println(output.Error)
	}

	b, err := parseFloatQueryParamAndTransform(r, "b", url)
	if err != nil {
		output.Error = err.Error()
		log.Println(output.Error)
	}

	if output.Error == "" {
		output.Result = a + b
	}

	calc_util.SendOutput(w, output)
}

func parseFloatQueryParamAndTransform(r *http.Request, paramName, url string) (num float64, err error) {
	num, err = calc_util.ParseFloatQueryParam(r, paramName)
	if err != nil {
		return 0, err
	}

	num, err = calc_util.CallUrlAndReturnFloat(fmt.Sprintf("%s?a=%f", url, num))
	if err != nil {
		return 0, err
	}

	return num, nil
}

func main() {
	port := flag.Int("port", 0, "TCP port for the HTTP server to listen on")
	url := flag.String("squareServerUrl", "", "URL to the 'calc_server_square' instance")

	flag.Parse()

	http.HandleFunc("/compute/add", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, *url)
	})

	http.ListenAndServe(":" + strconv.Itoa(*port), nil)
}
