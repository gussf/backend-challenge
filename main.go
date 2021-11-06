package main

import (
	"fmt"
	"log"
	"net/http"
)

func checkoutHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		rw.Write([]byte("Only POST method is allowed"))
	}

	fmt.Println("Hello world")
}

func main() {
	http.HandleFunc("/checkout", checkoutHandler)
	log.Fatal(http.ListenAndServe(":3000", nil))
}
