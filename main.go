package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	mux := mux.NewRouter()

	mux.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		json.NewEncoder(rw).Encode(map[string]interface{}{
			"status": "OK",
		})
	}).Methods("GET")

	fmt.Println("Server is running at PORT 8000")
	if err := http.ListenAndServe("0.0.0.0:8000", mux); err != nil {
		log.Fatal(err)
	}
}
