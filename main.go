package main

import (
  "net/http"
  "github.com/go-chi/chi"
  "log"
)

func main() {
  r := chi.NewRouter()
  r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    log.Printf("Received request")
    w.Write([]byte("welcome"))
  })
  http.ListenAndServe(":3000", r)
}
