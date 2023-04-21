package main

import (
  "log"
  "net/http"
)

// home is a handler function
func home(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Hello from snippetbox"))
  

  // if we want to show 404 error for undefined paths, we can add this condition to the catch all case.

  if r.URL.Path != "/" {
    http.NotFound(w,r)
    return
  }
}

func snippetView(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Display specific snippet...."))
}

func snippetCreate(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Create a new snippet..."))
}

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", home)
  mux.HandleFunc("/snippet/view", snippetView)
  mux.HandleFunc("/snippet/create", snippetCreate)

  log.Println("Starting server on 4000:")
  err := http.ListenAndServe(":4000", mux)

  log.Fatal(err)

}
