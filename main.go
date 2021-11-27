package main

import (
	"fmt"
	"go-postgres/router"
	"log"
	"net/http"
)

func main() {
	r := router.Router()
	// fs := http.FileServer(http.Dir("build"))
	// http.Handle("/", fs)
	fmt.Println("Starting server on the port 6070...")

	log.Fatal(http.ListenAndServe(":6070", r))
}
