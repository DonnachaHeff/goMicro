package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(rw http.ResponseWriter, r*http.Request){
		log.Println("Hello World")
		data, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(rw, "Request Failed", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(rw, "Hello %s\n", data)
	})

	http.HandleFunc("/goodbye", func(http.ResponseWriter, *http.Request){
		log.Println("Goodbye World")
	})

	http.ListenAndServe(":9090", nil)
}
