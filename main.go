package main

import (
  "log"
  "net/http"
)

// home is a handler function
func home(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Hello from snippetbox"))
}

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/", home)


  log.Println("Starting server on 4000:")
  err := http.ListenAndServe(":4000", mux)

  log.Fatal(err)

}
