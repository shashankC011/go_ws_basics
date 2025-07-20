package main

import (
	"log"
	"net/http"
)

func main() {
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		log.Fatal("error: ", err)
	}
}
