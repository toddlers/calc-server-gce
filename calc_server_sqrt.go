package main

import (
	"calc_util"
	"flag"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request, url string) {
	output := calc_util.ResultMsg{}

	a, err := calc_util.ParseFloatQueryParam(r, "a")
	if err != nil {
		output.Error = err.Error()
		log.Println(output.Error)
	}

	b, err := calc_util.ParseFloatQueryParam(r, "b")
	if err != nil {
		output.Error = err.Error()
		log.Println(output.Error)
	}

	num, err := calc_util.CallUrlAndReturnFloat(fmt.Sprintf("%s?a=%f&b=%f", url, a, b))
	if err != nil {
		output.Error = err.Error()
		log.Println(output.Error)
	}

	if output.Error == "" {
		output.Result = math.Sqrt(num)
	}

	calc_util.SendOutput(w, output)
}

func main() {
	port := flag.Int("port", 0, "TCP port for the HTTP server to listen on")
	url := flag.String("adderServerUrl", "", "URL to the 'calc_server_add' instance")

	flag.Parse()

	http.HandleFunc("/compute/sqrt", func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, *url)
	})

	http.ListenAndServe(":" + strconv.Itoa(*port), nil)
}
