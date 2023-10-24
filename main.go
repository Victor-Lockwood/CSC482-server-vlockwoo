package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	gorillamuxRouter := mux.NewRouter()
	gorillamuxRouter.HandleFunc("/vlockwoo/status", status).Methods(http.MethodGet)

	http.ListenAndServe(":8090", gorillamuxRouter)
}

func status(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "hello this is me flying gorilla\n")
}
