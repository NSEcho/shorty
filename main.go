package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lateralusd/shorty/db"
	"github.com/lateralusd/shorty/handler"
)

func main() {
	var err error
	database, err := db.InitDatabase()
	if err != nil {
		log.Fatal("Coult not initialize database", err)
	}

	env := &handler.Env{
		DB: database,
	}

	http.Handle("/", handler.Handler{env, handler.IndexPath})
	http.Handle("/shorty", handler.Handler{env, handler.ShortyPath})

	fmt.Println("[*] Starting the server")
	http.ListenAndServe(":8080", nil)
}
