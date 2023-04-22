package main

import (
	"log"
	"net/http"
)

func main() {

	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static"))
	// we strip the /static before we reach the file handler becauase
	// if we keep it then it searches in ./ui/static/static which in not present
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet/view", snippetView)
	mux.HandleFunc("/snippet/create", snippetCreate)

	log.Println("Starting server on 4000:")
	err := http.ListenAndServe(":4000", mux)

	log.Fatal(err)

}
