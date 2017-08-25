package main

import (
	"calc_util"
	"flag"
	"log"
	"net/http"
	"strconv"
)

func handler(w http.ResponseWriter, r *http.Request) {
	output := calc_util.ResultMsg{}

	param := r.FormValue("a")
	a, err := strconv.ParseFloat(param, 64)
	if err != nil {
		output.Error = "cannot parse query param 'a' as a float64: '" + param + "'"
		log.Println(output.Error)
	}
	log.Printf("a = %f", a)

	if output.Error == "" {
		output.Result = a * a
	}

	calc_util.SendOutput(w, output)
}

func main() {
	port := flag.Int("port", 0, "TCP port for the HTTP server to listen on")

	flag.Parse()

	http.HandleFunc("/compute/square", handler)

	http.ListenAndServe(":" + strconv.Itoa(*port), nil)
}
