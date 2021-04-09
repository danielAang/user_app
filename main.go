package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/usuario/", UserHandler)
	fmt.Println("Running...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
