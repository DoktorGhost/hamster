package main

import (
	"ham/internal/handlers"
	"log"
	"net/http"
)

func main() {
	handlers.InitRout()
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
